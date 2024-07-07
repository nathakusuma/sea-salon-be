package service

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/jwt"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type IReservationService interface {
	Create(req model.CreateReservationRequest, userClaims jwt.Claims) response.Response
	FindAvailableSchedules(req model.FindAvailableReservationSchedulesRequest) response.Response
	FindByUser(userClaims jwt.Claims) response.Response
	FindByDateAndBranch(req model.AdminFindReservationRequest) response.Response
}

type reservationService struct {
	r  repository.IReservationRepository
	sr repository.IServiceRepository
	br repository.IBranchRepository
}

func NewReservationService(r repository.IReservationRepository, sr repository.IServiceRepository, br repository.IBranchRepository) IReservationService {
	return &reservationService{r: r, sr: sr, br: br}
}

func (s *reservationService) Create(req model.CreateReservationRequest, userClaims jwt.Claims) response.Response {
	userID, err := ulid.Parse(userClaims.Subject)
	if err != nil {
		return response.New(400, "Fail to parse userID", err.Error())
	}

	startTime, err := time.Parse(time.Kitchen, req.StartTime)
	if err != nil {
		return response.New(400, "Fail to parse startTime", err.Error())
	}
	date, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return response.New(400, "Fail to parse date", err.Error())
	}
	startTime = time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)

	branch, service, errRes := s.validateBranchService(req.BranchID, req.ServiceID)
	if errRes.HttpCode != 0 {
		return errRes
	}

	duration := time.Duration(service.DurationMinute) * time.Minute

	// Check if the time is available
	isAvailable, err := s.isTimeAvailable(startTime, branch, service.ID, duration, date)
	if err != nil {
		return response.New(500, "Fail to check time availability", err.Error())
	}
	if !isAvailable {
		return response.New(404, "That time is not available", nil)
	}

	reservation := entity.Reservation{
		Model:     gorm.Model{},
		ID:        ulid.Make(),
		UserID:    userID,
		BranchID:  branch.ID,
		ServiceID: service.ID,
		StartTime: startTime,
	}

	id, err := s.r.Create(&reservation)
	if err != nil {
		return response.New(500, "Fail to create reservation", nil)
	}

	return response.New(201, "Reservation created", model.CreateReservationResponse{ID: id.String()})
}

func (s *reservationService) findAvailableStartTimes(branch entity.Branch, serviceID ulid.ULID, serviceDuration time.Duration, date time.Time) ([]time.Time, error) {
	openTime := time.Date(date.Year(), date.Month(), date.Day(), branch.OpenTime.Hour(), branch.OpenTime.Minute(), 0, 0, time.UTC)
	closeTime := time.Date(date.Year(), date.Month(), date.Day(), branch.CloseTime.Hour(), branch.CloseTime.Minute(), 0, 0, time.UTC)

	reservations, err := s.r.FindByTimeRange(branch.ID, serviceID, openTime, closeTime)
	if err != nil {
		return nil, err
	}

	var availableStartTimes []time.Time

	// Create a map to track booked times
	bookedTimes := make(map[time.Time]bool)
	for _, reservation := range reservations {
		bookedTimes[reservation.StartTime] = true
	}

	// disable booking for time that already passed
	branchTimeZone := time.FixedZone(branch.TimeZoneName, branch.TimeZoneOffset)
	currTimeInBranchZone := time.Now().In(branchTimeZone)
	currTimeFakeTz := time.Date(currTimeInBranchZone.Year(), currTimeInBranchZone.Month(), currTimeInBranchZone.Day(), currTimeInBranchZone.Hour(), currTimeInBranchZone.Minute(), currTimeInBranchZone.Second(), currTimeInBranchZone.Nanosecond(), time.UTC)
	for current := openTime; current.Before(currTimeFakeTz); current = current.Add(serviceDuration) {
		bookedTimes[current] = true
	}

	// Iterate over each possible schedule time from start to end
	for current := openTime; current.Before(closeTime); current = current.Add(serviceDuration) {
		if !bookedTimes[current] {
			availableStartTimes = append(availableStartTimes, current)
		}
	}

	return availableStartTimes, nil
}

