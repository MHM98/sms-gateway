package apimodel

type MessageRequest struct {
	UserID      uint64 `json:"user_id" validate:"gt=0"`
	Recipient   string `json:"recipient" validate:"required,max=20"`
	Body        string `json:"body" validate:"required,max=255"`
	ServiceType string `json:"service_type" validate:"required,oneof=normal express"`
}
