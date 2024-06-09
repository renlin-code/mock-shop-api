package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type ProductFileSystem struct {
	FileSystem *FileSystemStorage
}

func newProductFileSystem(fs *FileSystemStorage) *ProductFileSystem {
	return &ProductFileSystem{FileSystem: fs}
}

func (s *ProductFileSystem) UploadProductImage(productId int, handler *multipart.FileHeader, file multipart.File) (string, error) {
	productDir := fmt.Sprintf("%s/%s/%d/", basePath, productsDirectory, productId)

	err := os.RemoveAll("." + productDir)
	if err != nil {
		return "", errors_handler.StorageError("error deleting directory")
	}

	err = os.MkdirAll("."+productDir, os.ModePerm)
	if err != nil {
		return "", errors_handler.StorageError("error creating directory")
	}

	path := productDir + handler.Filename
	f, err := os.Create("." + path)
	if err != nil {
		return "", errors_handler.StorageError("error creating file")
	}
	defer f.Close()

	io.Copy(f, file)

	return fmt.Sprintf("%s/%s/%d/%s", s.FileSystem.config.MediaBaseUrl, productsDirectory, productId, handler.Filename), nil
}

func (s *ProductFileSystem) GetFilePath(productId int, fileName string) string {
	return fmt.Sprintf("%s/%s/%d/%s", basePath, productsDirectory, productId, fileName)
}
