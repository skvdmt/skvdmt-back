package internal

import "net/http"

// Delivery transport application interface
type Delivery interface {
	Router() http.Handler
	Close() error
}
