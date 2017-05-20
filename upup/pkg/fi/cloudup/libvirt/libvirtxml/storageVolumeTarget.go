package libvirtxml

type StorageVolumeTarget struct {
	node *Node
}

func newStorageVolumeTarget(node *Node) StorageVolumeTarget {
	return StorageVolumeTarget{
		node: node,
	}
}

func (s StorageVolumeTarget) Path() string {
	return s.node.getAttribute(nameForLocal("path"))
}

func (s StorageVolumeTarget) SetPath(value string) {
	s.node.setAttribute(nameForLocal("path"), value)
}

func (s StorageVolumeTarget) Format() StorageVolumeTargetFormat {
	node := s.node.ensureNode(nameForLocal("format"))
	return newStorageVolumeTargetFormat(node)
}
