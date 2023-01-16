package grpc

import (
	"context"
	"errors"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
	"strings"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	metric2 "github.com/dayterr/go_agent_metrics/internal/metric"
	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
)

func (gs GRPCServer) IPInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if gs.ts == nil {
		return handler(ctx, req)
	}

	addr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("parsing context error")
	}

	clientIP := strings.Split(addr.Addr.String(), ":")[0]

	if !gs.ts.Contains(net.ParseIP(clientIP)) {
		return nil, errors.New("ip not allowed")
	}

	return handler(ctx, req)

}

func (gs GRPCServer) PostMetric(ctx context.Context, req *pb.PostMetricRequest) (*pb.PostMetricResponse, error) {
	var resp pb.PostMetricResponse

	var m metric2.Metrics
	m.ID = req.Metrics.Id
	m.MType = req.Metrics.Type.String()
	m.Delta = &req.Metrics.Delta
	m.Value = &req.Metrics.Value

	if gs.key != "" {
		hashCheck := hash.EncryptMetric(m, gs.key)
		if !hash.CheckHash(m, hashCheck) {
			return &resp, errors.New("hash check failed")
		}
	}

	switch req.Metrics.Type.String() {
	case agent.GaugeType:
		gs.storage.SetGuage(ctx, req.Metrics.Id, &req.Metrics.Value)
		return &resp, nil
	case agent.CounterType:
		gs.storage.SetCounter(ctx, req.Metrics.Id, &req.Metrics.Delta)
	default:
		return &resp, errors.New("unknown type found")
	}

	return &resp, nil
}

func (gs GRPCServer) GetMetric(ctx context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var resp pb.GetMetricResponse

	var m metric2.Metrics
	m.ID = req.Name
	m.MType = req.Type

	switch req.Type {
	case agent.GaugeType:
		if gs.storage.CheckGaugeByName(ctx, req.Name) {
			v, err := gs.storage.GetGuageByID(ctx, req.Name)
			if err != nil {
				return &resp, err
			}
			resp.Metric.Value = v
			m.Value = &v
			if gs.key != "" {
				resp.Metric.Hash = hash.EncryptMetric(m, gs.key)
			}
			return &resp, nil
		}
	case agent.CounterType:
		if gs.storage.CheckCounterByName(ctx, req.Name) {
			v, err := gs.storage.GetCounterByID(ctx, req.Name)
			if err != nil {
				return &resp, err
			}
			resp.Metric.Delta = v
			m.Delta = &v
			if gs.key != "" {
				resp.Metric.Hash = hash.EncryptMetric(m, gs.key)
			}
			return &resp, nil
		}
	default:
		return &resp, errors.New("unknown type found")
	}

	if !gs.storage.CheckCounterByName(ctx, req.Name) && !gs.storage.CheckGaugeByName(ctx, req.Name) {
		return &resp, errors.New("metric not found")
	}

	return &resp, nil
}

