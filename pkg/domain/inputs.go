package domain

import (
	"errors"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	nameMinLength     = 2
	nameMaxLength     = 10
	passwordMinLength = 4
	passwordMaxLength = 12
)

type SignUpInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (i SignUpInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Length(nameMinLength, nameMaxLength)),
		validation.Field(&i.Email, validation.Required, is.Email),
	)
}

type ConfirmEmailInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (i ConfirmEmailInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Token, validation.Required),
		validation.Field(&i.Password, validation.Required, validation.Length(passwordMinLength, passwordMaxLength)),
	)
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i SignInInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email, validation.Required, is.Email),
		validation.Field(&i.Password, validation.Required, validation.Length(passwordMinLength, passwordMaxLength)),
	)
}

type UpdateProfileInput struct {
	Name       *string `json:"name"`
	ProfileImg *string `json:"profile_image"`
}

func (i UpdateProfileInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Length(nameMinLength, nameMaxLength)),
		validation.Field(&i.ProfileImg),
	)
}

type RecoveryPasswordInput struct {
	Email string `json:"email"`
}

func (i RecoveryPasswordInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email, validation.Required, is.Email),
	)
}

type UpdatePasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (i UpdatePasswordInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Token, validation.Required),
		validation.Field(&i.Password, validation.Required, validation.Length(passwordMinLength, passwordMaxLength)),
	)
}

type CreateOrderInput struct {
	Products []CreateOrderInputProduct `json:"products"`
}
type CreateOrderInputProduct struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

func (i CreateOrderInput) Validate() error {
	err := validation.ValidateStruct(&i,
		validation.Field(&i.Products, validation.Required),
	)
	if err != nil {
		return err
	}
	uniqueIDs := make(map[int]struct{})

	for _, product := range i.Products {
		if _, found := uniqueIDs[product.Id]; found {
			return errors.New("product id must to be unique")
		}
		uniqueIDs[product.Id] = struct{}{}
		err := validation.ValidateStruct(&product,
			validation.Field(&product.Id, validation.Required, validation.Min(1)),
			validation.Field(&product.Quantity, validation.Required, validation.Min(1)),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type ById []CreateOrderInputProduct

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (i *CreateOrderInput) Sort() {
	sort.Sort(ById(i.Products))
}
