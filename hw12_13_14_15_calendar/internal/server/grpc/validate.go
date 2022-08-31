package grpc

import (
	"context"
	"errors"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator func(req interface{}) error

func UnaryServerRequestValidatorInterceptor(validator Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validator(req); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "%s is rejected by validate middleware. Error: %v", info.FullMethod, err)
		}
		return handler(ctx, req)
	}
}

func ValidateReq(req interface{}) error {
	switch r := req.(type) {
	case *pb.AddEventRequest:
		if r.Event.Title == "" || r.Event.UserId == "" {
			return errors.New("middleware validator: you need Title and UserID for a new event")
		}
	case *pb.UpdateEventRequest:
		if r.EventId == "" {
			return errors.New("middleware validator: provide EventID, which you want update")
		}
		if r.UpdatedEvent.UserId == "" || r.UpdatedEvent.Title == "" {
			return errors.New("middleware validator: updating event has to have Title and UserID")
		}
	case *pb.DeleteEventRequest:
		if r.EventId == "" {
			return errors.New("middleware validator: provide EventID, which you want delete")
		}
	}

	return nil
}
