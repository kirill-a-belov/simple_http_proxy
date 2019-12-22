package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	c := config{}
	if err := c.parse(); err != nil {
		log.Fatal("parsing config", err)
	}

	const networkType = "tcp"
	listener, err := net.Listen(networkType, fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatal("listen starting", err)
	}
	defer listener.Close()

	for {
		newConnection, err := listener.Accept()
		if err != nil {
			if c.Logs {
				log.Print("connection accept error", err.Error())
			}

			continue
		}

		go func() {
			if err := handleConnection(newConnection); err != nil {
				if c.Logs {
					log.Print("connection handling error", err)
				}
			}
		}()
	}

}

const connectionHeaderSize = 1024

func handleConnection(inputConn net.Conn) error {
	if inputConn == nil {
		return nil
	}
	defer inputConn.Close()

	var b [connectionHeaderSize]byte
	n, err := inputConn.Read(b[:])
	if err != nil {
		return err
	}

	const minConnBytesSize = 50

	if n < minConnBytesSize {
		return errors.New(" non HTTP(S) traffic ")
	}

	var targetAddress targetAddress
	if err := targetAddress.parse(b); err != nil {
		return err
	}

	outputConn, err := net.Dial("tcp", targetAddress.address)
	if err != nil {
		return err
	}

	const proxyMethodName = "CONNECT"
	if targetAddress.method == proxyMethodName {
		if _, err := fmt.Fprint(inputConn, "HTTP/1.1 200 Connection established\r\n\r\n"); err != nil {
			return err
		}
	} else {
		if _, err := outputConn.Write(b[:n]); err != nil {
			return err
		}
	}

	go io.Copy(outputConn, inputConn)
	if _, err := io.Copy(inputConn, outputConn); err != nil {
		return err
	}

	return nil
}
