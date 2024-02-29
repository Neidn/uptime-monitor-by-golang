package lib

import "time"

func Delay(sec int) {
	time.Sleep(time.Duration(sec) * time.Second)
}
