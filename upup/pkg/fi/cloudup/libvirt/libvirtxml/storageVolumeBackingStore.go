package libvirtxml

type StorageVolumeBackingStore struct {
	node *Node
}

func newStorageVolumeBackingStore(node *Node) StorageVolumeBackingStore {
	return StorageVolumeBackingStore{
		node: node,
	}
}

func (s StorageVolumeBackingStore) Path() string {
	return s.node.getAttribute(nameForLocal("path"))
}

func (s StorageVolumeBackingStore) SetPath(value string) {
	s.node.setAttribute(nameForLocal("path"), value)
}

func (s StorageVolumeBackingStore) Format() StorageVolumeTargetFormat {
	node := s.node.ensureNode(nameForLocal("format"))
	return newStorageVolumeTargetFormat(node)
}
