package usecases

import (
	"fmt"

	"github.com/fadilmartias/dilz_code/apps/backend/app/jobs"
	job_tasks "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/tasks"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
)

type SendResetPasswordEmailUsecase struct {
}

func NewSendResetPasswordEmailUsecase() *SendResetPasswordEmailUsecase {
	return &SendResetPasswordEmailUsecase{}
}

func (uc *SendResetPasswordEmailUsecase) Execute(email string, token string, name string) error {
	// Send email
	task, err := job_tasks.NewEmailResetPasswordTask(email, token, name)
	if err != nil {
		logger.Error("failed to create task:", err)
	}
	info, err := jobs.AsynqClient.Enqueue(task)
	if err != nil {
		logger.Error("failed to enqueue task:", err)
	}
	fmt.Println("enqueued task:", info)
	return nil
}
