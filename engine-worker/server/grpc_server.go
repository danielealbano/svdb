package server

import (
	"github.com/danielealbano/svdb/engine-worker/collection"
	shared_support "github.com/danielealbano/svdb/shared/support"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type GrpcServer struct {
	listener   *net.Listener
	grpcServer *grpc.Server
	collection *collection.Collection
	running    *sync.WaitGroup
	done       chan struct{}
}

func NewGrpcServer(listener *net.Listener, coll *collection.Collection) *GrpcServer {
	grpcServer := grpc.NewServer()

	return &GrpcServer{
		listener:   listener,
		grpcServer: grpcServer,
		collection: coll,
		done:       make(chan struct{}),
		running:    &sync.WaitGroup{},
	}
}

func (s *GrpcServer) Start() {
	s.running.Add(1)
	go func() {
		s.running.Done()
		if err := s.grpcServer.Serve(*s.listener); err != nil {
			shared_support.Logger().Fatal().Msgf("serve: %v", err)
		}

		s.done <- struct{}{}
	}()

	s.running.Wait()
	shared_support.Logger().Info().Msgf("gRPC server listening on %s", (*s.listener).Addr())
}

func (s *GrpcServer) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *GrpcServer) Wait() {
	<-s.done
}

func (s *GrpcServer) DoneChannel() chan struct{} {
	return s.done
}
