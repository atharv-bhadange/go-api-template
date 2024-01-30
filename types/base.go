package types

type ServiceError struct {
	Message string
	Error   error
	Code    int
}
