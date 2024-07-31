package twix

import (
	"net/http"
)

// Context holds information about the current request and response
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
}

// Param retrieves a URL parameter value by key
func (c *Context) Param(key string) string {
	if value, ok := c.Params[key]; ok {
		return value
	}
	return ""
}
