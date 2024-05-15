package storage

import (
	"mime/multipart"
)

type Profile interface {
	UploadProfileImage(userId int, handler *multipart.FileHeader, file multipart.File) (string, error)
}

type Storage struct {
	Profile
}

func NewStorage(fs *FileSystemStorage) *Storage {
	return &Storage{
		Profile: newProfileLocal(fs),
	}
}
