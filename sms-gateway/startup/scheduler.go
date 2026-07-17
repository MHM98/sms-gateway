package startup

import (
	"fmt"
	"time"

	controllermodel "sms-gateway/models/controller"
	"sms-gateway/scheduler"

	"github.com/robfig/cron/v3"
)

const (
	normalDispatchSpec   = "@every 1m"
	normalDispatchLimit  = 500
	expressDispatchSpec  = "@every 10s"
	expressDispatchLimit = 1000
	dispatchJobTimeout   = 30 * time.Second
)

type dispatchJobConfig struct {
	Spec        string
	ServiceType controllermodel.ServiceType
	Limit       int
	Timeout     time.Duration
}

func buildScheduler(messageController scheduler.IMessageDispatcher) (*cron.Cron, error) {
	cronRunner := cron.New(
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	jobConfigs := []dispatchJobConfig{
		{
			Spec:        normalDispatchSpec,
			ServiceType: controllermodel.ServiceTypeNormal,
			Limit:       normalDispatchLimit,
			Timeout:     dispatchJobTimeout,
		},
		{
			Spec:        expressDispatchSpec,
			ServiceType: controllermodel.ServiceTypeExpress,
			Limit:       expressDispatchLimit,
			Timeout:     dispatchJobTimeout,
		},
	}

	for _, config := range jobConfigs {
		job := scheduler.NewMessageDispatchJob(
			messageController,
			config.ServiceType,
			config.Limit,
			config.Timeout,
		)
		if _, err := cronRunner.AddJob(config.Spec, job); err != nil {
			<-cronRunner.Stop().Done()
			return nil, fmt.Errorf("register %s message dispatch job: %w", config.ServiceType, err)
		}
	}

	return cronRunner, nil
}
