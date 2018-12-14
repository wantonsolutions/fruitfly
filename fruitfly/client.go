package fruitfly

import (
    "time"
    "log"
)

func Client(ip, port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    for true {
        log.Printf("Sending to at %s:%s",ip,port)
        time.Sleep(time.Second)
    }
}
