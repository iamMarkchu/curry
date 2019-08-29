package main

import (
	"curry/router"
	"log"
)

func main() {
	r := router.NewRouter()
	log.Fatal(r.Run(":8888"))
}
