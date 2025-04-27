package main

import (
	"fmt"
	"github.com/danielealbano/cvdb/engine-worker/collection"
	"github.com/danielealbano/cvdb/engine-worker/config"
	"github.com/danielealbano/cvdb/engine-worker/server"
	shared_support "github.com/danielealbano/cvdb/shared/support"
	log "github.com/phuslu/log"
	"net"
	"os"
	"time"
)

var (
	version       = ""
	commit        = ""
	buildDate     = ""
	builtBy       = ""
	goLangVersion = ""
)

func main() {
	shared_support.GenericMain(
		mainReal,
		version,
		commit,
		buildDate,
		builtBy,
		goLangVersion)
}

func mainReal() {
	var err error
	var cfg *config.Config
	var coll *collection.Collection
	var shardExists bool

	cfg, err = config.FromEnv()
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
	}

	// Update the logger level
	shared_support.Logger().Level = log.ParseLevel(cfg.LogLevel)

	// Initialize the collection
	collectionConfig := collection.NewCollectionConfig()
	collectionConfig.MaxSize, _ = config.ParseShardMaxSize(cfg.ShardMaxSize)
	collectionConfig.Dimensions = cfg.CollectionVectorDimensions
	collectionConfig.Quantization, _ = collection.ParseQuantization(cfg.CollectionQuantization)
	collectionConfig.Metric, _ = collection.ParseMetric(cfg.CollectionMetric)
	coll, err = collection.NewCollection(collectionConfig)
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
	}
	defer func(coll *collection.Collection) {
		err = coll.Destroy()
		if err != nil {
			shared_support.Logger().Error().Msg(err.Error())
		}
	}(coll)

	// Check if the pah exists
	shardExists = true
	if _, err = os.Stat(cfg.ShardPath); os.IsNotExist(err) {
		shardExists = false
		if cfg.ShardWriteable {
			shared_support.Logger().Info().Msgf("shard path %s does not exist, creating it", cfg.ShardPath)
		} else {
			shared_support.Logger().Error().Msgf("shard path %s does not exist", cfg.ShardPath)
			return
		}
	}

	// Load the shard if it exists
	if shardExists {
		err = coll.Load(cfg.ShardPath)
		if err != nil {
			shared_support.Logger().Error().Msg(err.Error())
			return
		}
	}

	// Initialize the gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
	}

	// Setup the gRPC server and register the collection grpc server implementation
	grpcServer := server.NewGrpcServer(&listener, coll)
	server.RegisterCollectionGrpcServerImplementation(grpcServer, cfg.ShardPath)

	// Start the gRPC server
	grpcServer.Start()

	// Wait for the service to be stopped or the grpc server to exit
	shared_support.WaitMultipleChannels(
		500*time.Millisecond,
		shared_support.StopSignal.Context.Done(),
		grpcServer.DoneChannel())

	shared_support.Logger().Info().Msg("gRPC server stopped or stop signal received, shutting down")

	// Stop the gRPC server
	grpcServer.Stop()

	// Wait for the gRPC server to finish
	grpcServer.Wait()

	// Save the shard if it is writeable
	if cfg.ShardWriteable {
		shared_support.Logger().Info().Msgf("saving shard to %s", cfg.ShardPath)
		err = coll.Save(cfg.ShardPath)
		if err != nil {
			shared_support.Logger().Error().Msg(err.Error())
			return
		}
		shared_support.Logger().Info().Msg("shard saved")
	}

	shared_support.Logger().Info().Msg("exiting...")
}
