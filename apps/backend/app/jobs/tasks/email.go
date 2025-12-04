package job_tasks

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailVerification  = "email:verification"
	TypeEmailResetPassword = "email:reset-password"
)

type EmailVerificationPayload struct {
	To    string
	Token string
	Name  string
}

type EmailResetPasswordPayload struct {
	To    string
	Token string
	Name  string
}

func NewEmailVerificationTask(to, token, name string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailVerificationPayload{To: to, Token: token, Name: name})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmailVerification, payload, asynq.MaxRetry(3)), nil
}

func NewEmailResetPasswordTask(to, token, name string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailResetPasswordPayload{To: to, Token: token, Name: name})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeEmailResetPassword, payload, asynq.MaxRetry(3)), nil
}
