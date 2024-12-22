package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"simplicity/storage"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testTimestamp, _ = time.Parse(time.RFC3339, "2024-12-22T18:37:56.871781+01:00")

func runTestServer() *httptest.Server {
	log.SetOutput(io.Discard)
	return httptest.NewServer(setupServer(func() time.Time {
		return testTimestamp
	}))
}

type Request struct {
	Path   string
	Method string
	Body   string
}

type Response struct {
	Code int
	Body string
}

type Call struct {
	Request  Request
	Response Response
}

type TestCase struct {
	Name  string
	Calls []Call
}

func Test_endpoints(t *testing.T) {
	cases := []TestCase{
		{
			Name: "it should return 200 when health is ok",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/health",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
					},
				},
			},
		},
		{
			Name: "it should return 404 when endpoint is not found",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/something",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 404,
					},
				},
			},
		},
		{
			Name: "it should return 400 when item is invalid",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"name": "item1"}`,
					},
					Response: Response{
						Code: 400,
					},
				},
			},
		},
		{
			Name: "it should return 405 when item post path is invalid",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 405,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 405,
					},
				},
			},
		},
		{
			Name: "it should successfully create an item",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 201,
					},
				},
			},
		},
		{
			Name: "it should fail when create an item dubliate",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 201,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item2"}`,
					},
					Response: Response{
						Code: 400,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
						Body: makeItem("1", "item1", ""),
					},
				},
			},
		},
		{
			Name: "it should successfully read an item",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item1", "description": "description1", "tags": ["tag1", "tag2"]}`,
					},
					Response: Response{
						Code: 201,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
						Body: makeItem("1", "item1", "description1", "tag1", "tag2"),
					},
				},
			},
		},
		{
			Name: "it should return 404 read item when item does not exist",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 404,
					},
				},
			},
		},
		{
			Name: "it should successfully list items",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 201,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/2",
						Method: http.MethodPost,
						Body:   `{"title": "item2"}`,
					},
					Response: Response{
						Code: 201,
					},
				},
				{
					Request: Request{
						Path:   "/api/item",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
						Body: makeItems(makeItem("1", "item1", ""), makeItem("2", "item2", "")),
					},
				},
				{
					Request: Request{
						Path:   "/api/item/",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
						Body: makeItems(makeItem("1", "item1", ""), makeItem("2", "item2", "")),
					},
				},
			},
		},
		{
			Name: "it should successfully delete existing item",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodPost,
						Body:   `{"title": "item1"}`,
					},
					Response: Response{
						Code: 201,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodDelete,
					},
					Response: Response{
						Code: 200,
					},
				},
				{
					Request: Request{
						Path:   "/api/item/",
						Method: http.MethodGet,
					},
					Response: Response{
						Code: 200,
						Body: "[]",
					},
				},
			},
		},
		{
			Name: "it should return 404 delete non-existing item",
			Calls: []Call{
				{
					Request: Request{
						Path:   "/api/item/1",
						Method: http.MethodDelete,
					},
					Response: Response{
						Code: 404,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := runTestServer()
			defer ts.Close()

			for _, call := range c.Calls {
				req, _ := http.NewRequest(call.Request.Method, fmt.Sprintf("%s%s", ts.URL, call.Request.Path), strings.NewReader(call.Request.Body))
				resp, err := http.DefaultClient.Do(req)

				assert.Nil(t, err)
				if call.Response.Code != 0 {
					assert.Equal(t, call.Response.Code, resp.StatusCode)
				}
				if call.Response.Body != "" {
					bodyBytes, err := io.ReadAll(resp.Body)
					assert.Nil(t, err)
					assert.JSONEq(t, call.Response.Body, string(bodyBytes))
				}
			}
		})
	}

}

func makeItem(id, title, description string, tags ...string) string {
	item := storage.Item{
		ItemMetadata: storage.ItemMetadata{
			ID:        id,
			CreatedAt: testTimestamp,
			UpdatedAt: testTimestamp,
		},
		ItemData: storage.ItemData{
			Title:       title,
			Description: description,
			Tags:        tags,
		},
	}

	json, _ := json.Marshal(item)
	return string(json)
}

func makeItems(items ...string) string {
	return fmt.Sprintf("[%s]", strings.Join(items, ","))
}
