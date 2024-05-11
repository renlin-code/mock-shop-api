package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllProducts(c *gin.Context) {
	categories, err := h.services.Product.GetAll()
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, categories)
}

func (h *Handler) getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, "invalid id param", http.StatusBadRequest)
		return
	}

	category, err := h.services.Product.GetById(id)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, category)
}
