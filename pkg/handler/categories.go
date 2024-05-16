package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
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
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
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
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	products, err := h.services.Category.GetProducts(catId)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, products)
}

func (h *Handler) adminCreateCategory(c *gin.Context) {
	r := c.Request
	var input domain.CreateCategoryInput
	input.Name = r.FormValue("name")
	input.Description = r.FormValue("description")

	available := r.FormValue("available")
	input.Available = available == "true"

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		Fail(c, "image_file: cannot be blank", http.StatusBadRequest)
		return
	}
	defer file.Close()

	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.services.Category.CreateCategory(input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OKId(c, id)
}

func (h *Handler) adminUpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	r := c.Request
	var input domain.UpdateCategoryInput

	name := r.FormValue("name")
	description := r.FormValue("description")
	availableValue := r.FormValue("available")
	available := availableValue != "false"

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		if err != http.ErrMissingFile {
			Fail(c, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}
	input.Name = &name
	input.Description = &description
	input.Available = &available
	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.services.Category.UpdateCategory(id, input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}
