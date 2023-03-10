// Code generated by hertz generator.
// @author fuxingyuan
// 此文件只可以修改hp!!!,其余修改请联系author
package main

import (
	"TikTok/biz/controller"
	"TikTok/biz/dao"
	mw "TikTok/biz/mw/jwt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/pprof"
)

func main() {
	h := server.Default(server.WithHostPorts("192.168.137.1:8081"), server.WithMaxRequestBodySize(1024*1024*1024))
	dao.Init()
	mw.Initjwt()
	register(h)
	controller.TimeTaskExec()
	pprof.Register(h)
	h.Spin()
}
