package scheduler

import (
	"context"
	controllermodel "sms-gateway/models/controller"
)

type IMessageDispatcher interface {
	DispatchPendingMessages(ctx context.Context, serviceType controllermodel.ServiceType, limit int) error
}
