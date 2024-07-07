package model

import "mime/multipart"

type CreateBranchRequest struct {
	Name           string                `json:"name" validate:"required,max=50,min=1"`
	Address        string                `json:"address" validate:"required,max=255,min=1"`
	MapsURL        string                `json:"mapsUrl" validate:"required"`
	Phone          string                `json:"phone" validate:"required,max=15,min=1"`
	OpenTime       string                `json:"openTime" validate:"required"`
	CloseTime      string                `json:"closeTime" validate:"required"`
	TimeZoneName   string                `json:"timeZoneName" validate:"required"`
	TimeZoneOffset string                `json:"timeZoneOffset" validate:"required"`
	ImageFile      *multipart.FileHeader `json:"-"`
}

type CreateBranchResponse struct {
	ID string `json:"id"`
}

type SetServicesToBranchRequest struct {
	ServiceIDs []string `json:"serviceIds" validate:"required"`
	BranchID   string   `json:"branchId" validate:"required"`
}

type FindBranchResponse struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Address      string                `json:"address"`
	MapsURL      string                `json:"mapsUrl"`
	Phone        string                `json:"phone"`
	OpenTime     string                `json:"openTime"`
	CloseTime    string                `json:"closeTime"`
	TimeZoneName string                `json:"timeZoneName"`
	Services     []FindServiceResponse `json:"services"`
	ImageURL     string                `json:"imageUrl"`
}
