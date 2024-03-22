package dto

type ErrorResponse struct {
	Code    int                    `json:"code"`
	Message int                    `json:"message"`
	Details map[string]interface{} `json:"details"`
}
