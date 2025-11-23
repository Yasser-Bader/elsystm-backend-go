package util

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func WorkerMakeRequest(body interface{}) {
	godotenv.Load()
	requestByte, _ := json.Marshal(body)
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest("POST", os.Getenv("SERVICE_WORKER"), requestReader)
	if err != nil {
		return
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
}
