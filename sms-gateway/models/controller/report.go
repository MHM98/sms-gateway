package controllermodel

import "time"

type UserMessageReportQuery struct {
	UserID   uint64
	LastSeen uint64
	To       time.Time
	From     time.Time
}
