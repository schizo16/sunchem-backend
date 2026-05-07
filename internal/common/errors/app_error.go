package errors

import "fmt"

type AppError struct {
	HttpStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInternalServer = &AppError{HttpStatus: 500, Code: "INTERNAL_ERROR", Message: "Lỗi hệ thống"}
	ErrBadRequest     = &AppError{HttpStatus: 400, Code: "BAD_REQUEST", Message: "Dữ liệu không hợp lệ"}
	ErrNotFound       = &AppError{HttpStatus: 404, Code: "NOT_FOUND", Message: "Không tìm thấy"}
	ErrUnauthorized   = &AppError{HttpStatus: 401, Code: "UNAUTHORIZED", Message: "Chưa đăng nhập"}
	ErrForbidden      = &AppError{HttpStatus: 403, Code: "FORBIDDEN", Message: "Không có quyền truy cập"}
)

func NewError(httpStatus int, code, message string) *AppError {
	return &AppError{HttpStatus: httpStatus, Code: code, Message: message}
}

func Wrap(err error, httpStatus int, code, message string) *AppError {
	return &AppError{HttpStatus: httpStatus, Code: code, Message: fmt.Sprintf("%s: %v", message, err)}
}
