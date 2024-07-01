package model

import "mime/multipart"

type CreateServiceRequest struct {
	Name           string                `json:"name" validate:"required"`
	Description    string                `json:"description" validate:"required"`
	DurationMinute int                   `json:"durationMinute" validate:"required"`
	Price          int                   `json:"price" validate:"required"`
	ImageFile      *multipart.FileHeader `json:"-"`
}

type CreateServiceResponse struct {
	ID string `json:"id"`
}

type FindServiceResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Price          int    `json:"price"`
	DurationMinute int    `json:"durationMinute"`
	ImageURL       string `json:"imageUrl"`
}
