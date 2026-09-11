package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/modules/auth"
	"github.com/kansaok/go-boilerplate/internal/util"
	"github.com/kansaok/go-boilerplate/internal/util/validators"
)

var validate = validators.Validate

// @Summary Register new user
// @Description Register a new user with email, password, and profile details
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param registerRequest body auth.RegisterRequestSchema true "Register Request Body"
// @Success 200 {object} auth.RegisterResponseSchema "Register data berhasil"
// @Failure 400 {object} util.APIResponse "Invalid payload"
// @Failure 422 {object} util.APIResponse "Error in validation"
// @Failure 500 {object} util.APIResponse "Register failed"
// @Router /v1/auth/register [post]
func Register(c *gin.Context) {
	var req auth.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if util.IsBodyTooLargeError(err) {
			util.RespondWithError(c, util.CodeRequestEntityTooLarge, "Request body terlalu besar", nil)
			return
		}
		util.RespondWithError(c, util.CodeBadRequest, util.MESSAGES["INVALID_PAYLOAD"], nil)
		return
	}

	validationErrors := validators.ValidateAndMapErrors(validate, req)
	if len(validationErrors) > 0 {
		util.RespondWithError(c, util.CodeValidationError, util.MESSAGES["ERROR_VALIDATION"], validationErrors)
		return
	}

	data, err := auth.RegisterUser(c.Request.Context(), req)
	if err != nil {
		util.RespondWithError(c, util.CodeUnknownError, util.MESSAGES["REGISTER_FAILED"], nil)
		return
	}

	util.RespondWithSuccess(c, data, util.MESSAGES["REGISTER_SUCCESS"])
}

// @Summary Login User
// @Description Login user with email and password
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param loginRequest body auth.LoginRequestSchema true "Login Request Body"
// @Success 200 {object} auth.LoginResponseSchema "Login berhasil"
// @Failure 400 {object} util.APIResponse "Invalid payload"
// @Failure 422 {object} util.APIResponse "Error in validation"
// @Failure 500 {object} util.APIResponse "Login failed"
// @Router /v1/auth/login [post]
func Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if util.IsBodyTooLargeError(err) {
			util.RespondWithError(c, util.CodeRequestEntityTooLarge, "Request body terlalu besar", nil)
			return
		}
		util.RespondWithError(c, util.CodeBadRequest, util.MESSAGES["INVALID_PAYLOAD"], nil)
		return
	}

	jwtConfig := config.LoadJWTConfig()

	token, err := auth.AuthenticateUser(req.Email, req.Password, jwtConfig)
	if err != nil {
		util.RespondWithError(c, util.CodeValidationError, util.MESSAGES["LOGIN_FAILED"], nil)
		return
	}
	loginResponse := auth.LoginResponse{
		Token: token,
	}

	util.RespondWithSuccess(c, loginResponse, util.MESSAGES["LOGIN_SUCCESS"])
}

// @Summary Get current user
// @Description Get the authenticated user's identity from the JWT
// @Tags Auth
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} util.APIResponse "User data berhasil"
// @Failure 401 {object} util.APIResponse "Unauthorized"
// @Router /v1/me [get]
func Me(c *gin.Context) {
	email, _ := c.Get("userEmail")
	subject, _ := c.Get("userSubject")

	util.RespondWithSuccess(c, gin.H{
		"email":   email,
		"subject": subject,
	}, util.MESSAGES["GET_SUCCESS"])
}
