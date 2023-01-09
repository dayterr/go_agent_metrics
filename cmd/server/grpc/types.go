package grpc

import (
	"log"

	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServiceServer
	storage storage.Storager
}

func NewGRPCServer(dsn string, isDB bool) (GRPCServer, error) {
	var s storage.Storager
	var err error
	if isDB {
		s, err = storage.NewDB(dsn)
		if err != nil {
			log.Println(err)
			return GRPCServer{}, err
		}
	} else {
		s = storage.NewIMS()
	}
	gs := GRPCServer{storage: s}
	return gs, nil
}