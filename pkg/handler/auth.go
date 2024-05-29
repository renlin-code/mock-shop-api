package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

// @Summary User Sign Up
// @Tags User Authorization
// @Description Create a user account. With this account the user can place orders. If the request is successful, the service sends an e-mail to the specified email address with an email confirmation token as a URL param "confToken". For example: https://store.com/confirm-email?confToken=eyJhbGciOiJIU1iIR5csdDIkwErXVCJ9. This token is required to confirm specified user email.
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.SignUpInput true "Account info"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-up [post]
func (h *Handler) userSignUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
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

// @Summary User Confirm Email
// @Tags User Authorization
// @Description Confirm the specified email when creating a user account and add a password for the account. If the request is successful, the user account is created and the user can log into it.
// @ID confirm-email
// @Accept json
// @Produce json
// @Param input body domain.ConfirmEmailInput true "Account info"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/confirm-email [post]
func (h *Handler) userConfirmEmail(c *gin.Context) {
	var input domain.ConfirmEmailInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
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

// @Summary User Sign In
// @Tags User Authorization
// @Description Log into an existing user account. If the request is successful, the service returns an authorization token.
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "Account access"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-in [post]
func (h *Handler) userSignIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
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

// @Summary User Recovery Password
// @Tags User Authorization
// @Description Recovery password. If the request is successful, the service sends an e-mail to the account email address with an email confirmation token as a URL param "confToken". For example: https://client.com/password-recovery?confToken=eyJhbGciOiJIU1iIR5csdDIkwErXVCJ9. This token is required to set a new password.
// @ID recovery-password
// @Accept json
// @Produce json
// @Param input body domain.RecoveryPasswordInput true "Account email"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/password-recovery [post]
func (h *Handler) recoveryUserPassword(c *gin.Context) {
	var input domain.RecoveryPasswordInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
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

// @Summary User Update Password
// @Tags User Authorization
// @Description Set new password. If the request is successful, the account password is changed for the specified new password.
// @ID update-password
// @Accept json
// @Produce json
// @Param input body domain.UpdatePasswordInput true "Account new password"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/password-update [put]
func (h *Handler) updateUserPassword(c *gin.Context) {
	var input domain.UpdatePasswordInput
	if err := c.BindJSON(&input); err != nil {
		Fail(c, bindJSONErrorText, http.StatusBadRequest)
		return
	}
	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.services.Authorization.UpdatePassword(input.Token, input.Password)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}
	OK(c)
}
