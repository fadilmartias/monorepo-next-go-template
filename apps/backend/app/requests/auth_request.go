package requests

type RegisterInput struct {
	Name                 string `json:"name" validate:"required,max=255"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Phone                string `json:"phone" validate:"required,max=255"`
	ReferralCode         string `json:"referral_code" validate:"omitempty,len=8"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginInput struct {
	Credential string `json:"credential" validate:"required,max=255"`
	Password   string `json:"password" validate:"required,min=8,max=255"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type ResetPasswordInput struct {
	Token                string `json:"token" validate:"required,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=255"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type VerifyEmailInput struct {
	Token string `json:"token" validate:"required,max=255"`
}

type Verify2FARequest struct {
	Code string `json:"code" validate:"required,max=255"`
}

type Login2FARequest struct {
	TempToken string `json:"temp_token" validate:"required,max=255"`
	Code      string `json:"code" validate:"required,max=255"`
}

type SendOTPRequest struct {
	TargetType string `json:"target_type" validate:"required,oneof=wa email"`
	Target     string `json:"target" validate:"max=255"`
	Subject    string `json:"subject" validate:"required,max=255"`
}

type VerifyOTPRequest struct {
	TargetType string `json:"target_type" validate:"required,oneof=wa email"`
	Target     string `json:"target" validate:"max=255"`
	Subject    string `json:"subject" validate:"required,max=255"`
	Otp        string `json:"otp" validate:"required,min=6,max=6"`
}
