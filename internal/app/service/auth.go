package service

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/bcrypt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(req model.RegisterRequest) response.Response
	Login(req model.LoginRequest) response.Response
}

type authService struct {
	ur      repository.IUserRepository
	jwtAuth jwt.IJWT
}

func NewAuthService(ur repository.IUserRepository, jwtAuth jwt.IJWT) IAuthService {
	return &authService{ur: ur, jwtAuth: jwtAuth}
}

func (s *authService) Register(req model.RegisterRequest) response.Response {
	hashedPassword, err := bcrypt.Hash(req.Password)
	if err != nil {
		return response.New(500, "Fail to hash password", nil)
	}

	user := entity.User{
		Model:       gorm.Model{},
		ID:          ulid.Make(),
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		IsAdmin:     false,
	}

	id, err := s.ur.Create(&user)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "23505" {
			return response.New(409, "Email already exists", nil)
		}
		return response.New(500, "Fail to create user", nil)
	}

	return response.New(201, "Successfully created user", model.RegisterResponse{ID: id.String()})
}

func (s *authService) Login(req model.LoginRequest) response.Response {
	user, err := s.ur.FindByEmail(req.Email)
	if err != nil {
		return response.New(500, "Fail to find user", nil)
	}

	if err := bcrypt.ValidateHash(req.Password, user.Password); err != nil {
		return response.New(401, "Invalid email or password", nil)
	}

	token, err := s.jwtAuth.Create(&user)
	if err != nil {
		return response.New(500, "Fail to create token", nil)
	}

	return response.New(200, "Successfully logged in", model.LoginResponse{
		Token:    token,
		FullName: user.FullName,
	})
}
