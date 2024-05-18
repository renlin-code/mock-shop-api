package domain

type Product struct {
	Id                int     `json:"id" db:"id"`
	CategoryId        int     `json:"category_id" db:"category_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Price             float32 `json:"price"`
	UndiscountedPrice float32 `json:"undiscounted_price" db:"undiscounted_price"`
	ImageUrl          string  `json:"image_url" db:"image_url"`
	Available         bool    `json:"available"`
	Stock             int     `json:"stock"`
}
