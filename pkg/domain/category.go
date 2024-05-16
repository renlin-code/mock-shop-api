package domain

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Available   bool   `json:"available"`
	ImageUrl    string `json:"image_url" db:"image_url"`
}
