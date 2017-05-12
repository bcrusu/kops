package libvirt

import "k8s.io/kops/upup/pkg/fi"

// LibvirtAPITarget represents target for libvirt.
type LibvirtAPITarget struct {
	Cloud *LibvirtCloud
}

var _ fi.Target = &LibvirtAPITarget{}

func NewLibvirtAPITarget(cloud *LibvirtCloud) *LibvirtAPITarget {
	return &LibvirtAPITarget{
		Cloud: cloud,
	}
}

func (t *LibvirtAPITarget) Finish(taskMap map[string]fi.Task) error {
	return nil
}

func (t *LibvirtAPITarget) ProcessDeletions() bool {
	return true
}
