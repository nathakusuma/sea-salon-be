package service

import (
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type IReviewService interface {
	Create(req model.CreateReviewRequest, claims *jwt.Claims) response.Response
	FindByLazyLoad(req model.FindReviewsLazyLoadRequest) response.Response
}

type reviewService struct {
	r  repository.IReviewRepository
	ur repository.IUserRepository
}

func NewReviewService(r repository.IReviewRepository, ur repository.IUserRepository) IReviewService {
	return &reviewService{r: r, ur: ur}
}

func (s *reviewService) Create(req model.CreateReviewRequest, claims *jwt.Claims) response.Response {
	userID, err := ulid.Parse(claims.Subject)
	if err != nil {
		return response.New(400, "Fail to parse userID", err.Error())
	}

	review := entity.Review{
		Model:      gorm.Model{},
		ID:         ulid.Make(),
		UserID:     userID,
		StarRating: req.StarRating,
		Comment:    req.Comment,
	}

	id, err := s.r.Create(&review)
	if err != nil {
		return response.New(500, "Fail to create review", nil)
	}

	return response.New(201, "Review created", model.CreateReviewResponse{ID: id.String()})
}

func (s *reviewService) FindByLazyLoad(req model.FindReviewsLazyLoadRequest) response.Response {
	var idPivot ulid.ULID
	if req.Action == "prev" || req.Action == "next" {
		var err error
		idPivot, err = ulid.Parse(req.IdPivot)
		if err != nil {
			return response.New(400, "Fail to parse id", err.Error())
		}
	}

	if req.Limit > 20 {
		return response.New(400, "You request too much data", nil)
	}

	reviews, err := s.r.FindByLazyLoad(idPivot, req.Action, req.Limit)
	if err != nil {
		return response.New(500, "Fail to fetch reviews", nil)
	}

	reviewsRes := make([]model.FindReviewResponse, len(reviews))
	for i, review := range reviews {
		reviewsRes[i] = model.FindReviewResponse{
			ID:           review.ID.String(),
			CustomerName: review.User.FullName,
			StarRating:   review.StarRating,
			Comment:      review.Comment,
			CreatedAt:    review.CreatedAt.Format(time.RFC3339),
		}
	}

	return response.New(200, "Reviews fetched", reviewsRes)
}
