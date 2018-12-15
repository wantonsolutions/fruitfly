package fruitfly

import (
    //"time"
    "log"
)

func Server(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    /*
    err, addr := net.ResolveUDPAddr("udp", "localhost:port")
    if err != nil {
        log.Fatal(err)
    }
    for true {
        conn, err := net.ListenUDP("udp",addr)
        if err != nil {
            log.Fatal(err)
        }
        handelConn(conn)
        log.Printf("Serving at :%s",port)
        time.Sleep(time.Second)
    }*/
}
