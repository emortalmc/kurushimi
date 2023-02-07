package frontend

import (
	"context"
	"github.com/google/uuid"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"kurushimi/internal/notifier"
	"kurushimi/internal/statestore"
	"kurushimi/pkg/pb"
	"net"
	"time"
)

var (
	logger *zap.SugaredLogger
)

type frontendService struct {
	pb.UnimplementedFrontendServer

	stateStore statestore.StateStore
	notifier   notifier.Notifier
}

func Init(ctx context.Context, stateStore statestore.StateStore, notifier notifier.Notifier) {
	logger = zap.S()

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpczap.UnaryServerInterceptor(logger.Desugar(), grpczap.WithLevels(func(code codes.Code) zapcore.Level {
			if code != codes.Internal && code != codes.Unavailable && code != codes.Unknown {
				return zapcore.DebugLevel
			} else {
				return zapcore.ErrorLevel
			}
		})),
	))

	pb.RegisterFrontendServer(s, &frontendService{
		stateStore: stateStore,
		notifier:   notifier,
	})

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
}

func (s *frontendService) CreateTicket(ctx context.Context, req *pb.CreateTicketRequest) (*pb.Ticket, error) {
	ticket := req.Ticket
	if ticket == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ticket is required")
	}
	if ticket.CreatedAt != nil {
		return nil, status.Errorf(codes.InvalidArgument, "CreatedAt is not allowed")
	}
	ticket.Id = uuid.New().String()
	ticket.CreatedAt = timestamppb.New(time.Now())
	err := s.stateStore.CreateTicket(ctx, req.Ticket)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create ticket %v", err)
	}
	err = s.stateStore.IndexTicket(ctx, req.Ticket)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to index ticket %v", err)
	}
	return req.Ticket, nil
}

func (s *frontendService) DeleteTicket(ctx context.Context, req *pb.DeleteTicketRequest) (*emptypb.Empty, error) {
	err := s.stateStore.DeleteTicket(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete ticket %v", err)
	}
	err = s.stateStore.UnIndexTicket(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to unindex ticket %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *frontendService) GetTicket(ctx context.Context, req *pb.GetTicketRequest) (*pb.Ticket, error) {
	ticket, err := s.stateStore.GetTicket(ctx, req.TicketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get ticket %v", err)
	}
	return ticket, nil
}

func (s *frontendService) WatchTicketCountdown(req *pb.WatchCountdownRequest, stream pb.Frontend_WatchTicketCountdownServer) error {
	finishNotifier := make(chan struct{})
	s.notifier.AddCountdownListener(req.TicketId, stream, finishNotifier)

	<-finishNotifier
	s.notifier.RemoveCountdownListener(req.TicketId)
	return nil
}

func (s *frontendService) WatchTicketAssignment(req *pb.WatchAssignmentRequest, stream pb.Frontend_WatchTicketAssignmentServer) error {
	finishNotifier := make(chan struct{})
	s.notifier.AddAssignmentListener(req.TicketId, stream, finishNotifier)

	<-finishNotifier
	s.notifier.RemoveAssignmentListener(req.TicketId)
	return nil
}
