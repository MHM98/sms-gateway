package scheduler

import (
	"context"
	"log"
	"time"

	controllermodel "sms-gateway/models/controller"
)

type messageDispatchJob struct {
	dispatcher  IMessageDispatcher
	serviceType controllermodel.ServiceType
	limit       int
	timeout     time.Duration
}

func NewMessageDispatchJob(dispatcher IMessageDispatcher, serviceType controllermodel.ServiceType, limit int, timeout time.Duration) *messageDispatchJob {
	return &messageDispatchJob{
		dispatcher:  dispatcher,
		serviceType: serviceType,
		limit:       limit,
		timeout:     timeout,
	}
}

func (j *messageDispatchJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), j.timeout)
	defer cancel()

	if err := j.dispatcher.DispatchPendingMessages(ctx, j.serviceType, j.limit); err != nil {
		log.Printf(
			"dispatch pending messages failed: service_type=%s limit=%d error=%v",
			j.serviceType,
			j.limit,
			err,
		)
	}
}
