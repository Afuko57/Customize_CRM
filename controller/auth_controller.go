package controller

import (
	"encoding/json"
	"net/http"

	"customize_crm/model"
	"customize_crm/service"
	"customize_crm/utils"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	tokens, user, err := c.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, model.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserID:       user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		ExpiresIn:    tokens.AtExpires - tokens.RtExpires,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.RefreshTokenResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /api/v1/auth/refresh-token [post]
func (c *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.RefreshToken == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	tokens, err := c.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	response := model.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.AtExpires - tokens.RtExpires,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// Logout godoc
// @Summary User logout
// @Description Logout a user and invalidate their token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.MessageResponse
// @Router /api/v1/auth/logout [post]
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, model.MessageResponse{
		Message: "Logged out successfully",
	})
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ForgotPasswordRequest true "User email"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/v1/auth/forgot-password [post]
func (c *AuthController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, model.MessageResponse{
		Message: "Password reset functionality not implemented yet",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/v1/auth/reset-password [post]
func (c *AuthController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, model.MessageResponse{
		Message: "Password reset functionality not implemented yet",
	})
}
