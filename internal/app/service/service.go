package service

import (
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/gcloud"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"log"
)

type IServiceService interface {
	Create(service *model.CreateServiceRequest) response.Response
	FindAll() response.Response
}

type serviceService struct {
	serviceRepository repository.IServiceRepository
	uploader          *gcloud.FileUploaderClient
}

func NewServiceService(serviceRepository repository.IServiceRepository, uploader *gcloud.FileUploaderClient) IServiceService {
	return &serviceService{serviceRepository: serviceRepository, uploader: uploader}
}

func (s *serviceService) Create(req *model.CreateServiceRequest) response.Response {
	imageFile, err := req.ImageFile.Open()
	if err != nil {
		return response.New(400, "Fail to open image", err.Error())
	}

	service := entity.Service{
		Model:          gorm.Model{},
		ID:             ulid.Make(),
		Name:           req.Name,
		Description:    req.Description,
		Price:          req.Price,
		DurationMinute: req.DurationMinute,
	}

	if err = s.uploader.UploadFile(imageFile, "img/services/"+service.ID.String()); err != nil {
		log.Println(err)
		return response.New(500, "Fail to upload image", nil)
	}

	id, err := s.serviceRepository.Create(&service)
	if err != nil {
		return response.New(500, "Fail to create service", nil)
	}

	return response.New(201, "Service created", model.CreateServiceResponse{ID: id.String()})
}

func (s *serviceService) FindAll() response.Response {
	services, err := s.serviceRepository.FindAll()
	if err != nil {
		return response.New(500, "Fail to fetch services", nil)
	}

	res := make([]model.FindServiceResponse, len(services))
	for i, service := range services {
		res[i] = model.FindServiceResponse{
			ID:             service.ID.String(),
			Name:           service.Name,
			Description:    service.Description,
			DurationMinute: service.DurationMinute,
			Price:          service.Price,
			ImageURL:       s.uploader.GetURL("img/services/" + service.ID.String()),
		}
	}

	return response.New(200, "Services fetched", res)
}
