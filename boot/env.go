package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/glog"
	"github.com/joho/godotenv"
)

func loadEnv() {
	var err error
	err = godotenv.Load()
	if err != nil {
		glog.Fatal("Error loading .env file")
	}

	// 动态修改配置
	config := g.Config()

	dsn := genv.Get("DATABASE_DSN")
	err = config.Set("database.link", dsn)
	if err != nil {
		glog.Fatal("Error set config")
	}
}
