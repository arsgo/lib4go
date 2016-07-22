package forever

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/colinyl/daemon"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/profile"
)

type forever struct {
	dm   daemon.Daemon
	log  logger.ILogger
	svs  service
	name string
	desc string
}
type service interface {
	Start() error
	Stop() error
}

func NewForever(svs service, log logger.ILogger, name string, desc string) *forever {
	dm, err := daemon.New(name, desc)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &forever{dm: dm, name: name, desc: desc, svs: svs, log: log}
}
func (f *forever) Start() {
	defer func() {
		if r := recover(); r != nil {
			f.log.Error(r, string(debug.Stack()))
		}
	}()
	result, err := f.run()
	if err != nil {
		fmt.Println("start error:", err)
		f.log.Error("start error:", err)
		return
	}
	fmt.Println(result)
	f.log.Info(result)
}

func (f *forever) run() (string, error) {

	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status | debug | pprof_mem | pprof_block", f.name)
	if len(os.Args) > 2 {
		command := os.Args[2]
		switch command {
		case "debug":
			f.dup(f.name)
		case "pprof_mem":
			defer profile.Start(profile.MemProfile).Stop()
		case "pprof_cpu":
			defer profile.Start(profile.CPUProfile).Stop()
		case "pprof_block":
			defer profile.Start(profile.BlockProfile).Stop()
		}
	}

	if len(os.Args) > 1 {
		command := os.Args[1]
		logger.SetDebug(false)
		switch command {
		case "install":
			return f.dm.Install()
		case "remove":
			return f.dm.Remove()
		case "start":
			return f.dm.Start()
		case "stop":
			return f.dm.Stop()
		case "status":
			return f.dm.Status()
		case "debug":
			f.dup(f.name)
		case "pprof_mem":
			defer profile.Start(profile.MemProfile).Stop()
		case "pprof_cpu":
			defer profile.Start(profile.CPUProfile).Stop()
		case "pprof_block":
			defer profile.Start(profile.BlockProfile).Stop()
		default:
			return usage, nil
		}
	}
	if err := f.svs.Start(); err != nil {
		return "", err
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			f.svs.Stop()
			f.dm.Start()
			return fmt.Sprintf("%s was killed", f.name), nil
		}
	}
	// never happen, but need to complete code
	return usage, nil
}
