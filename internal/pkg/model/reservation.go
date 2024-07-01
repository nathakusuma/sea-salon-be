package model

type CreateReservationRequest struct {
	ServiceID string `json:"serviceId" validate:"required,max=50,min=1"`
	StartTime string `json:"startTime" validate:"required"`
	Date      string `json:"date" validate:"required"`
}

type CreateReservationResponse struct {
	ID string `json:"id"`
}

type FindAvailableReservationSchedulesRequest struct {
	ServiceID string `form:"serviceId" validate:"required,max=50,min=1"`
	Date      string `form:"date" validate:"required"`
}

type FindAvailableReservationScheduleResponse struct {
	StartTime  string `json:"startTime"`
	FinishTime string `json:"finishTime"`
}

type FindReservationResponse struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	ServiceName string `json:"serviceName"`
	StartTime   string `json:"startTime"`
	FinishTime  string `json:"finishTime"`
}

type AdminFindReservationRequest struct {
	Date string `form:"date" validate:"required"`
}

type AdminFindReservationResponse struct {
	ID           string `json:"id"`
	CustomerName string `json:"customerName"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	ServiceName  string `json:"serviceName"`
	StartTime    string `json:"startTime"`
	FinishTime   string `json:"finishTime"`
}
