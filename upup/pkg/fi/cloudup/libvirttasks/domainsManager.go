package libvirttasks

import (
	"github.com/golang/glog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

// DomainsManager handles the libvirt domains owned by the cluster
//go:generate fitask -type=DomainsManager
type DomainsManager struct {
	ClusterName *string
	Domains     []*Domain
}

var _ fi.CompareWithID = &DomainsManager{}
var _ fi.HasDependencies = &DomainStart{}

func (d *DomainsManager) GetDependencies(tasks map[string]fi.Task) []fi.Task {
	var result []fi.Task

	for _, domain := range d.Domains {
		task, found := tasks["Domain/"+*domain.Name]
		if !found {
			glog.Fatalf("Unable to find Domain %s dependency for DomainsManager %s", domain.Name, *d.ClusterName)
		}

		result = append(result, task)
	}

	return result
}

func (d *DomainsManager) String() string {
	return fi.TaskAsString(d)
}

func (d *DomainsManager) CompareWithID() *string {
	return d.ClusterName
}

func (d *DomainsManager) Find(c *fi.Context) (*DomainsManager, error) {
	return nil, nil
}

func (d *DomainsManager) Run(c *fi.Context) error {
	return fi.DefaultDeltaRunMethod(d, c)
}

func (d *DomainsManager) CheckChanges(a, e, changes *DomainsManager) error {
	return nil
}

func (d *DomainsManager) RenderLibvirt(t *libvirt.LibvirtAPITarget, a, e, changes *DomainsManager) error {
	//TODO(bcrusu): _, err := t.Cloud.CreateVM(changes.Name, changes.VMTemplateName)

	// if err != nil {
	// 	return err
	// }
	return nil
}
