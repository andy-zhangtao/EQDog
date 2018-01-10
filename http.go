package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
	"strconv"
	"strings"
)

func SandHttp(method, path string) ([]byte, error) {
	client := &http.Client{}
	if debug {
		fmt.Printf("[SandHttp] [%s] [%s]\n", method, path)
	}

	req, err := http.NewRequest(method, path, nil)
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

	if resp.StatusCode != 200 {
		return nil, errors.New("Response Status Code " + strconv.Itoa(resp.StatusCode))
	}
	return content, nil
}

func SandHttpBody(method, path, body string) ([]byte, error) {
	client := &http.Client{}
	if debug {
		fmt.Printf("[SandHttp] [%s] [%s] [%s]\n", method, path, body)
	}

	req, err := http.NewRequest(method, path, strings.NewReader(body))
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

	if resp.StatusCode != 200 {
		return nil, errors.New("Response Status Code " + strconv.Itoa(resp.StatusCode))
	}
	return content, nil
}
