package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"simplicity/items"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var testTimestamp, _ = time.Parse(time.RFC3339, "2024-12-22T18:37:56.871781+01:00")

func runTestServer() *httptest.Server {
	log.SetOutput(io.Discard)
	registry := items.NewInMemoryRegistry(func() time.Time {
		return testTimestamp
	})
	return httptest.NewServer(setupServer(registry))
}

type Request struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Body   string `yaml:"body,omitempty"`
}

type Response struct {
	Code int    `yaml:"code"`
	Body string `yaml:"body,omitempty"`
}

type Call struct {
	Request  Request  `yaml:"request"`
	Response Response `yaml:"response"`
}

type TestCase struct {
	Name  string `yaml:"name"`
	Calls []Call `yaml:"calls"`
}

type TestSuite struct {
	TestCases []TestCase `yaml:"test_cases"`
}

func Test_endpoints(t *testing.T) {
	data, err := os.ReadFile("testcases.yaml")
	assert.NoError(t, err)

	var suite TestSuite
	err = yaml.Unmarshal(data, &suite)
	assert.NoError(t, err)

	for _, c := range suite.TestCases {
		t.Run(c.Name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()

			for _, call := range c.Calls {
				req, _ := http.NewRequest(call.Request.Method, fmt.Sprintf("%s%s", ts.URL, call.Request.Path), strings.NewReader(call.Request.Body))
				resp, err := http.DefaultClient.Do(req)
				body := readBody(resp)

				assert.Nil(t, err)
				if call.Response.Code != 0 {
					assert.Equal(t, call.Response.Code, resp.StatusCode, "body:%s", body)
				}
				if call.Response.Body != "" {
					assert.Nil(t, err)
					assert.JSONEq(t, call.Response.Body, body)
				}
			}
		})
	}
}

func readBody(resp *http.Response) string {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error reading body: " + err.Error()
	}
	return string(bodyBytes)
}
