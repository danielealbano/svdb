package main

import (
	"github.com/danielealbano/svdb/engine-worker/config"
	"github.com/danielealbano/svdb/engine-worker/program"
	shared_support "github.com/danielealbano/svdb/shared/support"
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

	cfg, err = config.FromEnv()
	if err != nil {
		shared_support.Logger().Error().Msg(err.Error())
		return
	}

	p := program.NewProgram(cfg)
	defer p.Shutdown()
	p.Start()
	p.Wait()

	shared_support.Logger().Info().Msg("gRPC server stopped or stop signal received, shutting down")
}
