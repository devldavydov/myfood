package model

type TemplateData struct {
	Nav  string
	Data any
}

type ResponseData struct {
	Error string `json:"error"`
	Data  any    `json:"data"`
}

func NewErrorResponse(err string) *ResponseData {
	return &ResponseData{Error: err}
}

func NewDataResponse(data any) *ResponseData {
	return &ResponseData{Data: data}
}

func NewOKResponse() *ResponseData {
	return &ResponseData{Data: "ok"}
}
