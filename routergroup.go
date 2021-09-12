package tion

import (
	"net/http"
	"sort"
	"strings"
)

type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	engine      *Engine
	router      *router
	middlewares []HandleFunc
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		parent: rg,
		engine: rg.engine,
		router: rg.router,
	}
	newGroup.engine.groups = append(newGroup.engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) Use(middlewares ...HandleFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

func (rg *RouterGroup) addRoute(method string, comp string, handler HandleFunc) {
	pattern := rg.prefix + comp
	rg.router.addRoute(method, pattern, handler)
}

func (rg *RouterGroup) Get(pattern string, handler HandleFunc) {
	rg.addRoute("GET", pattern, handler)
}

func (rg *RouterGroup) Post(pattern string, handler HandleFunc) {
	rg.addRoute("POST", pattern, handler)
}

func (rg *RouterGroup) Put(pattern string, handler HandleFunc) {
	rg.addRoute("PUT", pattern, handler)
}

func (rg *RouterGroup) Delete(pattern string, handler HandleFunc) {
	rg.addRoute("POST", pattern, handler)
}

func (rg *RouterGroup) getMiddlewares(node *node) []HandleFunc {
	if node == nil {
		return []HandleFunc{func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		}}
	}
	firstFilter := []*RouterGroup{}
	for _, v := range rg.engine.groups {
		if strings.HasPrefix(node.pattern, v.prefix+"/") {
			firstFilter = append(firstFilter, v)
		}
	}
	sort.Slice(firstFilter, func(i, j int) bool {
		return len(firstFilter[i].prefix) < len(firstFilter[j].prefix)
	})
	middlewares := []HandleFunc{}
	for _, item := range firstFilter {
		middlewares = append(middlewares, item.middlewares...)
	}
	middlewares = append(middlewares, node.handler)
	return middlewares
}
