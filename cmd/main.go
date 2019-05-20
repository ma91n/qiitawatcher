package cmd

import (
	"github.com/laqiiz/qiitawatcher/controller"
	"log"
)

// main endpoint
func main() {
	log.Printf("start")
	if err := controller.Execute(); err != nil {
		log.Fatalf("error=%v", err)
	}
	log.Printf("success")
}
