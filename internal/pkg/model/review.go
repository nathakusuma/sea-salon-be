package model

type CreateReviewRequest struct {
	CustomerName string `json:"customerName" validate:"required,max=50,min=1"`
	StarRating   int    `json:"starRating" validate:"required,min=1,max=5"`
	Comment      string `json:"comment" validate:"required,max=255,min=1"`
}

type CreateReviewResponse struct {
	ID string `json:"id"`
}

type FindReviewsLazyLoadRequest struct {
	Limit   int    `form:"limit" validate:"required,min=1,max=30"`
	IdPivot string `form:"idPivot"`
	Action  string `form:"action" validate:"oneof=prev next top ''"`
}

type FindReviewResponse struct {
	ID           string `json:"id"`
	CustomerName string `json:"customerName"`
	StarRating   int    `json:"starRating"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"createdAt"`
}

type FindReviewsLazyLoadResponse struct {
	FindReviewResponse `json:"reviews"`
	HasNext            bool `json:"hasNext"`
	HasPrev            bool `json:"hasPrev"`
}
