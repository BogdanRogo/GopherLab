package models

type OutResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SetKeyParams struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}
