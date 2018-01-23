package test

import (
	"net/http"
	"testing"
	"fmt"
	"os"

	"test/workload_generator/commands"

	"github.com/gavv/httpexpect"
)


func TestTransacService(t *testing.T) {
	port := os.Getenv("TRANS_PORT")
    host := os.Getenv("TRANS_HOST")
    url := fmt.Sprintf("http://%s:%s", host, port)

	e := httpexpect.New(t, url)

	user := "kevin"

	cmd := commands.Command{ Name: "ADD", Username: user, Amount: 87863.73 }
	endpoint := commands.FormatCommandEndpoint(cmd)

	e.GET(endpoint).
		Expect().
		Status(http.StatusOK)

	cmd = commands.Command{ Name: "BUY", Username: user, Symbol: "ABC", Amount: 5000 }
	endpoint = commands.FormatCommandEndpoint(cmd)

	e.GET(endpoint).
		Expect().
		Status(http.StatusOK)

	cmd = commands.Command{ Name: "COMMIT_BUY", Username: user }
	endpoint = commands.FormatCommandEndpoint(cmd)

	e.GET(endpoint).
		Expect().
		Status(http.StatusOK)

}