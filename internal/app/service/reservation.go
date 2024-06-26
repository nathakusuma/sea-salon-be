package service

import (
	"github.com/nathakusuma/sea-salon-be/internal/app/repository"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/model"
	"github.com/nathakusuma/sea-salon-be/internal/pkg/response"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type IReservationService interface {
	Create(req model.CreateReservationRequest) response.Response
	FindAvailableSchedules(req model.FindAvailableReservationSchedulesRequest) response.Response
}

type reservationService struct {
	r repository.IReservationRepository
}

func NewReservationService(r repository.IReservationRepository) IReservationService {
	return &reservationService{r: r}
}

func (s *reservationService) Create(req model.CreateReservationRequest) response.Response {
	startTime, err := time.Parse(time.Kitchen, req.StartTime)
	if err != nil {
		return response.New(400, "Fail to parse startTime", err.Error())
	}
	date, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return response.New(400, "Fail to parse date", err.Error())
	}
	startTime = time.Date(date.Year(), date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)

	// Check if the time is available
	isAvailable, err := s.isTimeAvailable(startTime, req.ServiceName, date)
	if err != nil {
		return response.New(500, "Fail to check time availability", err.Error())
	}
	if !isAvailable {
		return response.New(404, "That time is not available", nil)
	}

	reservation := entity.Reservation{
		Model:        gorm.Model{},
		ID:           ulid.Make(),
		CustomerName: req.CustomerName,
		PhoneNumber:  req.PhoneNumber,
		ServiceName:  req.ServiceName,
		StartTime:    startTime,
	}

	id, err := s.r.Create(&reservation)
	if err != nil {
		return response.New(500, "Fail to create reservation", nil)
	}

	return response.New(201, "Reservation created", model.CreateReservationResponse{ID: id.String()})
}

func (s *reservationService) findAvailableStartTimes(serviceName string, serviceDuration time.Duration, date time.Time) ([]time.Time, error) {
	openTime := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, time.UTC)
	closeTime := time.Date(date.Year(), date.Month(), date.Day(), 21, 0, 0, 0, time.UTC)

	reservations, err := s.r.FindByTimeRange(serviceName, openTime, closeTime)
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
	branchTimeZone := time.FixedZone("WIB", 7*60*60)
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

func (s *reservationService) FindAvailableSchedules(req model.FindAvailableReservationSchedulesRequest) response.Response {
	date, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return response.New(400, "Fail to parse date", err.Error())
	}

	duration := time.Hour

	startTimes, err := s.findAvailableStartTimes(req.ServiceName, duration, date)
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

func (s *reservationService) isTimeAvailable(timeToCheck time.Time, serviceName string, date time.Time) (bool, error) {
	availableTimes, err := s.findAvailableStartTimes(serviceName, time.Hour, date)
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
