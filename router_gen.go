// Code generated by hertz generator. DO NOT EDIT.
// @author fuxingyuan
package main

import (
	router "TikTok/biz/router"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// register registers all routers.
func register(r *server.Hertz) {

	router.GeneratedRegister(r)

	customizedRegister(r)
}
