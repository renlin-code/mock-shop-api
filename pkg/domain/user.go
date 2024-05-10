package domain

type User struct {
	Id         int    `json:"-" db:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	ProfileImg string `json:"profile_image" db:"profile_image"`
}
