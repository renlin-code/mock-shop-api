package errors_handler

import "errors"

// Type describes the kind of error in the app.
type Type string

const (
	// TypeBadRequest is used for HTTP 400-like errors.
	TypeBadRequest Type = "bad_request_error"
	// TypeNotFound is used for HTTP 404-like errors.
	TypeNotFound Type = "not_found_error"
	// TypeForbidden is used for HTTP 403-like errors.
	TypeForbidden Type = "forbidden"
	// TypeUnauthorized is used for HTTP 401-like errors.
	TypeUnauthorized Type = "unauthorized"

	// TypeNoRows is used for DB errors when query response is empty.
	TypeNoRows Type = "no_rows"
	// TypeAlreadyExists is used for DB errors when is violated unique constrain.
	TypeAlreadyExists Type = "already_exists"
	// TypeForeignKeyViolation is used for DB errors when is violated a foreign key.
	TypeForeignKeyViolation Type = "foreign_key_violation"
	// TypeConstrainViolation is used for DB errors when is violated a standard constrain.
	TypeConstrainViolation Type = "constrain_violation"

	// TypeStorageError is used for storage errors when there is an error creating file.
	TypeStorageError Type = "storage_error"
)

// AppError is an implementation of error with types to
// differentiate client and server errors.
type AppError struct {
	text    string
	errType Type
}

func (e AppError) Error() string {
	return e.text
}

// Type returns the type of the error.
func (e AppError) Type() Type {
	return e.errType
}

// ErrorIsType compare an error with an AppError type and returns a boolean.
func ErrorIsType(err error, errType Type) bool {
	appError := new(AppError)
	return errors.As(err, &appError) && appError.Type() == errType
}

// BadRequest returns an AppError with a TypeBadRequest type.
func BadRequest(text string) error {
	return &AppError{
		text:    text,
		errType: TypeBadRequest,
	}
}

// NotFound returns an AppError with a TypeNotFound type.
func NotFound(entity string) error {
	return &AppError{
		text:    entity + " not found",
		errType: TypeNotFound,
	}
}

// Forbidden returns an AppError with a TypeForbidden type.
func Forbidden(text string) error {
	return &AppError{
		text:    text,
		errType: TypeForbidden,
	}
}

// Unauthorized returns an AppError with a TypeUnauthorized type.
func Unauthorized(text string) error {
	return &AppError{
		text:    text,
		errType: TypeUnauthorized,
	}
}

// NoRows returns an AppError with a TypeNoRows type.
func NoRows() error {
	return &AppError{
		text:    "no rows",
		errType: TypeNoRows,
	}
}

// AlreadyExists returns an AppError with a TypeAlreadyExists type.
func AlreadyExists(entity string) error {
	return &AppError{
		text:    entity + " already exists",
		errType: TypeAlreadyExists,
	}
}

// ForeingKeyViolation returns an AppError with a TypeForeignKeyViolation type.
func ForeignKeyViolation() error {
	return &AppError{
		text:    "foreign key violated",
		errType: TypeForeignKeyViolation,
	}
}

// ConstrainViolation returns an AppError with a TypeConstrainViolation type.
func ConstrainViolation(constrainName string) error {
	return &AppError{
		text:    constrainName + " constrain violated",
		errType: TypeConstrainViolation,
	}
}

// StorageError returns an AppError with a TypeStorageError type.
func StorageError(text string) error {
	return &AppError{
		text:    text,
		errType: TypeStorageError,
	}
}
