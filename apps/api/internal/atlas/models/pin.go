package models

type PinOperationResult struct {
	Success bool     `json:"success"`
	Error   *PinError `json:"error"`
}

type PinError struct {
	Message string       `json:"message"`
	Code    PinErrorCode `json:"code"`
}

type PinErrorCode string

const (
	PinErrorWrongPin          PinErrorCode = "WRONG_PIN"
	PinErrorAlreadyEnabled    PinErrorCode = "PIN_ALREADY_ENABLED"
	PinErrorAlreadyDisabled   PinErrorCode = "PIN_ALREADY_DISABLED"
	PinErrorTooShort          PinErrorCode = "PIN_TOO_SHORT"
	PinErrorTooLong           PinErrorCode = "PIN_TOO_LONG"
	PinErrorSessionExpired    PinErrorCode = "SESSION_EXPIRED"
	PinErrorInternal          PinErrorCode = "INTERNAL_ERROR"
)

type PinEnableInput struct {
	Pin string `json:"pin"`
}

type PinDisableInput struct {
	CurrentPin string `json:"currentPin"`
}

type PinChangeInput struct {
	CurrentPin string `json:"currentPin"`
	NewPin     string `json:"newPin"`
}