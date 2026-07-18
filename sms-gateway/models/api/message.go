package apimodel

import "time"

type MessageRequest struct {
	UserID      uint64 `json:"user_id" validate:"gt=0"`
	Recipient   string `json:"recipient" validate:"required,max=20"`
	Body        string `json:"body" validate:"required,max=255"`
	ServiceType string `json:"service_type" validate:"required,oneof=normal express"`
}

type UserMessageReportRequest struct {
	UserID uint64 `json:"user_id" validate:"gt=0"`
	From   string `json:"from" validate:"required,datetime=2006-01-02"`
	To     string `json:"to" validate:"required,datetime=2006-01-02"`
}

type UserMessageReportResponse struct {
	UserID   uint64                 `json:"user_id"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	Messages MessagesReportResponse `json:"messages"`
}

type MessageReportResponse struct {
	ID                       uint64    `json:"id"`
	Recipient                string    `json:"recipient"`
	Body                     string    `json:"body"`
	ServiceType              string    `json:"service_type"`
	Status                   string    `json:"status"`
	CreatedAt                time.Time `json:"created_at"`
	SubmissionLatencySeconds *uint64   `json:"submission_latency_seconds"`
}

type MessagesReportResponse []MessageReportResponse
