package middleware

import (
	"context"
	"net/http"
	"strings"

	"customize_crm/service"
	"customize_crm/utils"

	"github.com/google/uuid"
)

type AuthMiddleware struct {
	userService *service.UserService
}

func NewAuthMiddleware(userService *service.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		user, err := m.userService.GetByID(r.Context(), userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		role, err := m.userService.GetRoleByID(r.Context(), user.RoleID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching user role")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "role", role.Name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "Admin" {
			utils.RespondWithError(w, http.StatusForbidden, "Admin permission required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
