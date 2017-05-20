package libvirtxml

import "github.com/pkg/errors"

type Domain struct {
	doc  Document
	root *Node
}

func NewDomainForXML(xmlDoc string) (Domain, error) {
	doc := Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return Domain{}, errors.Wrap(err, "failed to unmarshal domain XML document")
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("domain"))
	}

	return Domain{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s Domain) Name() string {
	return s.root.ensureNode(nameForLocal("name")).CharData
}

func (s Domain) SetName(value string) {
	s.root.ensureNode(nameForLocal("name")).CharData = value
}

func (s Domain) UUID() string {
	return s.root.ensureNode(nameForLocal("uuid")).CharData
}

func (s Domain) SetUUID(value string) {
	s.root.ensureNode(nameForLocal("uuid")).CharData = value
}

func (s Domain) Devices() DomainDevices {
	node := s.root.ensureNode(nameForLocal("devices"))
	return newDomainDevices(node)
}
