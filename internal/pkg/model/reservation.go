package model

type CreateReservationRequest struct {
	CustomerName string `json:"customerName" validate:"required,max=50,min=1"`
	PhoneNumber  string `json:"phoneNumber" validate:"required,max=15,min=1"`
	ServiceName  string `json:"serviceName" validate:"required,max=50,min=1"`
	StartTime    string `json:"startTime" validate:"required"`
	Date         string `json:"date" validate:"required"`
}

type CreateReservationResponse struct {
	ID string `json:"id"`
}

type FindAvailableReservationSchedulesRequest struct {
	ServiceName string `form:"serviceName" validate:"required,max=50,min=1"`
	Date        string `form:"date" validate:"required"`
}

type FindAvailableReservationScheduleResponse struct {
	StartTime  string `json:"startTime"`
	FinishTime string `json:"finishTime"`
}
