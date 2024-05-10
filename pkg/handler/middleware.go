package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "UserId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		Fail(c, "no authorization header provided", http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		Fail(c, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	userId, err := h.services.Authorization.ParseAuthToken(headerParts[1])

	if err != nil {
		Fail(c, err.Error(), http.StatusUnauthorized)
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		Fail(c, "user not found", http.StatusNotFound)
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		Fail(c, "invalid type for id", http.StatusBadRequest)
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
