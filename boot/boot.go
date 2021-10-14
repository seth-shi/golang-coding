package boot

import (
	_ "golang-coding/packed"
)

func init() {

	bootTime()
	loadEnv()
}
