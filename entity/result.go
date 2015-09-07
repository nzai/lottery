package entity

type Result struct {
	Success bool
	Message string
	Data    interface{}
}

func ResultD(data interface{}) *Result {
	return ResultSMD(data != nil, "", data)
}

func ResultS(success bool) *Result {
	return ResultSM(success, "")
}

func ResultSM(success bool, message string) *Result {
	return ResultSMD(success, message, nil)
}

func ResultSMD(success bool, message string, data interface{}) *Result {
	return &Result{Success: success, Message: message, Data: data}
}
