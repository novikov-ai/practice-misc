package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	//"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/grpclog"
)

var database app.Storage

type Service struct {
	pb.UnimplementedCalendarServer
}

func Start(ctx context.Context, st app.Storage, logger app.Logger) error {
	database = st
	if err := database.Connect(ctx); err != nil {
		return err
	}

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryServerRequestValidatorInterceptor(ValidateReq),
			UnaryServerLoggingInterceptor(logger)),
	)

	pb.RegisterCalendarServer(server, new(Service))

	logger.Info(fmt.Sprintf("starting protobuf server on %s\n", lsn.Addr().String()))

	if err = server.Serve(lsn); err != nil {
		return err
	}

	logger.Info("listening...")

	<-ctx.Done()
	return nil
}

func (s *Service) AddEvent(ctx context.Context, request *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	select {
	case <-ctx.Done():
		break
	default:
		createdID, err := database.Add(ctx, convertToEvent(request.Event))
		return &pb.AddEventResponse{CreatedId: createdID}, err
	}

	return &pb.AddEventResponse{}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, request *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	select {
	case <-ctx.Done():
		break
	default:
		err := database.Update(ctx, request.EventId, convertToEvent(request.UpdatedEvent))
		if err != nil {
			return &pb.EventResponse{Status: pb.Status_STATUS_FAILED, Message: err.Error()}, err
		}
	}

	return &pb.EventResponse{Status: pb.Status_STATUS_SUCCESS, Message: "event was updated"}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, request *pb.DeleteEventRequest) (*pb.EventResponse, error) {
	select {
	case <-ctx.Done():
		break
	default:
		err := database.Delete(ctx, request.EventId)
		if err != nil {
			return &pb.EventResponse{Status: pb.Status_STATUS_FAILED, Message: "event was not deleted"}, err
		}
	}

	return &pb.EventResponse{Status: pb.Status_STATUS_SUCCESS, Message: "event was successfully deleted"}, nil
}

func (s *Service) GetEventsForDay(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, database.GetEventsForDay)
}

func (s *Service) GetEventsForWeek(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, database.GetEventsForWeek)
}

func (s *Service) GetEventsForMonth(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, database.GetEventsForMonth)
}

func getEvents(ctx context.Context, request *pb.GetEventsRequest, action func(ctx context.Context, time time.Time) []models.Event) (*pb.GetEventsResponse, error) {
	pbEvents := make([]*pb.Event, 0)

	select {
	case <-ctx.Done():
		break
	default:
		events := action(ctx, request.FromDay.AsTime())
		for _, ev := range events {
			pbEvents = append(pbEvents, convertToPbEvent(ev))
		}
	}

	return &pb.GetEventsResponse{Events: pbEvents}, nil
}
