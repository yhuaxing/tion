package tion

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Rw          http.ResponseWriter
	Req         *http.Request
	Path        string
	Method      string
	StatusCode  int
	Params      map[string]string
	index       int
	middlewares []HandleFunc
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {

	return &Context{
		Rw:     w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.Form.Get(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Rw.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Rw.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Rw.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	if err := json.NewEncoder(c.Rw).Encode(obj); err != nil {
		http.Error(c.Rw, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Rw.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Rw.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.middlewares) {
		c.middlewares[c.index](c)
	}
}
