package grpc

import (
	"context"
	"fmt"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	pb "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/pb"
	"google.golang.org/grpc"
)

func UnaryServerLoggingInterceptor(logger app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch r := req.(type) {
		case *pb.AddEventRequest:
			logger.Info(fmt.Sprintf("event with ID:%s was created", r.Event.Id))
		case *pb.UpdateEventRequest:
			logger.Info(fmt.Sprintf("event with ID:%s was updated", r.EventId))
		case *pb.DeleteEventRequest:
			logger.Info(fmt.Sprintf("event with ID:%s was deleted", r.EventId))
		case *pb.GetEventsRequest:
			logger.Info(fmt.Sprintf("got events for timespan %s", r.FromDay))
		}

		return handler(ctx, req)
	}
}
