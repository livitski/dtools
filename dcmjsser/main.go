package main

import "log"
import "runtime"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	djs := DJsServ{}
	log.Println("info: gui located at address http://your ip:9978/index.html, for recive study use scsc_port=50000, aetitle=AE_DTLS")
	if err := djs.Start(9978); err != nil {
		log.Println(err)
	}

}
