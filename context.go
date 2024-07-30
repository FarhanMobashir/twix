package twix

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Params         map[string]string
}

func (c *Context) Param(key string) string {
	if value, ok := c.Params[key]; ok {
		return value
	}
	return ""
}

// Status sets the HTTP status code and returns the context for chaining
func (c *Context) Status(code int) *Context {
	c.ResponseWriter.WriteHeader(code)
	return c
}

// String writes a plain text response and returns the context for chaining
func (c *Context) String(response string) *Context {
	c.ResponseWriter.Write([]byte(response))
	return c
}

func (c *Context) JSON(code int, response interface{}) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	if err := json.NewEncoder(c.ResponseWriter).Encode(response); err != nil {
		http.Error(c.ResponseWriter, "Internal Server Error", http.StatusInternalServerError)
	}
}
