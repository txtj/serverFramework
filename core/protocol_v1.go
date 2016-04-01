package core

import (
	// "bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"encoding/binary"
)

/* 
//the protocol v1 be made of [header cmd len data]
// header: one byte
// cmd: two byte,
// len: two byte,length of the data
// data: the contents

//   	[x]	   [x][x]      [x][x]  [x][x][x][x]...
//   | byte | (uint32) |  (uint32)  |   (binary)
//   |1-byte|  2-byte  |   2-byte   |	 N-byte
//   --------------------------------------------...
//    header    cmd	      length         data
//
//   [0x05 cmd len data]
*/
type ProtocolV1 struct {
	ctx *context
}

func (p *ProtocolV1) IOLoop(conn net.Conn) error {
	fmt.Printf("[ProtocolV1::Loop] loop ...\n")
	var err error	
	var header byte
	var cmd uint32
	var length int32

	clientId := atomic.AddInt64(&p.ctx.core.clientIDSequence, 1)
	client := newClient(clientId, conn, p.ctx)

	tt := make([]byte,1024)
	for {
		_,err = io.ReadFull(conn,tt)
		if err != nil {
			return err
		}
		fmt.Printf("ProtocolV1 tt[%v]\n",tt)

		// header
		err = binary.Read(conn, binary.BigEndian,&header)
		if err != nil {
			return err
		}
		// if header != 0x05 {
		// 	return fmt.Errorf("ProtocolV1 header[%s] err\n",header)
		// }
		fmt.Printf("ProtocolV1 header[%v]\n",header)

		// cmd
		err = binary.Read(conn,binary.BigEndian,cmd)
		if err != nil {
			return err
		}
		fmt.Printf("ProtocolV1 cmd[%v]\n",cmd)

		// data
		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return err
		}
		data := make([]byte,length)
		_,err = io.ReadFull(conn,data)
		if err != nil {
			return err
		}
		fmt.Printf("ProtocolV1 len[%d] data[%v]\n",length,data)


		// // message binary data
		// buf := make([]byte, len)
		// _, err = io.ReadFull(r, buf)
		// if err != nil {
		// 	return nil, err
		// }

		// var line []byte
		// line, err = client.Reader.ReadSlice('\n')
		// if err != nil {
		// 	if err == io.EOF {
		// 		err = nil
		// 	} else {
		// 		err = fmt.Errorf("failed to read - %s\n", err)
		// 	}
		// 	break
		// } else if err == nil{
		// 	fmt.Printf("line-%v\n",line)
		// }

		// // trim the '\n'
		// line = line[:len(line)-1]
		// // optionally trim the '\r'
		// if len(line) > 0 && line[len(line)-1] == '\r' {
		// 	line = line[:len(line)-1]
		// }
		// params := bytes.Split(line, []byte(" "))
		// fmt.Printf("[ProtocolV1::Loop] line[%v] params[%v]\n", line, params)

		err = p.Send(client, []byte("string send to client"))
		if err != nil {
			err = fmt.Errorf("failed to send response - %s", err)
			break
		}

	}

	conn.Close()

	fmt.Printf("[ProtocolV1::Loop] loop exit err - %v\n",err)

	return err
}

func (p *ProtocolV1) Send(client *ClientV1, data []byte) error {
	client.writeLock.Lock()

	// var zeroTime time.Time
	// if client.HeartbeatInterval > 0 {
	// 	client.SetWriteDeadline(time.Now().Add(client.HeartbeatInterval))
	// } else {
	// 	client.SetWriteDeadline(zeroTime)
	// }

	// _, err := SendFramedResponse(client.Writer, frameType, data)
	_, err := SendResponse(client.Writer, data)
	if err != nil {
		client.writeLock.Unlock()
		return err
	}

	// if frameType != frameTypeMessage {
	// 	err = client.Flush()
	// }

	client.Writer.Flush()

	client.writeLock.Unlock()

	return err
}