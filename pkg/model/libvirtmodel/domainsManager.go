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

package libvirtmodel

import (
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/model"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirttasks"
)

type DomainsManagerModelBuilder struct {
	*LibvirtModelContext

	BootstrapScript *model.BootstrapScript
}

var _ fi.ModelBuilder = &DomainsManagerModelBuilder{}

func (b *DomainsManagerModelBuilder) Build(c *fi.ModelBuilderContext) error {
	clusterName := b.ClusterName()
	var allDomains []*libvirttasks.Domain

	for igIndex, ig := range b.InstanceGroups {
		domainTemplate := ig.Spec.MachineType
		backingStorageVolume := ig.Spec.Image

		instanceCount := int(fi.Int32Value(ig.Spec.MaxSize))
		for i := 1; i <= instanceCount; i++ {
			instanceName, err := b.instanceName(ig, igIndex+1, i)
			if err != nil {
				return err
			}

			domainTask := &libvirttasks.Domain{
				ClusterName:          &clusterName,
				Name:                 &instanceName,
				Template:             &domainTemplate,
				BackingStorageVolume: &backingStorageVolume,
			}

			c.AddTask(domainTask)
			allDomains = append(allDomains, domainTask)

			domainStartTask := &libvirttasks.DomainStart{
				Name:   &instanceName,
				Domain: domainTask,
			}

			c.AddTask(domainStartTask)
		}
	}

	domainsManager := &libvirttasks.DomainsManager{
		ClusterName: &clusterName,
		Domains:     allDomains,
	}
	c.AddTask(domainsManager)

	return nil
}

func (b *DomainsManagerModelBuilder) instanceName(ig *kops.InstanceGroup, igNumber int, domainNumber int) (string, error) {
	switch ig.Spec.Role {
	case kops.InstanceGroupRoleMaster, kops.InstanceGroupRoleNode, kops.InstanceGroupRoleBastion:
		name := b.ClusterName() + "." + ig.ObjectMeta.Name + "." + strconv.Itoa(igNumber) + "." + strconv.Itoa(domainNumber)
		return name, nil
	default:
		return "", errors.Errorf("unknown InstanceGroup Role: %v", ig.Spec.Role)
	}
}
