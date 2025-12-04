package job_tasks

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

const (
	TypeSendOTP = "send_otp"
)

type SendOTPPayload struct {
	TargetType string
	Target     string
	Subject    string
	OTP        string
	Name       string
}

func NewSendOTPTask(targetType string, target string, subject string, otp string, name string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(SendOTPPayload{TargetType: targetType, Target: target, Subject: subject, OTP: otp, Name: name})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSendOTP, payload, asynq.MaxRetry(3)), nil
}
