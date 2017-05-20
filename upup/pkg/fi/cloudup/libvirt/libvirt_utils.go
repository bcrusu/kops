package libvirt

import (
	"github.com/libvirt/libvirt-go"
	"github.com/pkg/errors"
	"k8s.io/kops/upup/pkg/fi/cloudup/libvirt/libvirtxml"
)

const uuidStringLength = 36

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
		return err
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
