package test

import (
	"net/http"
	"testing"
	"fmt"
	"os"

	"github.com/gavv/httpexpect"
)


func TestAddUser(t *testing.T) {
	port := os.Getenv("TRANS_PORT")
    host := os.Getenv("TRANS_HOST")
    url := fmt.Sprintf("http://%s:%s", host, port)

	e := httpexpect.New(t, url)
	e.GET("/api/addUser/oY01WVirLr/87863.730000").
		Expect().
		Status(http.StatusOK)
}