package controllermodel

type ServiceType string

const (
	ServiceTypeNormal  ServiceType = "normal"
	ServiceTypeExpress ServiceType = "express"
)

type Message struct {
	ID          uint64
	UserID      uint64
	Recipient   string
	Body        string
	ServiceType ServiceType
}
