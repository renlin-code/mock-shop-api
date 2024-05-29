package domain

import (
	"errors"
	"fmt"
	"mime/multipart"
	"sort"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	userNameMinLength = 2
	userNameMaxLength = 10

	passwordMinLength = 4
	passwordMaxLength = 12

	maxFileSize = 10 << 20 //10 MB

	categoryNameMinLength        = 2
	categoryNameMaxLength        = 15
	categoryDescriptionMinLength = 0
	categoryDescriptionMaxLength = 200

	productNameMinLength        = 2
	productNameMaxLength        = 100
	productDescriptionMinLength = 0
	productDescriptionMaxLength = 200
)

var allowedFileExtensions = [3]string{"jpg", "jpeg", "png"}

type SignUpInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (i SignUpInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Length(userNameMinLength, userNameMaxLength)),
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

type UpdateProfileInput struct {
	Name           *string               `json:"name"`
	ProfileImgFile *multipart.FileHeader `json:"profile_image_file"`
}

func (i UpdateProfileInput) Validate() error {
	if i.Name == nil && i.ProfileImgFile == nil {
		return errors.New("no fields provided")
	}
	if i.ProfileImgFile != nil {
		return validateFile(i.ProfileImgFile, maxFileSize, allowedFileExtensions[:])
	}
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Length(userNameMinLength, userNameMaxLength)),
	)
}

type DeleteProfileInput struct {
	Password string `json:"password"`
}

func (i DeleteProfileInput) Validate() error {
	return validation.ValidateStruct(&i,
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
			return errors.New("product id must be unique")
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

func validateFile(file *multipart.FileHeader, maxSize int64, allowedExtensions []string) error {
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds max size (%d bytes)", maxSize)
	}
	ext := strings.ToLower(strings.Split(file.Filename, ".")[1])
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}
	if !validExtension {
		return fmt.Errorf("file extension must be .%s", strings.Join(allowedExtensions, "/."))
	}
	return nil
}

type CreateCategoryInput struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	ImgFile     *multipart.FileHeader `json:"image_file"`
	Available   bool                  `json:"available"`
}

func (i CreateCategoryInput) Validate() error {
	err := validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Length(categoryNameMinLength, categoryNameMaxLength)),
		validation.Field(&i.Description, validation.Length(categoryDescriptionMinLength, categoryDescriptionMaxLength)),
	)
	if err != nil {
		return err
	}

	if i.ImgFile == nil || i.ImgFile.Size == 0 {
		return errors.New("image_file: can not be blank")
	}

	err = validateFile(i.ImgFile, maxFileSize, allowedFileExtensions[:])
	if err != nil {
		return err
	}
	return nil
}

type UpdateCategoryInput struct {
	Name        *string               `json:"name"`
	Description *string               `json:"description"`
	ImgFile     *multipart.FileHeader `json:"image_file"`
	Available   *bool                 `json:"available"`
}

func (i UpdateCategoryInput) Validate() error {
	if i.Name == nil && i.Description == nil && i.ImgFile == nil && i.Available == nil {
		return errors.New("no fields provided")
	}
	if i.ImgFile != nil {
		return validateFile(i.ImgFile, maxFileSize, allowedFileExtensions[:])
	}

	err := validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Length(categoryNameMinLength, categoryNameMaxLength)),
		validation.Field(&i.Description, validation.Length(categoryDescriptionMinLength, categoryDescriptionMaxLength)),
	)
	if err != nil {
		return err
	}
	return nil
}

type CreateProductInput struct {
	CategoryId        int                   `json:"category_id"`
	Name              string                `json:"name"`
	Description       string                `json:"description"`
	ImgFile           *multipart.FileHeader `json:"image_file"`
	Available         bool                  `json:"available"`
	Price             float32               `json:"price"`
	UndiscountedPrice float32               `json:"undiscounted_price"`
	Stock             int                   `json:"stock"`
}

func (i CreateProductInput) Validate() error {
	err := validation.ValidateStruct(&i,
		validation.Field(&i.CategoryId, validation.Required, validation.Min(1)),
		validation.Field(&i.Name, validation.Required, validation.Length(productNameMinLength, productNameMaxLength)),
		validation.Field(&i.Description, validation.Length(productDescriptionMinLength, productDescriptionMaxLength)),
		validation.Field(&i.Price, validation.Required, validation.Min(0.0)),
		validation.Field(&i.UndiscountedPrice, validation.Required, validation.Min(0.0)),
		validation.Field(&i.Stock, validation.Required, validation.Min(0)),
	)
	if err != nil {
		return err
	}

	if i.ImgFile == nil || i.ImgFile.Size == 0 {
		return errors.New("image_file: can not be blank")
	}

	err = validateFile(i.ImgFile, maxFileSize, allowedFileExtensions[:])
	if err != nil {
		return err
	}
	return nil
}

type UpdateProductInput struct {
	CategoryId        *int                  `json:"category_id"`
	Name              *string               `json:"name"`
	Description       *string               `json:"description"`
	ImgFile           *multipart.FileHeader `json:"image_file"`
	Available         *bool                 `json:"available"`
	Price             *float32              `json:"price"`
	UndiscountedPrice *float32              `json:"undiscounted_price"`
	Stock             *int                  `json:"stock"`
}

func (i UpdateProductInput) Validate() error {
	if i.CategoryId == nil &&
		i.Name == nil &&
		i.Description == nil &&
		i.ImgFile == nil &&
		i.Available == nil &&
		i.Price == nil &&
		i.UndiscountedPrice == nil &&
		i.Stock == nil {
		return errors.New("no fields provided")
	}
	if i.ImgFile != nil {
		return validateFile(i.ImgFile, maxFileSize, allowedFileExtensions[:])
	}
	err := validation.ValidateStruct(&i,
		validation.Field(&i.CategoryId, validation.Min(1)),
		validation.Field(&i.Name, validation.Length(productNameMinLength, productNameMaxLength)),
		validation.Field(&i.Description, validation.Length(productDescriptionMinLength, productDescriptionMaxLength)),
		validation.Field(&i.Price, validation.Min(0.0)),
		validation.Field(&i.UndiscountedPrice, validation.Min(0.0)),
		validation.Field(&i.Stock, validation.Min(0)),
	)
	if err != nil {
		return err
	}
	return nil
}

type PaginationParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (p PaginationParams) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Page, validation.Min(1)),
		validation.Field(&p.PageSize, validation.Min(1), validation.Max(100)),
	)
}
