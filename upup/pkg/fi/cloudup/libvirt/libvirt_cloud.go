package libvirt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kubernetes/federation/pkg/dnsprovider"
	k8scoredns "k8s.io/kubernetes/federation/pkg/dnsprovider/providers/coredns"
)

// LibvirtCloud represents a li cloud instance.
type LibvirtCloud struct {
	URI           string
	Connection    *libvirtConnection
	CoreDNSServer string
	DNSZone       string
}

var _ fi.Cloud = &LibvirtCloud{}

// ProviderID returns ID for libvirt type cloud provider.
func (c *LibvirtCloud) ProviderID() fi.CloudProviderID {
	return fi.CloudProviderLibvirt
}

// NewLibvirtCloud returns LibvirtCloud instance for given ClusterSpec.
func NewLibvirtCloud(spec *kops.ClusterSpec) (*LibvirtCloud, error) {
	uri := *spec.CloudConfig.LibvirtURI
	glog.V(2).Infof("Creating libvirt Cloud with uri(%s)", uri)

	c, err := newLibvirtConnection(uri)
	if err != nil {
		return nil, err
	}
	//TODO(bcrusu): close the connection

	result := &LibvirtCloud{
		URI:           uri,
		Connection:    c,
		CoreDNSServer: *spec.CloudConfig.LibvirtCoreDNSServer,
		DNSZone:       spec.DNSZone,
	}

	glog.V(2).Infof("Created libvirt Cloud successfully: %+v", result)
	return result, nil
}

// DNS returns dnsprovider interface for this cloud.
//TODO(bcrusu): review
func (c *LibvirtCloud) DNS() (dnsprovider.Interface, error) {
	var provider dnsprovider.Interface
	var err error
	var lines []string
	lines = append(lines, "etcd-endpoints = "+c.CoreDNSServer)
	lines = append(lines, "zones = "+c.DNSZone)
	config := "[global]\n" + strings.Join(lines, "\n") + "\n"
	file := bytes.NewReader([]byte(config))
	provider, err = dnsprovider.GetDnsProvider(k8scoredns.ProviderName, file)
	if err != nil {
		return nil, fmt.Errorf("Error building (k8s) DNS provider: %v", err)
	}

	return provider, nil
}

// FindVPCInfo is not supported.
func (c *LibvirtCloud) FindVPCInfo(id string) (*fi.VPCInfo, error) {
	glog.Warning("FindVPCInfo not supported by libvirt Cloud")
	return nil, nil
}
