package libvirt

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt/libvirtxml"
)

const uuidStringLength = 36

var random *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func getCapabilitiesXML(connect *libvirt.Connect) (libvirtxml.Capabilities, error) {
	xml, err := connect.GetCapabilities()
	if err != nil {
		return libvirtxml.Capabilities{}, errors.Wrapf(err, "failed to fetch libvirt capabilities")
	}

	return libvirtxml.NewCapabilitiesForXML(xml)
}

func lookupStoragePool(connect *libvirt.Connect, lookup string) (*libvirt.StoragePool, error) {
	if len(lookup) == uuidStringLength {
		if pool, _ := connect.LookupStoragePoolByUUIDString(lookup); pool != nil {
			return pool, nil
		}
	}

	if pool, _ := connect.LookupStoragePoolByName(lookup); pool != nil {
		return pool, nil
	}

	return nil, errors.Errorf("could not find storage pool '%s'", lookup)
}

func lookupDomain(connect *libvirt.Connect, lookup string) (*libvirt.Domain, error) {
	if len(lookup) == uuidStringLength {
		if domain, _ := connect.LookupDomainByUUIDString(lookup); domain != nil {
			return domain, nil
		}
	}

	if domain, _ := connect.LookupDomainByName(lookup); domain != nil {
		return domain, nil
	}

	return nil, errors.Errorf("could not find domain '%s'", lookup)
}

func getDomainXML(domain *libvirt.Domain) (libvirtxml.Domain, error) {
	xml, err := domain.GetXMLDesc(libvirt.DomainXMLFlags(0))
	if err != nil {
		return libvirtxml.Domain{}, errors.Wrapf(err, "failed to fetch domain XML description")
	}

	return libvirtxml.NewDomainForXML(xml)
}

func listAllDomains(connect *libvirt.Connect) ([]libvirtxml.Domain, error) {
	var result []libvirtxml.Domain

	flags := libvirt.CONNECT_LIST_DOMAINS_ACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_INACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_PERSISTENT |
		libvirt.CONNECT_LIST_DOMAINS_TRANSIENT |
		libvirt.CONNECT_LIST_DOMAINS_RUNNING |
		libvirt.CONNECT_LIST_DOMAINS_PAUSED |
		libvirt.CONNECT_LIST_DOMAINS_SHUTOFF |
		libvirt.CONNECT_LIST_DOMAINS_OTHER |
		libvirt.CONNECT_LIST_DOMAINS_MANAGEDSAVE |
		libvirt.CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE |
		libvirt.CONNECT_LIST_DOMAINS_AUTOSTART |
		libvirt.CONNECT_LIST_DOMAINS_NO_AUTOSTART |
		libvirt.CONNECT_LIST_DOMAINS_HAS_SNAPSHOT |
		libvirt.CONNECT_LIST_DOMAINS_NO_SNAPSHOT

	domains, err := connect.ListAllDomains(flags)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list domains")
	}

	for _, domain := range domains {
		domainXML, err := getDomainXML(&domain)
		if err != nil {
			return nil, err
		}

		result = append(result, domainXML)
		domain.Free()
	}

	return result, nil
}

func cloneStorageVolume(pool *libvirt.StoragePool, sourceName string, destName string) error {
	volumeXML, err := lookupStorageVolume(pool, sourceName)
	if err != nil {
		return err
	}

	volumeType := volumeXML.Type()
	if volumeType != "file" {
		errors.Errorf("cannot clone storage volume '%s' - unsupported volume type '%s'", sourceName, volumeType)
	}

	volumeXML.SetName(destName)
	volumeXML.SetKey("")

	targetXML := volumeXML.Target()
	targetXML.RemoveTimestamps()

	sourcePath := targetXML.Path()
	targetXML.SetPath("") // will be filled-in by libvirt

	{
		// set backing store as the souorce target
		backingStoreXML := volumeXML.BackingStore()
		backingStoreXML.SetPath(sourcePath)
		backingStoreXML.Format().SetType(targetXML.Format().Type())
		backingStoreXML.RemoveTimestamps()
	}

	// switch to a format that supports backing store
	switch targetXML.Format().Type() {
	case "raw":
		targetXML.Format().SetType("qcow2")
	}

	xmlString, err := volumeXML.MarshalToXML()
	if err != nil {
		return err
	}

	storageVol, err := pool.StorageVolCreateXML(xmlString, libvirt.StorageVolCreateFlags(0))
	if err != nil {
		return errors.Wrapf(err, "failed to clone storage volume '%s' to '%s'", sourceName, destName)
	}
	defer storageVol.Free()

	return err
}

