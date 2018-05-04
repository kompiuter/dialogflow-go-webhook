package dialogflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func prepare(intent string) map[string]interface{} {
	return map[string]interface{}{
		"queryResult": map[string]interface{}{
			"intent": map[string]string{
				"displayName": intent,
			},
		},
	}
}

func TestHandler_RouteIntent(t *testing.T) {
	const intent = "books.get"

	input := prepare(intent)
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewReader(b))

	rr := httptest.NewRecorder()
	handler := NewHandler()
	handler.Register(intent, func(ctx context.Context, dfr *Request) (*Fulfillment, int) {
		return &Fulfillment{
			FulfillmentText: "hello",
		}, http.StatusOK
	})

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"fulfillmentText":"hello"}`
	got := strings.TrimSpace(rr.Body.String())
	if got != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestHandler_MultipleIntents(t *testing.T) {
	tests := []struct {
		intent string
		status int
	}{
		{
			intent: "books.get",
			status: 200,
		},
		{
			intent: "books.set",
			status: 400,
		},
	}

	handler := NewHandler()
	handler.Register(tests[0].intent, func(ctx context.Context, dfr *Request) (*Fulfillment, int) {
		return &Fulfillment{FulfillmentText: tests[0].intent}, tests[0].status
	})
	handler.Register(tests[1].intent, func(ctx context.Context, dfr *Request) (*Fulfillment, int) {
		return &Fulfillment{FulfillmentText: tests[1].intent}, tests[1].status
	})

	for _, tc := range tests {
		t.Run(tc.intent, func(t *testing.T) {
			b, err := json.Marshal(prepare(tc.intent))
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.status {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.status)
			}
		})

	}
}
