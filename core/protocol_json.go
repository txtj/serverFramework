package core

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/TaXingTianJi/serverFramework/protocol"
)

type ProtocolJson struct {
}

func init() {
	protocol.Register("  json", &ProtocolV1{})
}

func (p *ProtocolJson) IOLoop(conn net.Conn) error {
	fmt.Printf("[ProtocolJson::Loop] loop ...\n")
	var err error
	var line []byte

	clientId := atomic.AddInt64(&ServerApp.clientIDSequence, 1)
	client := NewClient(clientId, conn)

	fmt.Printf("[ProtocolJson::Loop] clientid[%d]\n", client.ID)

	for {
		line, err = client.Reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read [%s]\n", err)
			}
			break
		}

		line = line[:len(line)-1]
		fmt.Printf("[ProtocolJson::Loop] line[%v]\n", line)
	}

	conn.Close()

	fmt.Printf("[ProtocolJson::Loop] loop exit\n")

	return err
}
