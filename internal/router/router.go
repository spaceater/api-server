package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// use example:
// addr("/user",nihao).use(login)
// addrGroup("/user",
//		addr("/login",nihao)
//)

// 支持分组路由以及中间间，分组路由可直接全局添加中间价，提供默认中间件
// 注册只梳理关系，最后init挂数

var routerMap []*Router

type Router struct {
	path       string
	middleware []http.HandlerFunc
}

type RouterGroup []*Router

func (rt *Router) next(w http.ResponseWriter, r *http.Request) {
	for index := len(rt.middleware) - 1; index > -1; index-- {
		_, exists := r.Header[http.CanonicalHeaderKey("middleware-break")]
		if exists {
			return
		}
		rt.middleware[index](w, r)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("err req", r)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			log.Println("pen server error:", string(buf))
			result := struct {
				Msg     string      `json:"msg"`
				Code    int         `json:"code"`
				Data    interface{} `json:"data"`
				Success bool        `json:"success"`
			}{}
			result.Success = false
			result.Msg = fmt.Sprintf("%s", r)
			res, _ := json.Marshal(result)
			w.Write(res)
		}
	}()
}

func Url(path string, handlerFunc http.HandlerFunc) *Router {
	router := &Router{
		path:       path,
		middleware: []http.HandlerFunc{handlerFunc},
	}
	routerMap = append(routerMap, router)
	return router
}

func (r *Router) Use(hf http.HandlerFunc) *Router {
	r.middleware = append(r.middleware, hf)
	return r
}

func UrlGroup(pre string, routers ...*Router) *RouterGroup {
	var AddrList RouterGroup
	for i := range routers {
		r := routers[i]
		r.path = pre + r.path
		AddrList = append(AddrList, r)
	}
	return &AddrList
}

func (rg RouterGroup) Use(hf http.HandlerFunc) RouterGroup {
	for i := range rg {
		j := rg[i]
		(*j).middleware = append((*j).middleware, hf)
	}
	return rg
}

func SetBreak(r *http.Request) {
	r.Header.Set("middleware-break", "")
}

// 输出生成的路由信息
func Debug() {
	log.Println("----路由信息-----")
	routers := FilterMultiple()
	for i := range routers {
		r := routers[i]
		log.Println("注册路径:", r.path)
	}
	log.Println("----路由调试end----")
}

// 过滤重复路由
func FilterMultiple() map[string]*Router {
	routers := make(map[string]*Router)
	for i := range routerMap {
		r := routerMap[i]
		routers[r.path] = r
	}
	return routers
}

func Init() {
	routers := FilterMultiple()
	for i := range routers {
		r := routers[i]
		http.HandleFunc(r.path, r.next)
	}
}
