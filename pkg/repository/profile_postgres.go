package repository

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
)

type ProfilePostgres struct {
	db *sqlx.DB
	s  *storage.Storage
}

func newProfilePostgres(db *sqlx.DB, s *storage.Storage) *ProfilePostgres {
	return &ProfilePostgres{db, s}
}

func (r *ProfilePostgres) GetProfile(userId int) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf(`SELECT 
		id, 
		name, 
		email, 
		profile_image 
	FROM %s WHERE id=$1`, usersTable)
	err := r.db.Get(&user, query, userId)
	if err == sql.ErrNoRows {
		return user, errors_handler.NoRows()
	}

	return user, err
}

func (r *ProfilePostgres) GetFilePath(userId int, fileName string) string {
	return r.s.Profile.GetFilePath(userId, fileName)
}

func (r *ProfilePostgres) UpdateProfile(userId int, input domain.UpdateProfileInput, file multipart.File) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.ProfileImgFile != nil && file != nil {
		url, err := r.s.UploadProfileImage(userId, input.ProfileImgFile, file)

		if err != nil {
			return err
		}

		setValues = append(setValues, fmt.Sprintf("profile_image=$%d", argId))
		args = append(args, url)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d RETURNING id", usersTable, setQuery, argId)
	args = append(args, userId)

	var id int
	err = tx.QueryRow(query, args...).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors_handler.NoRows()
		}

		return err
	}

	return tx.Commit()
}

func (r *ProfilePostgres) CreateOrder(userId int, products []domain.CreateOrderInputProduct) (int, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var orderId int
	orderDate := time.Now()
	createOrderQuery := fmt.Sprintf(`INSERT INTO %s (
		user_id, 
		date
	) VALUES ($1, $2) RETURNING id`, ordersTable)
	row := tx.QueryRow(createOrderQuery, userId, orderDate)
	if err := row.Scan(&orderId); err != nil {
		return 0, err
	}
	createOrderedProductsQuery := fmt.Sprintf(`INSERT INTO %s (
		order_id, 
		product_id, 
		name, 
		description, 
		price, 
		undiscounted_price, 
		image_url, 
		quantity
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, orderedProductsTable)

	stmt, err := tx.Prepare(createOrderedProductsQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, product := range products {
		updateStockQuery := fmt.Sprintf(`UPDATE %s SET stock = stock - $1 WHERE id = $2 RETURNING 
			name, 
			description, 
			price, 
			undiscounted_price, 
			image_url, 
			stock
		`, productsTable)
		var productFromTable domain.Product

		row := tx.QueryRow(updateStockQuery, product.Quantity, product.Id)
		err := row.Scan(
			&productFromTable.Name,
			&productFromTable.Description,
			&productFromTable.Price,
			&productFromTable.UndiscountedPrice,
			&productFromTable.ImageUrl,
			&productFromTable.Stock)

		if err != nil {
			if err == sql.ErrNoRows {
				return 0, errors_handler.NoRows()
			}
			if pqErr, ok := err.(*pq.Error); ok && strings.Contains(pqErr.Message, "violates check constraint \"stock\"") {
				return 0, errors_handler.ConstrainViolation("stock")
			}
			return 0, err
		}
		_, err = stmt.Exec(
			orderId,
			product.Id,
			productFromTable.Name,
			productFromTable.Description,
			productFromTable.Price,
			productFromTable.UndiscountedPrice,
			productFromTable.ImageUrl,
			product.Quantity)
		if err != nil {
			return 0, err
		}
	}
	return orderId, tx.Commit()
}

func (r *ProfilePostgres) GetAllOrders(userId, limit, offset int) ([]domain.Order, error) {
	query := fmt.Sprintf(`
		WITH order_total_cost AS (
			SELECT order_id, SUM(price * quantity) AS total_cost
				FROM %s opt
				GROUP BY opt.order_id
				ORDER BY opt.order_id
		)
		SELECT 
			ot.id AS order_id,     
			ot.user_id, 
			ot.date,     
			opt.id, 
			opt.product_id,     
			opt.name, 
			opt.description,     
			opt.price, 
			opt.undiscounted_price,     
			opt.image_url, 
			opt.quantity,
			otct.total_cost
		FROM 
			( SELECT * 
			FROM %s ot
			WHERE ot.user_id = $1
			ORDER BY ot.id
			LIMIT $2
			OFFSET $3
			) ot
		INNER JOIN %s opt
		ON ot.id = opt.order_id
		INNER JOIN order_total_cost otct
		ON otct.order_id = ot.id
		ORDER BY ot.id, opt.id;
	`, orderedProductsTable, ordersTable, orderedProductsTable)

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var product domain.OrderedProduct
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.Date,
			&product.Id,
			&product.ProductId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.UndiscountedPrice,
			&product.ImageUrl,
			&product.Quantity,
			&order.TotalCost)
		if err != nil {
			return nil, err
		}

		order.Products = append(order.Products, product)

		var exists bool
		for i, o := range orders {
			if o.Id == order.Id {
				orders[i].Products = append(orders[i].Products, product)
				exists = true
				break
			}
		}

		if !exists {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (r *ProfilePostgres) GetOrderById(userId, orderId int) (domain.Order, error) {
	query := fmt.Sprintf(`
		WITH order_total_cost AS (
			SELECT order_id, SUM(price * quantity) AS total_cost
				FROM %s op
				WHERE order_id=$1
				GROUP BY order_id
		)
			SELECT 
				ot.id AS order_id,     
				ot.user_id, 
				ot.date,     
				opt.id, 
				opt.product_id,     
				opt.name, 
				opt.description,     
				opt.price, 
				opt.undiscounted_price,     
				opt.image_url, 
				opt.quantity,
				otct.total_cost
			FROM %s ot
			INNER JOIN %s opt
			ON opt.order_id = ot.id 
			INNER JOIN order_total_cost otct
			ON otct.order_id = ot.id
			WHERE ot.id=$1 AND ot.user_id=$2
		`, orderedProductsTable, ordersTable, orderedProductsTable)

	var order domain.Order

	rows, err := r.db.Query(query, orderId, userId)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.OrderedProduct
		err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.Date,
			&product.Id,
			&product.ProductId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.UndiscountedPrice,
			&product.ImageUrl,
			&product.Quantity,
			&order.TotalCost)
		if err != nil {
			return order, err
		}
		product.OrderId = orderId
		order.Products = append(order.Products, product)
	}
	if order.Id == 0 {
		return order, errors_handler.NoRows()
	}
	return order, nil
}

func (r *ProfilePostgres) DeleteProfile(userId int, password string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1 AND password_hash=$2 RETURNING id`, usersTable)
	err = tx.QueryRow(query, userId, password).Scan(&id)
	if err == sql.ErrNoRows {
		return errors_handler.NoRows()
	}

	err = r.s.DeleteProfileImage(userId)

	return err
}
