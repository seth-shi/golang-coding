package main

import (
	_ "golang-coding/boot"
	_ "golang-coding/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
