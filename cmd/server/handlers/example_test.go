package handlers

import (
	"log"

	//"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

func Example() {
	handler := NewAsyncHandler("", "", false)
	stor := storage.NewIMS()
	var v1 float64
	v1 = 353808
	var v2 float64
	v2 = 3865
	stor.SetGuage("Alloc", &v1)
	stor.SetGuage("BuckHashSys", &v2)

	_, err := handler.MarshallMetrics()
	if err != nil {
		log.Println("something went wrong")
	}
}
