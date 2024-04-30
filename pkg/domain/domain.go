package domain

type User struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfileImg string `json:"profile_image"`
}

type Category struct {
	Id          int
	Name        string
	Description string
	ImageUrl    string
}

type Order struct {
	Id     int
	UserId int
	Date   string
}

type Product struct {
	Id          int
	CategoryId  int
	Name        string
	Description string
	Price       float32
	SalePrice   float32
	ImagesUrl   []string
	Available   bool
	Stock       int
}
