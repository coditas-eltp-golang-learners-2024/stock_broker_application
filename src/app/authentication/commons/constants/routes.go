package constants

// authentication route Constants
const (
	SwaggerRoute              = "/swagger/*any"
	AuthenticationHealthCheck = "/authentication/health-check"
	Metrics                   = "/metrics"
	SignUp                    = "auth/signup"
	ValidateOTP               = "auth/validate-otp"
	ForgotPassword            = "auth/forgot-password"
	ChangePassword            = "auth/change-password"
	SignIn                    = "auth/signin"
	GetStockById			  = "/stocks/:id"
)
