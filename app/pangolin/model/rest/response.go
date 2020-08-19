package rest

import (
	"pangolin/app/pangolin/model/db"
)

type BasicResponse struct {
	Meta *ResponseMeta `json:"meta"`
	Data interface{}   `json:"data"`
}

type ResponseMeta struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	InnerError error  `json:"inner_error"`
}

type ListTunnelConfigResponseData struct {
	Rows       []*db.TunnelConfig `json:"rows"`
	TotalCount int                `json:"total_count"`
}
