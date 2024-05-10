package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

type ProfilePostgres struct {
	db *sqlx.DB
}

func newProfilePostgres(db *sqlx.DB) *ProfilePostgres {
	return &ProfilePostgres{db}
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

	return user, err
}

func (r *ProfilePostgres) UpdateProfile(userId int, input domain.UpdateProfileInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.ProfileImg != nil {
		setValues = append(setValues, fmt.Sprintf("profile_image=$%d", argId))
		args = append(args, *input.ProfileImg)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d", usersTable, setQuery, argId)
	args = append(args, userId)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfilePostgres) UpdatePassword(userId int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE id=$2", usersTable)

	_, err := r.db.Exec(query, password, userId)
	return err
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
		images_urls, 
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
			images_urls, 
			stock
		`, productsTable)
		var productFromTable domain.Product

		row := tx.QueryRow(updateStockQuery, product.Quantity, product.Id)
		err := row.Scan(
			&productFromTable.Name,
			&productFromTable.Description,
			&productFromTable.Price,
			&productFromTable.UndiscountedPrice,
			&productFromTable.ImagesUrls,
			&productFromTable.Stock)

		if err != nil {
			return 0, err
		}
		if productFromTable.Stock < 0 {
			return 0, errors.New("not enough stock")
		}

		_, err = stmt.Exec(
			orderId,
			product.Id,
			productFromTable.Name,
			productFromTable.Description,
			productFromTable.Price,
			productFromTable.UndiscountedPrice,
			productFromTable.ImagesUrls,
			product.Quantity)
		if err != nil {
			return 0, err
		}
	}
	return orderId, tx.Commit()
}

func (r *ProfilePostgres) GetAllOrders(userId int) ([]domain.Order, error) {
	query := fmt.Sprintf(`
		WITH order_products AS (
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
				opt.images_urls, 
				opt.quantity
			FROM %s ot
			INNER JOIN %s opt ON opt.order_id = ot.id
			WHERE ot.user_id = $1
		),
		order_total_cost AS (
			SELECT 
				order_id, 
				SUM(price) AS total_cost
			FROM order_products
			GROUP BY order_id
		)
		SELECT 
			op.order_id, 
			op.user_id, 
			op.date, 
			op.id, 
			op.product_id,
			op.name, 
			op.description, 
			op.price, 
			op.undiscounted_price, 
			op.images_urls, 
			op.quantity, 
			otc.total_cost
		FROM order_products op
		LEFT JOIN order_total_cost otc ON op.order_id = otc.order_id
	`, ordersTable, orderedProductsTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersMap = make(map[int]domain.Order)
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
			&product.ImagesUrls,
			&product.Quantity,
			&order.TotalCost)
		if err != nil {
			return nil, err
		}
		if existingOrder, found := ordersMap[order.Id]; found {
			existingOrder.Products = append(existingOrder.Products, product)
			ordersMap[order.Id] = existingOrder
		} else {
			order.Products = append(order.Products, product)
			ordersMap[order.Id] = order
		}
	}
	for _, order := range ordersMap {
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *ProfilePostgres) GetOrderById(userId, orderId int) (domain.Order, error) {
	query := fmt.Sprintf(`
		WITH order_products AS (
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
				opt.images_urls, 
				opt.quantity
			FROM %s ot
			INNER JOIN %s opt ON opt.order_id = ot.id
			WHERE ot.id = $1 AND ot.user_id = $2
		),
		order_total_cost AS (
			SELECT order_id, SUM(price) AS total_cost
			FROM order_products
			GROUP BY order_id
		)
		SELECT 
			op.order_id, 
			op.user_id, 
			op.date, 
			op.id, 
			op.product_id, 
			op.name, 
			op.description, 
			op.price, 
			op.undiscounted_price, 
			op.images_urls, 
			op.quantity, 
			otc.total_cost
		FROM order_products op
		LEFT JOIN order_total_cost otc ON op.order_id = otc.order_id
	`, ordersTable, orderedProductsTable)

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
			&product.ImagesUrls,
			&product.Quantity,
			&order.TotalCost)
		if err != nil {
			return order, err
		}
		product.OrderId = orderId
		order.Products = append(order.Products, product)
	}
	if order.Id == 0 {
		return order, errors.New("no rows")
	}
	return order, nil
}
