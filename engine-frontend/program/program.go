package program

import (
	"fmt"
	"github.com/danielealbano/svdb/engine-frontend/config"
	"github.com/danielealbano/svdb/engine-frontend/server"
	"github.com/danielealbano/svdb/shared/collection"
	"github.com/danielealbano/svdb/shared/grpc_server"
	"github.com/danielealbano/svdb/shared/support"
	"github.com/phuslu/log"
	"net"
	"os"
	"time"
)

type Program struct {
	config           *config.Config
	collectionConfig *shared_collection.CollectionConfig
	server           *shared_grpc_server.GrpcServer
	running          bool
}

func NewProgram(config *config.Config) *Program {
	return &Program{
		config: config,
	}
}

func (p *Program) GetConfig() *config.Config {
	return p.config
}

func (p *Program) updateLoggerLevel() {
	shared_support.Logger().Level = log.ParseLevel(p.config.LogLevel)
}

func (p *Program) setupGrpcServer() (*shared_grpc_server.GrpcServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.config.Host, p.config.Port))
	if err != nil {
		return nil, err
	}

	grpcServer := shared_grpc_server.NewGrpcServer(&listener)
	server.RegisterFrontendGrpcServerImplementation(grpcServer, p.collectionConfig)

	return grpcServer, nil
}

func (p *Program) Shutdown() {
	if p.running {
		shared_support.Logger().Info().Msg("shutting down gRPC server (if still running)")
		p.server.Stop()
		p.server.Wait()
		shared_support.Logger().Info().Msg("gRPC server stopped")
	}

	// TODO: Ensure all the shards are save and closed
}

func (p *Program) Start() {
	var err error

	// Update the logger level
	p.updateLoggerLevel()

	// Check if the path exists
	if _, err = os.Stat(p.config.CollectionPath); os.IsNotExist(err) {
		// TODO: Initialize the collection
		shared_support.Logger().Info().Msgf("the collection does not exist, initializing it")
	} else {
		// TODO: Load the collection
		//p.collection, err = p.initializeCollection(collectionExists)
		//if err != nil {
		//	shared_support.Logger().Error().Msg(err.Error())
		//	return
		//}
	}

	// Start the gRPC server
	p.server, err = p.setupGrpcServer()
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
	}
	p.server.Start()
}

func (p *Program) Wait() {
	shared_support.WaitMultipleChannels(
		500*time.Millisecond,
		shared_support.StopSignal.Context.Done(),
		p.server.DoneChannel())
}
