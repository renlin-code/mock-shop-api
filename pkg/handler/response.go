package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// OKId send a successful response without only a token in the body.
func OKToken(c *gin.Context, token string) {
	c.JSON(http.StatusOK, newResponse(true, "", map[string]string{
		"token": token,
	}))
}

// OK send a failed response.
func Fail(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, newResponse(false, message, nil))
}
