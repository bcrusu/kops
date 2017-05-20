package libvirttasks

import (
	"github.com/golang/glog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

// VirtualMachine represents a libvirt VM
//go:generate fitask -type=VirtualMachine
type VirtualMachine struct {
	Name           *string
	VMTemplateName *string
}

var _ fi.CompareWithID = &VirtualMachine{}
var _ fi.HasName = &VirtualMachine{}

// GetName returns the Name of the object, implementing fi.HasName
func (vm *VirtualMachine) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *VirtualMachine) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *VirtualMachine) String() string {
	return fi.TaskAsString(o)
}

// CompareWithID is returning name of this VirtualMachine.
func (e *VirtualMachine) CompareWithID() *string {
	glog.V(4).Info("VirtualMachine.CompareWithID invoked!")
	return e.Name
}

// Find is a no-op.
func (e *VirtualMachine) Find(c *fi.Context) (*VirtualMachine, error) {
	glog.V(4).Info("VirtualMachine.Find invoked!")
	return nil, nil
}

// Run executes DefaultDeltaRunMethod for this task.
func (e *VirtualMachine) Run(c *fi.Context) error {
	glog.V(4).Info("VirtualMachine.Run invoked!")
	return fi.DefaultDeltaRunMethod(e, c)
}

// CheckChanges is a no-op.
func (_ *VirtualMachine) CheckChanges(a, e, changes *VirtualMachine) error {
	glog.V(4).Info("VirtualMachine.CheckChanges invoked!")
	return nil
}

// RenderLibvirt executes the actual VM creation.
func (_ *VirtualMachine) RenderLibvirt(t *libvirt.LibvirtAPITarget, a, e, changes *VirtualMachine) error {
	glog.V(4).Infof("VirtualMachine.RenderLibvirt invoked with a(%+v) e(%+v) and changes(%+v)", a, e, changes)
	//TODO(bcrusu): _, err := t.Cloud.CreateVM(changes.Name, changes.VMTemplateName)

	// if err != nil {
	// 	return err
	// }
	return nil
}
