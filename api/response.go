package api

type Response struct {
	Status int64       `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func NewSuccessResponse(data interface{}) Response {
	return Response{
		Status: 0,
		Data:   data,
	}
}

func NewFailResponse(err error, data interface{}) Response {
	return Response{
		Status: 1,
		Error:  err.Error(),
		Data:   data,
	}
}
