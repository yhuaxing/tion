package tion

import (
	"net/http"
)

// 定义动作
type HandleFunc func(*Context)

// 定义引擎
type Engine struct {
	*RouterGroup
	groups []*RouterGroup
}

// 初始化引擎
func New() *Engine {
	router := newRouter()
	engine := &Engine{
		RouterGroup: &RouterGroup{
			prefix: "",
			router: router,
		},
	}
	engine.RouterGroup.engine = engine
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	node, params := e.router.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
	}
	c.middlewares = e.getMiddlewares(node)
	c.Next()
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}
