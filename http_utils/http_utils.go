package httputils

import (
	"net/http"
	"time"
)

func InitClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
