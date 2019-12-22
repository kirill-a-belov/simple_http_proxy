package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

type targetAddress struct {
	address string
	method  string
}

func (a *targetAddress) parse(input [connectionHeaderSize]byte) error {
	var host string
	if _, err := fmt.Sscanf(string(input[:bytes.IndexByte(input[:], '\n')]), "%s%s", &a.method, &host); err != nil {
		return err
	}

	hostPortURL, err := url.Parse(host)
	if err != nil {
		return err
	}

	if hostPortURL.Opaque == "443" {
		a.address = hostPortURL.Scheme + ":443"
	} else {
		if strings.Index(hostPortURL.Host, ":") == -1 {
			a.address = hostPortURL.Host + ":80"
		} else {
			a.address = hostPortURL.Host
		}
	}

	return nil
}

