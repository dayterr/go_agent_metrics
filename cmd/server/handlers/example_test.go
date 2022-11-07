package handlers

import (
	"context"
	"log"

	//"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

func Example() {
	ctx := context.Background()
	handler := NewAsyncHandler("", "", false)
	stor := storage.NewIMS()
	var v1 float64 = 353808
	var v2 float64 = 3865
	stor.SetGuage(ctx, "Alloc", &v1)
	stor.SetGuage(ctx, "BuckHashSys", &v2)

	_, err := handler.MarshallMetrics()
	if err != nil {
		log.Println("something went wrong")
	}
}
