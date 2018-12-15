package fruitfly

import (
    "log"
    "net"
)

const PHOTO_COUNT = 3
const MODEL_FILE = "models/crap.gob"
const PHOTO_FILE = "data/Aging_Study_1/dec_7_2018/1.jpg"


func Client(ip, port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    //localProcessing()
    RemoteSynchronous()
}

func localProcessing() {
    //full client side processing
    //read photo in from a file (add a stub for actually execing it taking it)
    //process the photo
    histogram := ReadModel(MODEL_FILE)
    for i:=0;i<PHOTO_COUNT;i++ {
        im := openImage(PHOTO_FILE)
        bounds := im.Bounds()
        score := calculateScore(bounds,histogram,im)
        rypeness, _, cs := processRypenessRGB(score,bounds,im)
        log.Printf("%s, %s\n",rypeness,cs.String())
    }
}

func RemoteSynchronous() {

}

    
