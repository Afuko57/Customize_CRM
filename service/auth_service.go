package service

import (
	"context"
	"customize_crm/model"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userService *UserService
	jwtSecret   string
	accessExp   time.Duration
	refreshExp  time.Duration
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

func NewAuthService(userService *UserService) *AuthService {
	accessMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRY_MINUTES"))
	if accessMinutes == 0 {
		accessMinutes = 15
	}

	refreshDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY_DAYS"))
	if refreshDays == 0 {
		refreshDays = 7
	}

	return &AuthService{
		userService: userService,
		jwtSecret:   os.Getenv("JWT_SECRET"),
		accessExp:   time.Duration(accessMinutes) * time.Minute,
		refreshExp:  time.Duration(refreshDays) * 24 * time.Hour,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*TokenDetails, *model.User, error) {
	user, err := s.userService.Authenticate(ctx, username, password)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.CreateTokens(user.ID.String())
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration time")
	}
	if time.Unix(int64(exp), 0).Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID")
	}

	return s.CreateTokens(userID)
}

func (s *AuthService) CreateTokens(userID string) (*TokenDetails, error) {
	td := &TokenDetails{
		AccessUUID:  uuid.New().String(),
		RefreshUUID: uuid.New().String(),
		AtExpires:   time.Now().Add(s.accessExp).Unix(),
		RtExpires:   time.Now().Add(s.refreshExp).Unix(),
	}

	atClaims := jwt.MapClaims{
		"sub": userID,
		"exp": td.AtExpires,
		"jti": td.AccessUUID,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken

	rtClaims := jwt.MapClaims{
		"sub": userID,
		"exp": td.RtExpires,
		"jti": td.RefreshUUID,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = refreshToken

	return td, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	return nil
}
