package main

import (
	"fmt"
	"log"
	"net/http"
	heartbeat "user/helper/Heartbeat"
	"user/router"
)

func main() {
	r := router.Router()
	fmt.Println("Server running in port 3000")
	go heartbeat.Beat()
	log.Fatal(http.ListenAndServe(":3000", r))
	return
}
