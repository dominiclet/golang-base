package httpresp

type StandardDataResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StandardResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
