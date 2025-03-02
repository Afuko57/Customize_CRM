package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"customize_crm/model"
	"customize_crm/service"
	"customize_crm/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserController struct {
	userService *service.UserService
}

type CreateUserRequest struct {
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	RoleID     uuid.UUID `json:"role_id"`
	Department *string   `json:"department,omitempty"`
	IsActive   bool      `json:"is_active"`
}

type UpdateUserRequest struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Department *string   `json:"department,omitempty"`
	RoleID     uuid.UUID `json:"role_id"`
	IsActive   bool      `json:"is_active"`
}

type DeleteUsersRequest struct {
	IDs []uuid.UUID `json:"ids"`
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the current authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.User
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/users/me [get]
func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	user, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateCurrentUser godoc
// @Summary Update current user
// @Description Update the current authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} model.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/users/me [patch]
func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Department = req.Department

	if err := c.userService.Update(r.Context(), user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.User
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/users [get]
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAll(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "New user data"
// @Success 201 {object} model.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Required fields are missing")
		return
	}

	user := &model.User{
		Username:   req.Username,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		RoleID:     req.RoleID,
		Department: req.Department,
		IsActive:   req.IsActive,
	}

	if err := c.userService.Create(r.Context(), user, req.Password); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "username") {
				utils.RespondWithError(w, http.StatusBadRequest, "Username already exists")
			} else if strings.Contains(err.Error(), "email") {
				utils.RespondWithError(w, http.StatusBadRequest, "Email already exists")
			} else {
				utils.RespondWithError(w, http.StatusBadRequest, "User already exists")
			}
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	user, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user by ID (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} model.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/users/{id} [patch]
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.userService.GetByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Department = req.Department
	user.RoleID = req.RoleID
	user.IsActive = req.IsActive

	if err := c.userService.Update(r.Context(), user); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUsers godoc
// @Summary Delete multiple users
// @Description Delete multiple users by IDs (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DeleteUsersRequest true "User IDs to delete"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/users [delete]
func (c *UserController) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var req DeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.IDs) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No user IDs provided")
		return
	}

	if err := c.userService.Delete(r.Context(), req.IDs); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting users")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Users deleted successfully"})
}
