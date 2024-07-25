package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var (
	t = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 5,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     300 * time.Second,
	}
	readBytesSize = 4096
)

type HttpClient struct {
	Client     *http.Client
	MaxRetries int
}

func NewHttpClient(n int) *HttpClient {
	return &HttpClient{
		MaxRetries: n,
		Client: &http.Client{
			Transport: t,
			Timeout:   300 * time.Second,
		},
	}
}

func (hc *HttpClient) DoWait(req *http.Request, jsonResp interface{}) (statusCode int, err error) {
	if jsonResp == nil {
		return http.StatusBadRequest, errors.New("JSON receiver not configured")
	}

	resp, err := hc.DoSend(req)
	if err != nil || resp == nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	lgr := logrus.WithFields(logrus.Fields{
		"req-url":    req.URL,
		"robot-uuid": req.Header.Get("robot-uuid"),
		"res-status": resp.Status,
	})
	statusCode = resp.StatusCode
	err = json.NewDecoder(resp.Body).Decode(jsonResp)
	if err == nil {
		lgr.Info(jsonResp)
	} else {
		lgr.Error(err)
	}

	return
}

func (hc *HttpClient) Download(req *http.Request) (r []byte, statusCode int, err error) {
	resp, err := hc.DoSend(req)
	if err != nil || resp == nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if code := resp.StatusCode; code < 200 || code > 299 {
		statusCode = code

		var rb []byte
		if rb, err = io.ReadAll(resp.Body); err == nil {
			err = fmt.Errorf("response has status:%s and body:%q", resp.Status, rb)
		}

	} else {
		r, err = io.ReadAll(resp.Body)
	}

	return
}

func (hc *HttpClient) DoSend(req *http.Request) (resp *http.Response, err error) {
	if resp, err = hc.Client.Do(req); err == nil {
		return
	}

	maxRetries := hc.MaxRetries
	backoff := 1000 * time.Millisecond

	for retries := 1; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = hc.Client.Do(req); err == nil {
			break
		}
	}
	if err != nil {
		logrus.Error(err)
	}
	return
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(t); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
