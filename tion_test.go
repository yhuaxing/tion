package tion

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	r := New()
	r.Get("/hello", func(c *Context) {
		c.String(http.StatusOK, c.Path)
	})
	r.Get("/say/:text", func(c *Context) {
		c.String(http.StatusOK, c.Param("text"))
	})
	v1 := r.Group("/v1")
	{
		v1.Use(func(c *Context) {
			log.Println("startTime:", time.Now())
			c.Next()
			log.Println("endTime:", time.Now())
		})
		v1.Get("/say/:text", func(c *Context) {
			log.Println("HELLO")
			c.String(http.StatusOK, c.Param("text"))
		})
		vp := v1.Group("/p/:page")
		vp.Use(func(c *Context) {
			log.Println("startTime2:", time.Now())
			c.Next()
			log.Println("endTime2:", time.Now())
		})
		vp.Get("/name", func(c *Context) {
			c.String(http.StatusOK, c.Param("page"))
		})
	}
	go func() {
		log.Println("启动服务")
		err := r.Run(":8080")
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(5 * time.Second)
	log.Println("发起请求")
	resp, err := http.Get("http://localhost:8080/v1/p/1/name")
	if err != nil {
		t.Error("HTTP请求失败")
	}
	defer resp.Body.Close()
	if bs, err := ioutil.ReadAll(resp.Body); err != nil || string(bs) != "1" {
		t.Error("接口返回不期望的值")
	}
	log.Println("Done!")
}
