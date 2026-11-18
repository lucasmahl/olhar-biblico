package utils

type Result struct {
	Results     interface{} `json:"results"`
	Skip        uint64      `json:"skip"`
	Limit       uint64      `json:"limit"`
	TotalRecord int64       `json:"total_record"`
	Paginable   bool        `json:"paginable"`
}
