package dropsonde

import (
	"errors"
	"github.com/cloudfoundry-incubator/dropsonde/emitter"
	"github.com/cloudfoundry-incubator/dropsonde/heartbeat"
	"sync"
)

var DefaultEmitterRemoteAddr = "localhost:42420"
var HeartbeatEmitterRemoteAddr = "localhost:42421"

var heartbeatState struct {
	sync.Mutex
	stopChannel chan<- interface{}
}

func Initialize(origin string) error {
	if len(origin) == 0 {
		return errors.New("Cannot initialize dropsonde without an origin")
	}

	if emitter.DefaultEmitter == nil {
		udpEmitter, err := emitter.NewUdpEmitter(DefaultEmitterRemoteAddr, origin)
		if err != nil {
			return err
		}

		emitter.DefaultEmitter, err = emitter.NewInstrumentedEmitter(udpEmitter)
		if err != nil {
			return err
		}
	}

	heartbeatState.Lock()
	defer heartbeatState.Unlock()

	if heartbeatState.stopChannel != nil {
		return nil
	}

	if heartbeatEventSource, ok := emitter.DefaultEmitter.(heartbeat.HeartbeatEventSource); ok {
		var err error
		if heartbeat.HeartbeatEmitter == nil {
			heartbeat.HeartbeatEmitter, err = emitter.NewTcpEmitter(HeartbeatEmitterRemoteAddr, origin)
			if err != nil {
				return err
			}
		}

		heartbeatState.stopChannel, err = heartbeat.BeginGeneration(heartbeatEventSource)
		if err != nil {
			return err
		}
	}

	return nil
}

func Cleanup() {
	heartbeatState.Lock()
	defer heartbeatState.Unlock()

	if heartbeatState.stopChannel != nil {
		close(heartbeatState.stopChannel)
		heartbeatState.stopChannel = nil
	}
}
