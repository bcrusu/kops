/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package libvirttasks

import (
	"github.com/golang/glog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

// DomainStart starts a suspended/defined domain
//go:generate fitask -type=DomainStart
type DomainStart struct {
	Name   *string
	Domain *Domain
}

var _ fi.HasName = &DomainStart{}
var _ fi.HasDependencies = &DomainStart{}

func (d *DomainStart) GetDependencies(tasks map[string]fi.Task) []fi.Task {
	domain := tasks["Domain/"+*d.Domain.Name]
	if domain == nil {
		glog.Fatalf("Unable to find Domain %s dependency for DomainStart %s", *d.Domain.Name, *d.Name)
	}

	return []fi.Task{domain}
}

func (d *DomainStart) GetName() *string {
	return d.Name
}

func (d *DomainStart) SetName(name string) {
	d.Name = &name
}

func (d *DomainStart) Run(c *fi.Context) error {
	glog.Info("VMPowerOn.Run invoked!")
	return fi.DefaultDeltaRunMethod(d, c)
}

func (d *DomainStart) Find(c *fi.Context) (*DomainStart, error) {
	glog.Info("VMPowerOn.Find invoked!")
	return nil, nil
}

func (d *DomainStart) CheckChanges(a, e, changes *DomainStart) error {
	glog.Info("VMPowerOn.CheckChanges invoked!")
	return nil
}

func (d *DomainStart) RenderLibvirt(t *libvirt.LibvirtAPITarget, a, e, changes *DomainStart) error {
	//TODO(bcrusu)
	return nil
}
