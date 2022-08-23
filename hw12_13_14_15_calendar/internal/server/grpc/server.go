package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/server"
	"github.com/r-egorov/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type Service struct {
	UnimplementedEventsServer
	app        server.Application
	logg       server.Logger
	host, port string
	serv       *grpc.Server
}

func NewService(logger server.Logger, app server.Application, host, port string) *Service {
	service := &Service{
		app:  app,
		logg: logger,
		host: host,
		port: port,
	}

	grpcServer := grpc.NewServer(nil)
	RegisterEventsServer(grpcServer, service)

	service.serv = grpcServer
	return service
}

func (s *Service) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	errChan := make(chan error)
	go func() {
		if err := s.serv.Serve(lsn); err != nil {
			errChan <- err
		}
	}()

	s.logg.Info(fmt.Sprintf("gRPC serving at %s:%s", s.host, s.port))
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
	}
	return nil
}

func (s *Service) Stop(ctx context.Context) {
	s.logg.Info("stopping gRPC server...")
	s.serv.GracefulStop()
}

func (s *Service) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	nativeEvent, err := mapEventProtoToNative(req.Event)
	if err != nil {
		return nil, err
	}
	savedEvent, err := s.app.SaveEvent(ctx, *nativeEvent)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{
		Event: mapEventNativeToProto(&savedEvent),
	}, nil
}

func (s *Service) Update(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	nativeEvent, err := mapEventProtoToNative(req.Event)
	if err != nil {
		return nil, err
	}
	updatedEvent, err := s.app.UpdateEvent(ctx, *nativeEvent)
	if err != nil {
		return nil, err
	}
	return &UpdateResponse{
		Event: mapEventNativeToProto(&updatedEvent),
	}, nil
}

func (s *Service) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	nativeEvent, err := mapEventProtoToNative(req.Event)
	if err != nil {
		return nil, err
	}
	err = s.app.DeleteEvent(ctx, nativeEvent.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteResponse{
		Id: req.Event.Id,
	}, nil
}

func (s *Service) ListDay(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	events, err := s.app.ListEventsDay(ctx, req.PeriodStart.AsTime())
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Events: mapEventsNativeToProto(events),
	}, nil
}

func (s *Service) ListWeek(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	events, err := s.app.ListEventsWeek(ctx, req.PeriodStart.AsTime())
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Events: mapEventsNativeToProto(events),
	}, nil
}

func (s *Service) ListMonth(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	events, err := s.app.ListEventsMonth(ctx, req.PeriodStart.AsTime())
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Events: mapEventsNativeToProto(events),
	}, nil
}

func mapEventProtoToNative(event *Event) (*storage.Event, error) {
	eventID, err := uuid.FromBytes([]byte(event.GetId()))
	if err != nil {
		return nil, err
	}
	ownerID, err := uuid.FromBytes([]byte(event.GetOwnerId()))
	if err != nil {
		return nil, err
	}
	return &storage.Event{
		ID:          eventID,
		Title:       event.Title,
		DateTime:    event.Datetime.AsTime(),
		Duration:    event.Duration.AsDuration(),
		Description: event.Description,
		OwnerID:     ownerID,
	}, nil
}

func mapEventNativeToProto(event *storage.Event) *Event {
	return &Event{
		Id:          event.ID.String(),
		Title:       event.Title,
		Datetime:    timestamppb.New(event.DateTime),
		Duration:    durationpb.New(event.Duration),
		Description: event.Description,
		OwnerId:     event.OwnerID.String(),
	}
}

func mapEventsNativeToProto(events []storage.Event) []*Event {
	eventsProto := make([]*Event, len(events))
	for _, event := range events {
		eventsProto = append(eventsProto, mapEventNativeToProto(&event))
	}
	return eventsProto
}
