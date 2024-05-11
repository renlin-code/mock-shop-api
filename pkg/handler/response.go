package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/errors_handler"
	"github.com/sirupsen/logrus"
)

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func newResponse(success bool, message string, data interface{}) *response {
	return &response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// OK send a successful response with body.
func Response(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, newResponse(true, "", data))
}

// OK send a successful response without body.
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, newResponse(true, "", nil))
}

// OKId send a successful response without only an id in the body.
func OKId(c *gin.Context, id int) {
	c.JSON(http.StatusOK, newResponse(true, "", map[string]int{
		"id": id,
	}))
}

// OKToken send a successful response without only a token in the body.
func OKToken(c *gin.Context, token string) {
	c.JSON(http.StatusOK, newResponse(true, "", map[string]string{
		"token": token,
	}))
}

// Fail send a failed response.
func Fail(c *gin.Context, message string, statusCode int) {
	c.AbortWithStatusJSON(statusCode, newResponse(false, message, nil))
}

// Fail send a failed response with appError.
func AbortWithAppErr(c *gin.Context, err error) {
	appError := new(errors_handler.AppError)
	if errors.As(err, &appError) { // client error
		c.AbortWithStatusJSON(errTypeStatusCode(appError.Type()), newResponse(false, err.Error(), nil))
		return
	}

	logrus.Errorf("server error: %s", err.Error())
	Fail(c, err.Error(), http.StatusInternalServerError)
}

// errTypeStatusCode return a http status code depending on an appError error type.
func errTypeStatusCode(errType errors_handler.Type) int {
	switch errType {
	case errors_handler.TypeBadRequest:
		return http.StatusBadRequest
	case errors_handler.TypeNotFound:
		return http.StatusNotFound
	case errors_handler.TypeForbidden:
		return http.StatusForbidden
	case errors_handler.TypeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusBadRequest
	}
}
