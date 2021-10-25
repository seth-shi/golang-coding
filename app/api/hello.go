package api

import (
	"github.com/gogf/gf/net/ghttp"
	"golang-coding/app/service"
)

var Hello = helloApi{}

type helloApi struct {}


func (*helloApi) Ping(r *ghttp.Request) {
	r.Response.Writeln("PONG")
}


func (*helloApi) Index(r *ghttp.Request) {
	r.Response.Writeln("Hello World! Application Start Time :" + service.BootTime)
}
