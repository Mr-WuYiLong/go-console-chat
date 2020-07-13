package exception

// Exception 用户错误结构体
type Exception struct {
	// ErrorType byte
	errorMsg string
}

// NewException 实例化Exception
func NewException(errorMsg string) Exception {
	return Exception{
		errorMsg: errorMsg,
	}
}

// Error 错误信息
func (exception Exception) Error() string {
	return exception.errorMsg
}