func lookupStorageVolume(pool *libvirt.StoragePool, volumeName string) (libvirtxml.StorageVolume, error) {
	volume, err := pool.LookupStorageVolByName(volumeName)
	if err != nil {
		return libvirtxml.StorageVolume{}, errors.Errorf("could not find storage volume '%s'", volumeName)
	}
	defer volume.Free()

	xml, err := volume.GetXMLDesc(0)
	if err != nil {
		return libvirtxml.StorageVolume{}, errors.Wrapf(err, "failed to fetch XML description for storage volume '%s'", volumeName)
	}

	return libvirtxml.NewStorageVolumeForXML(xml)
}

func randomMACAddressNoConflict(connect *libvirt.Connect) (string, error) {
	uri, err := connect.GetURI()
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch libvirt connection uri")
	}

	allDomains, err := listAllDomains(connect)
	if err != nil {
		return "", err
	}

	for i := 0; i < 256; i++ {
		mac, err := randomMACAddress(uri)
		if err != nil {
			return "", err
		}

		if hasConflictingMACAddress(allDomains, mac) {
			continue
		}

		return mac, nil
	}

	return "", errors.New("failed to generate non-conflicting MAC address")
}

func randomMACAddress(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse libvirt connection uri")
	}

	var mac []byte

	if isQemuURL(url) {
		mac = []byte{0x52, 0x54, 0x00}
	} else if isXenURL(url) {
		mac = []byte{0x00, 0x16, 0x3E}
	}

	for len(mac) < 6 {
		b := random.Uint32()
		mac = append(mac, byte(b))
	}

	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]), nil
}

func isQemuURL(url *url.URL) bool {
	return strings.HasPrefix(url.Scheme, "qemu")
}

func isXenURL(url *url.URL) bool {
	return strings.HasPrefix(url.Scheme, "xen") ||
		strings.HasPrefix(url.Scheme, "libxl")
}

func hasConflictingMACAddress(domains []libvirtxml.Domain, mac string) bool {
	mac = strings.ToLower(mac)

	for _, domain := range domains {
		interfaces := domain.Devices().Interfaces()
		for _, iface := range interfaces {
			ifaceMAC := iface.MACAddress()
			if strings.ToLower(ifaceMAC) == mac {
				return true
			}
		}
	}

	return false
}

// inspiration drawn from https://github.com/virt-manager/virt-manager/blob/master/virtinst/cloner.py
func createDomain(connect *libvirt.Connect, name string, domainTemplateXML string, diskPath string) error {
	domainXML, err := libvirtxml.NewDomainForXML(domainTemplateXML)
	if err != nil {
		return err
	}

	domainXML.SetID("")
	domainXML.SetUUID("")
	domainXML.SetName(name)

	// Set the graphics device port to auto, in order to avoid conflicts
	graphics := domainXML.Devices().Graphics()
	for _, graphic := range graphics {
		graphic.SetPort(-1)
	}

	// generate random MAC address for network interfaces
	interfaces := domainXML.Devices().Interfaces()
	for _, iface := range interfaces {
		mac, err := randomMACAddressNoConflict(connect)
		if err != nil {
			return err
		}

		iface.SetTargetDevice("")
		iface.SetMACAddress(mac)
	}

	// reset path for guest agent channel
	channels := domainXML.Devices().Channels()
	for _, channel := range channels {
		if channel.Type() != "unix" {
			continue
		}

		// will be set by libvirt
		channel.SetSourcePath("")
	}

	if domainXML.Devices().Emulator() == "" {
		setEmulator(connect, domainXML)
	}

	setDiskPath(domainXML, diskPath)

	xml, err := domainXML.MarshalToXML()
	if err != nil {
		return err
	}

	domain, err := connect.DomainCreateXML(xml, libvirt.DomainCreateFlags(0))
	if err != nil {
		return errors.Wrapf(err, "failed to create domain '%s'", name)
	}
	defer domain.Free()

	return nil
}

func setEmulator(connect *libvirt.Connect, domain libvirtxml.Domain) error {
	capabilities, err := getCapabilitiesXML(connect)
	if err != nil {
		return err
	}

	hostArch := capabilities.Host().CPU().Arch()
	guests := capabilities.Guests()

	var emulator string
	for _, guest := range guests {
		if guest.Arch().Name() != hostArch {
			continue
		}

		emulator = guest.Arch().Emulator()

		if guest.OSType() == "hvm" {
			// found hardware-assisted vm - use this emulator
			break
		}
	}

	if emulator == "" {
		return errors.Errorf("found no guest matching host architecture '%s'", hostArch)
	}

	domain.Devices().SetEmulator(emulator)
	return nil
}

func setDiskPath(domain libvirtxml.Domain, diskPath string) error {
	disks := domain.Devices().Disks()
	if len(disks) != 1 {
		return errors.Errorf("multiple disks detected for domain '%s' - single disk domain templates are supported atm.", domain.Name())
	}

	disk := disks[0]
	disk.Source().SetFile(diskPath)

	return nil
}
