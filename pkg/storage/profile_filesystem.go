package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
)

type ProfileFileSystem struct {
	FileSystem *FileSystemStorage
}

func newProfileFileSystem(fs *FileSystemStorage) *ProfileFileSystem {
	return &ProfileFileSystem{FileSystem: fs}
}

func (s *ProfileFileSystem) UploadProfileImage(userId int, handler *multipart.FileHeader, file multipart.File) (string, error) {
	userDir := fmt.Sprintf("/data/users/%d/", userId)

	err := os.RemoveAll("." + userDir)
	if err != nil {
		return "", errors_handler.StorageError("error deleting directory")
	}

	err = os.MkdirAll("."+userDir, os.ModePerm)
	if err != nil {
		return "", errors_handler.StorageError("error creating directory")
	}

	path := userDir + handler.Filename
	f, err := os.Create("." + path)
	if err != nil {
		return "", errors_handler.StorageError("error creating file")
	}
	defer f.Close()

	io.Copy(f, file)

	return s.FileSystem.config.BaseUrl + path, nil
}
