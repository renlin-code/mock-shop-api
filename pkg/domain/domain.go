package domain

type User struct {
	Id         int    `json:"-" db:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfileImg string `json:"profile_image"`
}

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url" db:"image_url"`
}

type Order struct {
	Id     int
	UserId int
	Date   string
}

type Product struct {
	Id          int     `json:"id" db:"id"`
	CategoryId  int     `json:"category_id" db:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SalePrice   float32 `json:"sale_price" db:"sale_price"`
	ImagesUrls  string  `json:"images_urls" db:"images_urls"`
	Available   bool    `json:"available"`
	Stock       int     `json:"stock"`
}
