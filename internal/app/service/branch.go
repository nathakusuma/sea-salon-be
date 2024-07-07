package service

import (
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/gcloud"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/helper"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"log"
	"time"
)

type IBranchService interface {
	Create(req *model.CreateBranchRequest) response.Response
	FindAll() response.Response
	SetServices(req model.SetServicesToBranchRequest) response.Response
}

type branchService struct {
	r        repository.IBranchRepository
	uploader *gcloud.FileUploaderClient
}

func NewBranchService(r repository.IBranchRepository, uploader *gcloud.FileUploaderClient) IBranchService {
	return &branchService{r: r, uploader: uploader}
}

func (s *branchService) Create(req *model.CreateBranchRequest) response.Response {
	imageFile, err := req.ImageFile.Open()
	if err != nil {
		return response.New(400, "Fail to open image", err.Error())
	}

	openTime, err1 := time.Parse("15:04", req.OpenTime)
	closeTime, err2 := time.Parse("15:04", req.CloseTime)
	if err1 != nil || err2 != nil {
		return response.New(400, "Fail to parse time", nil)
	}

	if openTime.After(closeTime) {
		return response.New(400, "Open time must be before close time", nil)
	}

	offset, err := helper.TimeZoneOffsetToSeconds(req.TimeZoneOffset)
	if err != nil {
		return response.New(400, "Fail to parse time zone offset", err.Error())
	}

	branch := entity.Branch{
		Model:          gorm.Model{},
		ID:             ulid.Make(),
		Name:           req.Name,
		Address:        req.Address,
		MapsURL:        req.MapsURL,
		Phone:          req.Phone,
		OpenTime:       openTime,
		CloseTime:      closeTime,
		TimeZoneName:   req.TimeZoneName,
		TimeZoneOffset: offset,
		Services:       nil,
	}

	if err = s.uploader.UploadFile(imageFile, "img/branches/"+branch.ID.String()); err != nil {
		log.Println(err)
		return response.New(500, "Fail to upload image", nil)
	}

	id, err := s.r.Create(&branch)
	if err != nil {
		return response.New(500, "Fail to create branch", nil)
	}

	return response.New(201, "Branch created", model.CreateBranchResponse{ID: id.String()})
}

func (s *branchService) FindAll() response.Response {
	branches, err := s.r.FindAll()
	if err != nil {
		return response.New(500, "Fail to fetch branches", nil)
	}

	res := make([]model.FindBranchResponse, len(branches))
	for i, branch := range branches {
		services := make([]model.FindServiceResponse, len(branch.Services))
		for j, service := range branch.Services {
			services[j] = model.FindServiceResponse{
				ID:             service.ID.String(),
				Name:           service.Name,
				Description:    service.Description,
				Price:          service.Price,
				DurationMinute: service.DurationMinute,
				ImageURL:       s.uploader.GetURL("img/services/" + service.ID.String()),
			}
		}
		res[i] = model.FindBranchResponse{
			ID:           branch.ID.String(),
			Name:         branch.Name,
			Address:      branch.Address,
			MapsURL:      branch.MapsURL,
			Phone:        branch.Phone,
			OpenTime:     branch.OpenTime.Format(time.Kitchen),
			CloseTime:    branch.CloseTime.Format(time.Kitchen),
			TimeZoneName: branch.TimeZoneName,
			Services:     services,
			ImageURL:     s.uploader.GetURL("img/branches/" + branch.ID.String()),
		}
	}

	return response.New(200, "Branches fetched", res)
}

func (s *branchService) SetServices(req model.SetServicesToBranchRequest) response.Response {
	branchID, err := ulid.Parse(req.BranchID)
	if err != nil {
		return response.New(400, "Fail to parse branchID", err.Error())
	}

	serviceIDs := make([]ulid.ULID, len(req.ServiceIDs))
	for i, id := range req.ServiceIDs {
		serviceID, err := ulid.Parse(id)
		if err != nil {
			return response.New(400, "Fail to parse serviceID", err.Error())
		}
		serviceIDs[i] = serviceID
	}

	err = s.r.SetServices(branchID, serviceIDs)
	if err != nil {
		return response.New(500, "Fail to set services", nil)
	}

	return response.New(200, "Services set", nil)
}
