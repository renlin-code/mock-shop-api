package storage

import (
	"mime/multipart"
)

type Profile interface {
	UploadProfileImage(userId int, handler *multipart.FileHeader, file multipart.File) (string, error)
	GetFilePath(userId int, fileName string) string
}

type Category interface {
	UploadCategoryImage(categoryId int, handler *multipart.FileHeader, file multipart.File) (string, error)
	GetFilePath(categoryId int, fileName string) string
}

type Product interface {
	UploadProductImage(productId int, handler *multipart.FileHeader, file multipart.File) (string, error)
	GetFilePath(productId int, fileName string) string
}

type Storage struct {
	Profile
	Category
	Product
}

func NewStorage(fs *FileSystemStorage) *Storage {
	return &Storage{
		Profile:  newProfileFileSystem(fs),
		Category: newCategoryFileSystem(fs),
		Product:  newProductFileSystem(fs),
	}
}
