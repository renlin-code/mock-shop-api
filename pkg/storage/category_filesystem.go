package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type CategoryFileSystem struct {
	FileSystem *FileSystemStorage
}

func newCategoryFileSystem(fs *FileSystemStorage) *CategoryFileSystem {
	return &CategoryFileSystem{FileSystem: fs}
}

func (s *CategoryFileSystem) UploadCategoryImage(categoryId int, handler *multipart.FileHeader, file multipart.File) (string, error) {
	categoryDir := fmt.Sprintf("%s/%s/%d/", basePath, categoriesDirectory, categoryId)

	err := os.RemoveAll("." + categoryDir)
	if err != nil {
		return "", errors_handler.StorageError("error deleting directory")
	}

	err = os.MkdirAll("."+categoryDir, os.ModePerm)
	if err != nil {
		return "", errors_handler.StorageError("error creating directory")
	}

	path := categoryDir + handler.Filename
	f, err := os.Create("." + path)
	if err != nil {
		return "", errors_handler.StorageError("error creating file")
	}
	defer f.Close()

	io.Copy(f, file)

	return fmt.Sprintf("%s/%s/%d/%s", s.FileSystem.config.MediaBaseUrl, categoriesDirectory, categoryId, handler.Filename), nil
}

func (s *CategoryFileSystem) GetFilePath(categoryId int, fileName string) string {
	return fmt.Sprintf("%s/%s/%d/%s", basePath, categoriesDirectory, categoryId, fileName)
}
