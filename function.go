package qiitawatcher

import (
	"net/http"
	"log"
	"time"
	"os"

	"github.com/laqiiz/qiitawatcher/controller"
)

// http endpoint for Google Cloud Function
func Receive(w http.ResponseWriter, r *http.Request) {
	log.Printf("request received")

	// override created time
	beforeDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	os.Setenv("CREATED", beforeDate)

	if err := controller.Execute(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("success"))
	log.Printf("finished")
}
