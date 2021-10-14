package boot

import (
	"github.com/gogf/gf/os/glog"
	"golang-coding/app/service"
	"time"
)

func bootTime() {
	service.BootTime = time.Now().Format("2006-01-02 15:04:05")
	glog.Info("application boot time :" + service.BootTime)
}
