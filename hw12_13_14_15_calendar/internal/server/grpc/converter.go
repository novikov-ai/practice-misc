package grpc

import (
	pb "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/pb"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/models"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertToEvent(event *pb.Event) models.Event {
	return models.Event{
		ID:             event.Id,
		Title:          event.Title,
		Description:    event.Description,
		DateTime:       event.DateTime.AsTime(),
		Duration:       event.Duration.AsDuration(),
		NotifiedBefore: event.NotifiedBefore.AsDuration(),
	}
}

func convertToPbEvent(event models.Event) *pb.Event {
	return &pb.Event{
		Id:             event.ID,
		Title:          event.Title,
		Description:    event.Description,
		DateTime:       timestamppb.New(event.DateTime),
		Duration:       durationpb.New(event.Duration),
		NotifiedBefore: durationpb.New(event.NotifiedBefore),
	}
}
