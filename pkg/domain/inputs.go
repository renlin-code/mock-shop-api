package domain

import "sort"

type SignUpInput struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type ConfirmEmailInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,max=30"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,max=30"`
}

type UpdateProfileInput struct {
	Name       *string `json:"name"`
	ProfileImg *string `json:"profile_image"`
}

type UpdatePasswordInput struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateOrderInputProduct struct {
	Id       int `json:"id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
}

type CreateOrderInput struct {
	Products []CreateOrderInputProduct `json:"products" validate:"required"`
}

type ById []CreateOrderInputProduct

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (i *CreateOrderInput) Sort() {
	sort.Sort(ById(i.Products))
}
