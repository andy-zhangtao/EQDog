package main

import (
	"io/ioutil"
	"net/http"
)

func SandHttp(method, path string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
