package libvirt

import (
	libvirt "github.com/libvirt/libvirt-go" //TODO(bcrusu) add libvirt to vendor directory
)

//TODO(bcrusu): Add retry functionality
type libvirtConnection struct {
	connect *libvirt.Connect
}

func newLibvirtConnection(uri string) (*libvirtConnection, error) {
	connect, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, err
	}

	return &libvirtConnection{connect}, nil
}

func (c *libvirtConnection) Close() {
	c.Close()
}
