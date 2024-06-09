package storage

type Config struct {
	MediaBaseUrl string
}

type FileSystemStorage struct {
	config Config
}

func NewFileSystemStorage(cfg Config) *FileSystemStorage {
	return &FileSystemStorage{config: cfg}
}

const (
	basePath            = "/data"
	usersDirectory      = "users"
	categoriesDirectory = "categories"
	productsDirectory   = "products"
)
