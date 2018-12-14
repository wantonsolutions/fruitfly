package fruitfly

import (
    "time"
    "log"
)

func Server(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    for true {
        log.Printf("Serving at :%s",port)
        time.Sleep(time.Second)
    }
}
