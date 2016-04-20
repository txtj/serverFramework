package core

import (
	"net"
	"os"
	"sync"
	"time"

	"serverFramework/internal/moduledb"
	"serverFramework/internal/modulelogic"
	"serverFramework/internal/util"
)

var (
	// BeeApp is an application instance
	Server *ServerCore
)

type ServerCore struct {
	sync.RWMutex
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64
	startTime        time.Time
	tcpListener      net.Listener
	wg               util.WaitGroupWrapper

	db    chan *moduledb.DBModuler
	logic chan *modulelogic.LogicModuler
}

func init() {
	Server = New()
}

func New() *ServerCore {
	sc := &ServerCore{
		startTime: time.Now(),
	}

	return sc
}

func (sc *ServerCore) Run() {
	tcpListener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		BeeLogger.Error("listen [%s] failed ->%s", "localhost", err)
		os.Exit(0)
	}

	sc.Lock()
	sc.tcpListener = tcpListener
	sc.Unlock()

	Info("server listen on 8888")

	ctx := &context{sc}

	handle := &ServerHandler{ctx: ctx}
	// start accept routine
	sc.wg.Wrap(func() {
		HandleAccept(sc.tcpListener, handle)
	})

	sc.wg.Wait()
}