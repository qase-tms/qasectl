package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = `https://api.qase.io/`

type api struct {
	cfg Config
}

func NewApiFromConfig() (api, error) {
	cfg, err := getConfig()
	if err != nil {
		return api{}, err
	}

	return api{
		cfg: cfg,
	}, nil
}

func (a api) CreateRun(title, description string) error {
	url := fmt.Sprintf("%s/v1/run/%s", baseURL, a.cfg.ProjectCode)

	var request = struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}{
		Title:       title,
		Description: description,
	}

	buff, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(buff))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Token", a.cfg.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buff, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response struct {
		Status       bool   `json:"status"`
		ErrorMessage string `json:"errorMessage"`
	}

	err = json.Unmarshal(buff, &response)
	if err != nil {
		return err
	}

	if !response.Status {
		return fmt.Errorf("response status is false: %s", response.ErrorMessage)
	}

	return nil
}
