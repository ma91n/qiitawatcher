package qiitawatcher

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laqiiz/qiitawatcher/controller"
)

// http endpoint for Google Cloud Function
func Receive(w http.ResponseWriter, r *http.Request) {
	log.Printf("request received")

	// override created time
	beforeDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if err := os.Setenv("CREATED", beforeDate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if err := controller.Execute(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte("success"))
	log.Printf("finished")
}
