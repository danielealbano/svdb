package shared_support

import (
	"golang.org/x/net/context"
	"os/signal"
	"syscall"
)

type SignalContextAndStop struct {
	Context context.Context
	Cancel  context.CancelFunc
}

var StopSignal *SignalContextAndStop
var Usr1Signal *SignalContextAndStop
var Usr2Signal *SignalContextAndStop

func catchTerminationSignals() *SignalContextAndStop {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	return &SignalContextAndStop{
		Context: ctx,
		Cancel:  cancel,
	}
}

func catchUsr1Signal() *SignalContextAndStop {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGUSR1)

	return &SignalContextAndStop{
		Context: ctx,
		Cancel:  cancel,
	}
}

func catchUsr2Signal() *SignalContextAndStop {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGUSR2)

	return &SignalContextAndStop{
		Context: ctx,
		Cancel:  cancel,
	}
}

func ResetSignals() {
	signal.Reset(syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
}

func SetupSignalsCatching() {
	StopSignal = catchTerminationSignals()
	Usr1Signal = catchUsr1Signal()
	Usr2Signal = catchUsr2Signal()
}
