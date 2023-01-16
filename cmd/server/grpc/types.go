package grpc

import (
	"log"
	"net"

	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type GRPCServer struct {
	pb.UnimplementedMetricsServiceServer
	storage storage.Storager
	key string
	ts *net.IPNet
}

func NewGRPCServer(dsn string, key, ts string) (GRPCServer, error) {
	var s storage.Storager
	var err error
	if dsn != "" {
		s, err = storage.NewDB(dsn)
		if err != nil {
			log.Println(err)
			return GRPCServer{}, err
		}
	} else {
		s = storage.NewIMS()
	}
	gs := GRPCServer{storage: s, key: key}
	if ts != "" {
		_, cidr, err := net.ParseCIDR(ts)
		if err != nil {
			return GRPCServer{}, err
		}
		gs.ts = cidr
	}
	return gs, nil
}