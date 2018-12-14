package main

import (
	"github.com/wantonsolutions/fruitfly/fruitfly"
    "log"
    "flag"
)

var server = flag.Bool("server", false, "Launch in server mode, must set ip field")
var client = flag.Bool("client", false, "Launch in client mode, must set serverIP feild")
var process = flag.Bool("process", false, "Launch in process mode")

//client specific args
var serverIP = flag.String("serverIP", "", "Ip of the server")
var serverPort = flag.String("serverPort", "", "port of the server")

//server specific args

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    flag.Parse()
    checkExclusion()
    if *server {
        checkServerArgs()
        fruitfly.Server(*serverPort)
    } else if *client {
        checkClientArgs()
        fruitfly.Client(*serverIP,*serverPort)
    } else if *process {
        checkProcessArgs()
        fruitfly.Process()
    } else {
        flag.PrintDefaults()
        log.Fatal("I'm sensing your new to fruitfly")
    }
}

func checkClientArgs() {
    //potentially check that the client args were set
    return
}

func checkServerArgs() {
    //potentially check that the server args were set
    return
}

func checkProcessArgs() {
    //potentailly check that the processing args were set
    return
}

func checkExclusion() {
    fatal := false
    if *client {
        if *server || *process {
            fatal = true
        }
    } else if *server {
        if *client || *process {
            fatal = true
        }
    } else if *process {
        if *client || *server {
            fatal = true
        }
    }
    if fatal {
        log.Fatal("Chose only one operational mode to run fruitfly in")
    }
}
            
