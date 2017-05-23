package libvirttasks

import (
	"github.com/golang/glog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

// Domain represents a libvirt Domain (i.e. VM)
//go:generate fitask -type=Domain
type Domain struct {
	ClusterName          *string
	Name                 *string
	Template             *string
	BackingStorageVolume *string
}

var _ fi.CompareWithID = &Domain{}
var _ fi.HasName = &Domain{}

// GetName returns the Name of the object, implementing fi.HasName
func (vm *Domain) GetName() *string {
	return vm.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (d *Domain) SetName(name string) {
	d.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (d *Domain) String() string {
	return fi.TaskAsString(d)
}

// CompareWithID is returning name of this VirtualMachine.
func (d *Domain) CompareWithID() *string {
	glog.V(4).Info("VirtualMachine.CompareWithID invoked!")
	return d.Name
}

// Find is a no-op.
func (d *Domain) Find(c *fi.Context) (*Domain, error) {
	glog.V(4).Info("VirtualMachine.Find invoked!")
	return nil, nil
}

// Run executes DefaultDeltaRunMethod for this task.
func (d *Domain) Run(c *fi.Context) error {
	glog.V(4).Info("VirtualMachine.Run invoked!")
	return fi.DefaultDeltaRunMethod(d, c)
}

// CheckChanges is a no-op.
func (d *Domain) CheckChanges(a, e, changes *Domain) error {
	glog.V(4).Info("VirtualMachine.CheckChanges invoked!")
	return nil
}

// RenderLibvirt executes the actual VM creation.
func (d *Domain) RenderLibvirt(t *libvirt.LibvirtAPITarget, a, e, changes *Domain) error {
	glog.V(4).Infof("VirtualMachine.RenderLibvirt invoked with a(%+v) e(%+v) and changes(%+v)", a, e, changes)
	//TODO(bcrusu): _, err := t.Cloud.CreateVM(changes.Name, changes.VMTemplateName)

	// if err != nil {
	// 	return err
	// }
	return nil
}
