/*
Copyright 2016 The Kubernetes Authors.

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
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

//TODO(bcrusu):
type StorageVolume struct {
	Pool       *string
	Name       *string
	VolumeType *string
	SizeGB     *int64
	Encrypted  *bool
}

var _ fi.CompareWithID = &StorageVolume{}

func (e *StorageVolume) CompareWithID() *string {
	return e.Name
}

func (e *StorageVolume) Find(c *fi.Context) (*StorageVolume, error) {
	return nil, nil
}

func (e *StorageVolume) Run(c *fi.Context) error {
	return fi.DefaultDeltaRunMethod(e, c)
}

func (_ *StorageVolume) CheckChanges(a, e, changes *StorageVolume) error {
	return nil
}

func (_ *StorageVolume) RenderLibvirt(t *libvirt.LibvirtAPITarget, a, e, changes *StorageVolume) error {
	return nil
}
