package domain

type User struct {
	Id         int    `json:"-" db:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	ProfileImg string `json:"profile_image" db:"profile_image"`
}

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url" db:"image_url"`
}

type Product struct {
	Id                int     `json:"id" db:"id"`
	CategoryId        int     `json:"category_id" db:"category_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float32 `json:"price"`
	UndiscountedPrice float32 `json:"undiscounted_price" db:"undiscounted_price"`
	ImagesUrls        string  `json:"images_urls" db:"images_urls"`
	Available         bool    `json:"available"`
	Stock             int     `json:"stock"`
}

type Order struct {
	Id        int              `json:"id" db:"id"`
	UserId    int              `json:"user_id" db:"user_id"`
	Date      string           `json:"date"`
	Products  []OrderedProduct `json:"products" db:"products"`
	TotalCost float32          `json:"total_cost" db:"total_cost"`
}

type OrderedProduct struct {
	Id                int     `json:"id" db:"id"`
	OrderId           int     `json:"order_id" db:"order_id"`
	ProductId         int     `json:"product_id" db:"product_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float32 `json:"price"`
	UndiscountedPrice float32 `json:"undiscounted_price" db:"undiscounted_price"`
	ImagesUrls        string  `json:"images_urls" db:"images_urls"`
	Quantity          int     `json:"quantity"`
}
