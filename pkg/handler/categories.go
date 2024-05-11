package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllCategories(c *gin.Context) {
	categories, err := h.services.Category.GetAll()
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, categories)
}

func (h *Handler) getCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, "invalid id param", http.StatusBadRequest)
		return
	}

	category, err := h.services.Category.GetById(id)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, category)
}

func (h *Handler) getCategoryProducts(c *gin.Context) {
	catId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, "invalid id param", http.StatusBadRequest)
		return
	}

	products, err := h.services.Category.GetProducts(catId)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, products)
}
