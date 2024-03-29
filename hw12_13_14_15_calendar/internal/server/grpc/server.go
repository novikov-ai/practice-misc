package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/models"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	pb "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/pb"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedCalendarServer
	database app.Storage
}

func Start(ctx context.Context, st app.Storage, logger app.Logger, config configs.Config) error {
	lsn, err := net.Listen("tcp", ":"+config.GetServerConfig().PortGRPC)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryServerRequestValidatorInterceptor(ValidateReq),
			UnaryServerLoggingInterceptor(logger)),
	)

	if err := st.Connect(ctx); err != nil {
		return err
	}

	service := new(Service)
	service.database = st

	pb.RegisterCalendarServer(server, service)

	logger.Info(fmt.Sprintf("starting protobuf server on %s", lsn.Addr().String()))

	if err = server.Serve(lsn); err != nil {
		return err
	}

	logger.Info("listening...")

	<-ctx.Done()
	return nil
}

func (s *Service) AddEvent(ctx context.Context, request *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	createdID, err := s.database.Add(ctx, convertToEvent(request.Event))
	return &pb.AddEventResponse{CreatedId: createdID}, err
}

func (s *Service) UpdateEvent(ctx context.Context, request *pb.UpdateEventRequest) (*pb.EventResponse, error) {
	err := s.database.Update(ctx, request.EventId, convertToEvent(request.UpdatedEvent))
	if err != nil {
		return &pb.EventResponse{Status: pb.Status_STATUS_FAILED, Message: err.Error()}, err
	}

	return &pb.EventResponse{Status: pb.Status_STATUS_SUCCESS, Message: "event was updated"}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, request *pb.DeleteEventRequest) (*pb.EventResponse, error) {
	err := s.database.Delete(ctx, request.EventId)
	if err != nil {
		return &pb.EventResponse{Status: pb.Status_STATUS_FAILED, Message: "event was not deleted"}, err
	}

	return &pb.EventResponse{Status: pb.Status_STATUS_SUCCESS, Message: "event was successfully deleted"}, nil
}

func (s *Service) GetEventsForDay(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, s.database.GetEventsForDay)
}

func (s *Service) GetEventsForWeek(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, s.database.GetEventsForWeek)
}

func (s *Service) GetEventsForMonth(ctx context.Context, request *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return getEvents(ctx, request, s.database.GetEventsForMonth)
}

func getEvents(ctx context.Context, request *pb.GetEventsRequest, action func(ctx context.Context, time time.Time) []models.Event) (*pb.GetEventsResponse, error) {
	pbEvents := make([]*pb.Event, 0)

	events := action(ctx, request.FromDay.AsTime())
	for _, ev := range events {
		pbEvents = append(pbEvents, convertToPbEvent(ev))
	}

	return &pb.GetEventsResponse{Events: pbEvents}, nil
}
