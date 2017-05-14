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

package resources

import (
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt"
)

type clusterDiscoveryLibvirt struct {
	cloud       *libvirt.LibvirtCloud
	clusterName string
}

type libvirtListFn func() ([]*ResourceTracker, error)

func (c *AwsCluster) listResourcesLibvirt() (map[string]*ResourceTracker, error) {
	resources := make(map[string]*ResourceTracker)

	d := &clusterDiscoveryLibvirt{
		cloud:       c.Cloud.(*libvirt.LibvirtCloud),
		clusterName: c.ClusterName,
	}

	listFunctions := []libvirtListFn{
		d.listVMs,
	}

	for _, fn := range listFunctions {
		trackers, err := fn()
		if err != nil {
			return nil, err
		}
		for _, t := range trackers {
			resources[GetResourceTrackerKey(t)] = t
		}
	}

	return resources, nil
}

func (d *clusterDiscoveryLibvirt) listVMs() ([]*ResourceTracker, error) {
	var trackers []*ResourceTracker

	//TODO(bcrusu)

	return trackers, nil
}
