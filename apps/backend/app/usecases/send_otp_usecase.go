package usecases

import (
	"github.com/fadilmartias/dilz_code/apps/backend/app/jobs"
	job_tasks "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/tasks"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
)

type SendOTPUsecase struct {
}

func NewSendOTPUsecase() *SendOTPUsecase {
	return &SendOTPUsecase{}
}

func (uc *SendOTPUsecase) Execute(targetType string, target, subject, otp string, name string) error {
	task, err := job_tasks.NewSendOTPTask(targetType, target, subject, otp, name)
	if err != nil {
		logger.Error("error create task:", err)
		return err
	}
	_, err = jobs.AsynqClient.Enqueue(task)
	if err != nil {
		logger.Error("error enqueue task:", err)
		return err
	}
	return nil
}
