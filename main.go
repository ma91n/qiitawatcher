package main

import (
	"log"
	"qiitawatcher/controller"
)

// main endpoint
func main() {
	log.Printf("start")
	if err := controller.Execute(); err != nil {
		log.Fatalf("error=%v", err)
	}
	log.Printf("success")
}
