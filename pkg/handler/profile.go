package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

// @Summary Get User Account
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Get user account.
// @ID get-user-account
// @Accept json
// @Produce json
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/ [get]
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

// @Summary Update User Account
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Update user account info.
// @ID update-user-account
// @Accept  multipart/form-data
// @Produce json
// @Param name formData string false "User name"
// @Param profile_image_file formData file false "User profile image"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/ [put]
func (h *Handler) updateUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	r := c.Request
	var input domain.UpdateProfileInput

	name := r.FormValue("name")
	if name != "" {
		input.Name = &name
	}

	file, handler, err := r.FormFile("profile_image_file")
	if err != nil {
		if err != http.ErrMissingFile {
			Fail(c, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}

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

// @Summary Create Order
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Create a new order. If the request is successful, a new order is added to the user's list of orders and the stock of products in the catalog is updated.
// @ID create-order
// @Accept json
// @Produce json
// @Param input body domain.CreateOrderInput true "Order info"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/orders [post]
func (h *Handler) userCreateOrder(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input domain.CreateOrderInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
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

// @Summary Get Orders
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Get all user's orders.
// @ID get-orders
// @Accept json
// @Produce json
// @Param page query string false "Pagination: page number"
// @Param pageSize query string false "Pagination: amount of items per page"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/orders [get]
func (h *Handler) userGetAllOrder(c *gin.Context) {
	var params domain.PaginationParams
	if err := c.BindQuery(&params); err != nil {
		Fail(c, bindPaginationParamsErrorText, http.StatusBadRequest)
		return
	}
	if err := params.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	limit, offset := computePaginationParams(params)
	orders, err := h.services.Profile.GetAllOrders(userId, limit, offset)

	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	Response(c, orders)
}

// @Summary Get Order By Id
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Get user's order by id.
// @ID get-order-by-id
// @Accept json
// @Produce json
// @Param id path int true "Order id"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/orders/{id} [get]
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

// @Summary Delete User Profile
// @Security ApiKeyAuth
// @Tags User Profile
// @Description Delete user profile.
// @ID delete-user-account
// @Accept json
// @Produce json
// @Param input body domain.DeleteProfileInput true "Account password"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /profile/ [delete]
func (h *Handler) deleteUserProfile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input domain.DeleteProfileInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.services.Profile.DeleteProfile(userId, input.Password)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}

func (h *Handler) mediaGetUserImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}
	fileName := c.Param("file-name")

	filePath := h.services.Profile.GetFilePath(id, fileName)

	c.File(filePath)
}
