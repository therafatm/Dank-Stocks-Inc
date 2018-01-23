package test

import (
	"net/http"
	"testing"
	"fmt"
	"os"

	"github.com/gavv/httpexpect"
)


func TestQuote(t *testing.T) {
	port := os.Getenv("QUOTE_SERVER_PORT")
    host := os.Getenv("QUOTE_SERVER_HOST")
    url := fmt.Sprintf("http://%s:%s", host, port)

	e := httpexpect.New(t, url)
	e.GET("/api/getQuote/test_user/AMD").
		Expect().
		Status(http.StatusOK)
}