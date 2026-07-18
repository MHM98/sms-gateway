package controllermodel

import "time"

type ServiceType string

const (
	ServiceTypeNormal  ServiceType = "normal"
	ServiceTypeExpress ServiceType = "express"
)

type Message struct {
	ID                       uint64
	UserID                   uint64
	ChargeAmount             uint64
	Recipient                string
	Body                     string
	ServiceType              ServiceType
	Status                   string
	CreatedAt                time.Time
	SubmissionLatencySeconds *uint64
}

type Messages []Message
