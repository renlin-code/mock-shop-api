package storage

type Config struct {
	BaseUrl string
}

type FileSystemStorage struct {
	config Config
}

func NewFileSystemStorage(cfg Config) *FileSystemStorage {
	return &FileSystemStorage{config: cfg}
}
