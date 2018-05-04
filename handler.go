package dialogflow

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// Handler provides a simple way to route and handle requests by intent
// This handler is based on Campoy's apiai project: github.com/campoy/apiai/blob/master/apiai.go
type Handler struct {
	mu       sync.RWMutex
	handlers map[string]IntentHandler
}

// NewHandler returns a new empty handler
func NewHandler() *Handler {
	return &Handler{handlers: make(map[string]IntentHandler)}
}

// Register registers the handler for a given intent
func (h *Handler) Register(intent string, handler IntentHandler) {
	h.mu.Lock()
	h.handlers[intent] = handler
	h.mu.Unlock()
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var dfr *Request

	if err := json.NewDecoder(req.Body).Decode(&dfr); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	intent := dfr.QueryResult.Intent.DisplayName
	h.mu.RLock()
	handler, ok := h.handlers[intent]
	h.mu.RUnlock()
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(context.Background(), httpRequestKey, req)
	dff, status := handler(ctx, dfr)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	if err := json.NewEncoder(rw).Encode(dff); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

// IntentHandler handles an intent, by returning a Fulfillment object and an HTTP status code.
type IntentHandler func(ctx context.Context, dfr *Request) (*Fulfillment, int)

type key string

const httpRequestKey key = "httprequest"

// HTTPRequest returns the HTTP request associated to the given context or nil.
func HTTPRequest(ctx context.Context) *http.Request {
	req, ok := ctx.Value(httpRequestKey).(*http.Request)
	if !ok {
		return nil
	}
	return req
}
