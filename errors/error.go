package errors

import "github.com/gin-gonic/gin"

type Error interface {
	StatusCode() StatusCode
	Error() string
	StackTrace() string
	Data() gin.H
}

// -------------------------------------------------------------------------- //
type StatusCode int

func (this StatusCode) Int() int {
	return int(this)
}

func (this StatusCode) Is2xxSuccess() bool {
	return 200 <= this && this < 300
}

func (this StatusCode) Is4xxClientError() bool {
	return 400 <= this && this < 500
}

func (this StatusCode) Is5xxServerError() bool {
	return 500 <= this && this < 600
}

func (this StatusCode) Is501NotImplemented() bool {
	return this == 501
}

// -------------------------------------------------------------------------- //
type DefaultError struct {
	statusCode StatusCode
	error      string
	stackTrace string
	data       gin.H
}

func NewError(statusCode int, errMsg string, stackTrace string, data gin.H) DefaultError {
	return DefaultError{
		statusCode: StatusCode(statusCode),
		error:      errMsg,
		stackTrace: stackTrace,
		data:       data,
	}
}

func (this DefaultError) StatusCode() StatusCode {
	return this.statusCode
}

func (this DefaultError) Error() string {
	return this.error
}

func (this DefaultError) StackTrace() string {
	return this.stackTrace
}

func (this DefaultError) Data() gin.H {
	return this.data
}
