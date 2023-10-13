// Package udp is the udp display.
package udp

import (
	"context"
	"net"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the udp transporter.
type Transporter struct {
	//address *net.UDPAddr
	// udpServer net.PacketConn
	udpServer net.Conn
	formater  serialize.Serializer
}

func New(ctx context.Context, formater serialize.Serializer, addr string) (*Transporter, error) {
	// splitter := strings.Split(addr, ":")

	// p := make([]byte, 2048)
	udpServer, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	// _, err = bufio.NewReader(conn).Read(p)
	// if err == nil {
	// 	fmt.Printf("%s\n", p)
	// } else {
	// 	fmt.Printf("Some error %v\n", err)
	// }
	// conn.Close()

	// udpServer, err := net.ListenPacket("udp", fmt.Sprintf(":%s", splitter[len(splitter)-1]))
	// if err != nil {
	// 	return nil, err
	// }

	go func() {
		<-ctx.Done()
		_ = udpServer.Close()
	}()

	// address, err := net.ResolveUDPAddr("udp4", addr)
	// if err != nil {
	// 	return nil, err
	// }

	return &Transporter{
		//address:   address,
		formater:  formater,
		udpServer: udpServer,
	}, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *dump.Aircraft) error {
	data, err := t.formater.Serialize(ac)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	data = append(data, '\n')

	_, err = t.udpServer.Write(data)

	return err
}
