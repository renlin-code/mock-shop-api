package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

func (h *Handler) getUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	user, err := h.services.Profile.GetProfile(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input domain.UpdateProfileInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateInput(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Profile.UpdateProfile(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) recoveryUserPassword(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	err = h.services.Profile.RecoveryPassword(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}

func (h *Handler) updateUserPassword(c *gin.Context) {
	var input domain.UpdatePasswordInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateInput(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Profile.UpdatePassword(input.Token, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) userCreateOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input domain.CreateOrderInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateInput(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.Sort()
	id, err := h.services.Profile.CreateOrder(userId, input.Products)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) userGetAllOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	orders, err := h.services.Profile.GetAllOrders(userId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) userGetOrderById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	order, err := h.services.Profile.GetOrderById(userId, orderId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, order)
}
