package libvirtxml

import (
	"strconv"

	"github.com/golang/glog"
)

type DomainGraphic struct {
	node *Node
}

func newDomainGraphic(node *Node) DomainGraphic {
	return DomainGraphic{
		node: node,
	}
}

func (s DomainDisk) Port() int {
	str := s.node.getAttribute(nameForLocal("port"))
	port, err := strconv.Atoi(str)
	if err != nil {
		port = 0
		glog.Warningf("ignoring invalid domain graphics port '%s'", str)
	}
	return port
}

func (s DomainDisk) SetPort(value int) {
	str := strconv.FormatInt(int64(value), 10)
	s.node.setAttribute(nameForLocal("type"), str)
}
