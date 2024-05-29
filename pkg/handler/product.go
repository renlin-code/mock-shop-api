package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

// @Summary Get Products
// @Tags Products
// @Description Get all products.
// @ID get-products
// @Accept json
// @Produce json
// @Param page query string false "Pagination: page number"
// @Param pageSize query string false "Pagination: amount of items per page"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/products [get]
func (h *Handler) getAllProducts(c *gin.Context) {
	var params domain.PaginationParams
	if err := c.BindQuery(&params); err != nil {
		Fail(c, bindPaginationParamsErrorText, http.StatusBadRequest)
		return
	}
	if err := params.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := h.services.Product.GetAll(computePaginationParams(params))
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, categories)
}

// @Summary Get Product By Id
// @Tags Products
// @Description Get product by id.
// @ID get-product-by-id
// @Accept json
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/products/{id} [get]
func (h *Handler) getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	category, err := h.services.Product.GetById(id)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, category)
}

// @Summary Create Product
// @Security ApiKeyAuth
// @Tags Admin
// @Description Create a new product.
// @ID create-product
// @Accept  multipart/form-data
// @Produce json
// @Param category_id formData int true "Category id"
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param price formData number true "Product actual price"
// @Param undiscounted_price formData number true "Product price without any discount"
// @Param stock formData int true "Product stock"
// @Param available formData boolean true "Product is available"
// @Param image_file formData file true "Product image"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admin/products [post]
func (h *Handler) adminCreateProduct(c *gin.Context) {
	r := c.Request
	var input domain.CreateProductInput

	catIdInt, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		Fail(c, "invalid category_id", http.StatusBadRequest)
		return
	}
	input.CategoryId = catIdInt

	input.Name = r.FormValue("name")
	input.Description = r.FormValue("description")

	available := r.FormValue("available")
	input.Available = available == "true"

	stock, err := strconv.Atoi(r.FormValue("stock"))
	if err != nil {
		Fail(c, "invalid stock value", http.StatusBadRequest)
		return
	}
	input.Stock = stock

	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		Fail(c, "invalid price value", http.StatusBadRequest)
		return
	}
	input.Price = float32(price)

	price, err = strconv.ParseFloat(r.FormValue("undiscounted_price"), 32)
	if err != nil {
		Fail(c, "invalid undiscounted_price value", http.StatusBadRequest)
		return
	}
	input.UndiscountedPrice = float32(price)

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		if err == http.ErrMissingFile {
			Fail(c, "image_file: can not be blank", http.StatusBadRequest)
			return
		}
		FailAndHandleErr(c, err)
		return
	}
	defer file.Close()

	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.services.Product.CreateProduct(input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OKId(c, id)
}

// @Summary Update Product
// @Security ApiKeyAuth
// @Tags Admin
// @Description Update product.
// @ID update-product
// @Accept  multipart/form-data
// @Produce json
// @Param id path int true "Product id"
// @Param category_id formData int false "Category id"
// @Param name formData string false "Product name"
// @Param description formData string false "Product description"
// @Param price formData number false "Product actual price"
// @Param undiscounted_price formData number false "Product price without any discount"
// @Param stock formData int false "Product stock"
// @Param available formData boolean false "Product is available"
// @Param image_file formData file false "Product image"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admin/products/{id} [put]
func (h *Handler) adminUpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	r := c.Request
	var input domain.UpdateProductInput

	categoryId := r.FormValue("category_id")
	if categoryId != "" {
		catIdInt, err := strconv.Atoi(categoryId)
		if err != nil {
			Fail(c, "invalid category id", http.StatusBadRequest)
		}
		input.CategoryId = &catIdInt
	}

	name := r.FormValue("name")
	if name != "" {
		input.Name = &name
	}

	description := r.FormValue("description")
	if description != "" {
		input.Description = &description
	}

	availableField := r.FormValue("available")
	available := availableField != "false"
	input.Available = &available

	stockField := r.FormValue("stock")
	if stockField != "" {
		stock, err := strconv.Atoi(stockField)
		if err != nil {
			Fail(c, "invalid stock value", http.StatusBadRequest)
		}
		input.Stock = &stock
	}

	price := r.FormValue("price")
	if price != "" {
		price64, err := strconv.ParseFloat(price, 32)
		if err != nil {
			Fail(c, "invalid price value", http.StatusBadRequest)
		}

		price := float32(price64)
		input.Price = &price
	}

	undiscountedPrice := r.FormValue("undiscounted_price")
	if undiscountedPrice != "" {
		price64, err := strconv.ParseFloat(undiscountedPrice, 32)
		if err != nil {
			Fail(c, "invalid price value", http.StatusBadRequest)
		}

		price := float32(price64)
		input.UndiscountedPrice = &price
	}

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		if err != http.ErrMissingFile {
			Fail(c, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}
	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.services.Product.UpdateProduct(id, input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}
