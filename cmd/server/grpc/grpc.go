package grpc

import (
	"context"
	"errors"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
)

func (gs *GRPCServer) PostMetric(ctx context.Context, req *pb.PostMetricRequest) (*pb.PostMetricResponse, error) {
	var resp pb.PostMetricResponse

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

func (gs *GRPCServer) GetMetric(ctx context.Context, req *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var resp pb.GetMetricResponse

	switch req.Type {
	case agent.GaugeType:
		if gs.storage.CheckGaugeByName(ctx, req.Name) {
			v, err := gs.storage.GetGuageByID(ctx, req.Name)
			if err != nil {
				return &resp, err
			}
			resp.Metric.Value = v
			return &resp, nil
		}
	case agent.CounterType:
		if gs.storage.CheckCounterByName(ctx, req.Name) {
			v, err := gs.storage.GetCounterByID(ctx, req.Name)
			if err != nil {
				return &resp, err
			}
			resp.Metric.Delta = v
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