func (s *reservationService) validateBranchService(branchIDStr string, serviceIDStr string) (entity.Branch, entity.Service, response.Response) {
	branchID, err := ulid.Parse(branchIDStr)
	if err != nil {
		return entity.Branch{}, entity.Service{}, response.New(400, "Fail to parse branchID", nil)
	}

	branch, err := s.br.FindByID(branchID)
	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if ok && pgErr.Code == "42703" {
			return entity.Branch{}, entity.Service{}, response.New(404, "That branch does not exist", nil)
		}
		return entity.Branch{}, entity.Service{}, response.New(500, "Fail to find branch", nil)
	}

	serviceID, err := ulid.Parse(serviceIDStr)
	if err != nil {
		return entity.Branch{}, entity.Service{}, response.New(400, "Fail to parse serviceID", nil)
	}

	var service *entity.Service
	for _, s := range branch.Services {
		if s.ID == serviceID {
			service = s
			break
		}
	}
	if service == nil {
		return entity.Branch{}, entity.Service{}, response.New(404, "That service does not exist in the branch", nil)
	}

	return branch, *service, response.Response{}
}

func (s *reservationService) FindAvailableSchedules(req model.FindAvailableReservationSchedulesRequest) response.Response {
	date, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return response.New(400, "Fail to parse date", err.Error())
	}

	branch, service, errRes := s.validateBranchService(req.BranchID, req.ServiceID)
	if errRes.HttpCode != 0 {
		return errRes
	}

	duration := time.Duration(service.DurationMinute) * time.Minute

	startTimes, err := s.findAvailableStartTimes(branch, service.ID, duration, date)
	if err != nil {
		return response.New(500, "Fail to find available schedules", err.Error())
	}

	availableSchedules := make([]model.FindAvailableReservationScheduleResponse, 0)
	for _, startTime := range startTimes {
		availableSchedules = append(availableSchedules, model.FindAvailableReservationScheduleResponse{
			StartTime:  startTime.Format(time.Kitchen),
			FinishTime: startTime.Add(duration).Format(time.Kitchen),
		})
	}

	return response.New(200, "Available schedules fetched", availableSchedules)
}

func (s *reservationService) isTimeAvailable(timeToCheck time.Time, branch entity.Branch, serviceID ulid.ULID, serviceDuration time.Duration, date time.Time) (bool, error) {
	availableTimes, err := s.findAvailableStartTimes(branch, serviceID, serviceDuration, date)
	if err != nil {
		return false, err
	}

	// check if the timeToCheck is in the availableTimes
	for _, availableTime := range availableTimes {
		if availableTime.Equal(timeToCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (s *reservationService) FindByUser(userClaims jwt.Claims) response.Response {
	userID, err := ulid.Parse(userClaims.Subject)
	if err != nil {
		return response.New(400, "Fail to parse userID", err.Error())
	}

	reservations, err := s.r.FindByUserID(userID)
	if err != nil {
		return response.New(500, "Fail to find reservations", err.Error())
	}

	res := make([]model.FindReservationResponse, len(reservations))
	for i, reservation := range reservations {
		duration := time.Duration(reservation.Service.DurationMinute) * time.Minute
		res[i] = model.FindReservationResponse{
			ID:          reservation.ID.String(),
			Date:        reservation.StartTime.Format(time.DateOnly),
			BranchName:  reservation.Branch.Name,
			ServiceName: reservation.Service.Name,
			Time:        fmt.Sprintf("%s - %s %s", reservation.StartTime.Format(time.Kitchen), reservation.StartTime.Add(duration).Format(time.Kitchen), reservation.Branch.TimeZoneName),
		}
	}

	return response.New(200, "Reservations fetched", res)
}

func (s *reservationService) FindByDateAndBranch(req model.AdminFindReservationRequest) response.Response {
	date, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return response.New(400, "Fail to parse date", err.Error())
	}

	branchID, err := ulid.Parse(req.BranchID)
	if err != nil {
		return response.New(400, "Fail to parse branchID", err.Error())
	}

	reservations, err := s.r.FindByDateAndBranch(date, branchID)
	if err != nil {
		return response.New(500, "Fail to find reservations", nil)
	}

	res := make([]model.AdminFindReservationResponse, len(reservations))
	for i, reservation := range reservations {
		duration := time.Duration(reservation.Service.DurationMinute) * time.Minute
		res[i] = model.AdminFindReservationResponse{
			ID:           reservation.ID.String(),
			CustomerName: reservation.User.FullName,
			Email:        reservation.User.Email,
			PhoneNumber:  reservation.User.PhoneNumber,
			BranchName:   reservation.Branch.Name,
			ServiceName:  reservation.Service.Name,
			Time:         fmt.Sprintf("%s - %s %s", reservation.StartTime.Format(time.Kitchen), reservation.StartTime.Add(duration).Format(time.Kitchen), reservation.Branch.TimeZoneName),
		}
	}

	return response.New(200, "Reservations fetched", res)
}
