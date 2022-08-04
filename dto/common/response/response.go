package response

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func New(code int, data interface{}, msg string) Result {
	// 开始时间
	return Result{
		Code: code,
		Data: data,
		Msg:  msg,
	}
}

func Ok() Result {
	return New(SUCCESS, map[string]interface{}{}, "success")
}

func OkWithMessage(message string) Result {
	return New(SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(data interface{}) Result {
	return New(SUCCESS, data, "success")
}

func OkWithDetailed(data interface{}, message string) {
	New(SUCCESS, data, message)
}

func Fail() Result {
	return New(ERROR, map[string]interface{}{}, "success")
}

func FailWithMessage(message string) Result {
	return New(ERROR, map[string]interface{}{}, message)
}

func FailWithDetailed(data interface{}, message string) Result {
	return New(ERROR, data, message)
}
