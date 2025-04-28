package program

import (
	"fmt"
	"github.com/danielealbano/svdb/engine-worker/config"
	"github.com/danielealbano/svdb/engine-worker/server"
	"github.com/danielealbano/svdb/shared/collection"
	"github.com/danielealbano/svdb/shared/grpc_server"
	"github.com/danielealbano/svdb/shared/support"
	"github.com/phuslu/log"
	"net"
	"os"
	"time"
)

type Program struct {
	config     *config.Config
	collection *shared_collection.Collection
	server     *shared_grpc_server.GrpcServer
	running    bool
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

func (p *Program) initializeCollection(shardExists bool) (*shared_collection.Collection, error) {
	// Setup the collection configuration
	collectionConfig := shared_collection.NewCollectionConfig()
	collectionConfig.MaxSize, _ = config.ParseShardMaxSize(p.config.ShardMaxSize)
	collectionConfig.Dimensions = p.config.CollectionVectorDimensions
	collectionConfig.Quantization, _ = shared_collection.ParseQuantization(p.config.CollectionQuantization)
	collectionConfig.Metric, _ = shared_collection.ParseMetric(p.config.CollectionMetric)

	// Initialize the collection
	coll, err := shared_collection.NewCollection(collectionConfig)
	if err != nil {
		return nil, err
	}

	// Load the shard if it exists
	if shardExists {
		err = coll.Load(p.config.ShardPath)
		if err != nil {
			return nil, err
		}
	}

	return coll, nil
}

func (p *Program) setupGrpcServer() (*shared_grpc_server.GrpcServer, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.config.Host, p.config.Port))
	if err != nil {
		return nil, err
	}

	grpcServer := shared_grpc_server.NewGrpcServer(&listener)
	server.RegisterCollectionGrpcServerImplementation(grpcServer, p.collection, p.config.ShardPath)

	return grpcServer, nil
}

func (p *Program) Shutdown() {
	var err error

	defer func() {
		if p.collection != nil {
			err = p.collection.Destroy()
			if err != nil {
				shared_support.Logger().Error().Msg(err.Error())
			}
		}
	}()

	if p.running {
		shared_support.Logger().Info().Msg("shutting down gRPC server (if still running)")
		p.server.Stop()
		p.server.Wait()
		shared_support.Logger().Info().Msg("gRPC server stopped")
	}

	// Save the shard if it is writeable
	if p.collection != nil && p.config.ShardWriteable {
		shared_support.Logger().Info().Msgf("saving shard to %s", p.config.ShardPath)
		err = p.collection.Save(p.config.ShardPath)
		if err != nil {
			shared_support.Logger().Error().Msg(err.Error())
			return
		}
		shared_support.Logger().Info().Msg("shard saved")
	}
}

func (p *Program) Start() {
	var err error
	var shardExists bool

	// Update the logger level
	p.updateLoggerLevel()

	// Check if the path exists
	shardExists = true
	if _, err = os.Stat(p.config.ShardPath); os.IsNotExist(err) {
		shardExists = false
		if p.config.ShardWriteable {
			shared_support.Logger().Info().Msgf("shard path %s does not exist, creating it", p.config.ShardPath)
		} else {
			shared_support.Logger().Error().Msgf("shard path %s does not exist", p.config.ShardPath)
			return
		}
	}

	// Initialize the collection
	p.collection, err = p.initializeCollection(shardExists)
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
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
