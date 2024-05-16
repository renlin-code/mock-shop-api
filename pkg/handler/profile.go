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
		FailAndHandleErr(c, err)
		return
	}
	Response(c, user)
}

func (h *Handler) updateUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	r := c.Request
	var input domain.UpdateProfileInput

	name := r.FormValue("name")
	file, handler, err := r.FormFile("profile_image_file")
	if err != nil {
		if err != http.ErrMissingFile {
			Fail(c, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}

	input.Name = &name
	input.ProfileImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.Profile.UpdateProfile(userId, input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	OK(c)
}

func (h *Handler) userCreateOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input domain.CreateOrderInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	input.Sort()
	id, err := h.services.Profile.CreateOrder(userId, input.Products)

	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	OKId(c, id)
}

func (h *Handler) userGetAllOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	orders, err := h.services.Profile.GetAllOrders(userId)

	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	Response(c, orders)
}

func (h *Handler) userGetOrderById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	order, err := h.services.Profile.GetOrderById(userId, orderId)

	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	Response(c, order)
}
