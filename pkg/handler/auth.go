package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

func (h *Handler) userSignUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.services.UserSignUp(input.Name, input.Email)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}

func (h *Handler) userConfirmEmail(c *gin.Context) {
	var input domain.ConfirmEmailInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.services.CreateUser(input.Token, input.Password)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OKId(c, id)
}

func (h *Handler) userSignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.services.Authorization.GenerateAuthToken(input.Email, input.Password)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OKToken(c, token)
}

func (h *Handler) recoveryUserPassword(c *gin.Context) {
	var input domain.RecoveryPasswordInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.services.Authorization.RecoveryPassword(input.Email)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}

func (h *Handler) updateUserPassword(c *gin.Context) {
	var input domain.UpdatePasswordInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.services.Authorization.UpdatePassword(input.Token, input.Password)
	if err != nil {
		Fail(c, err.Error(), http.StatusInternalServerError)
		return
	}
	OK(c)
}
