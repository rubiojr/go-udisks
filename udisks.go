package udisks

import (
	"github.com/godbus/dbus/v5"

	"github.com/godbus/dbus/v5/introspect"
)

type Client struct {
	conn *dbus.Conn
}

type Drive struct {
	Vendor         string
	Model          string
	Serial         string
	Id             string
	MediaRemovable bool
	Ejectable      bool
	MediaAvailable bool
	ConnectionBus  string
	SiblingId      string
	Seat           string
	Removable      bool
	Size           uint64
}

type BlockDevice struct {
	UUID                string
	Device              string
	Id                  string
	Drive               *Drive
	Filesystems         []Filesystem
	CryptoBackingDevice *CryptoBackingDevice
}

type CryptoBackingDevice struct {
	Path                string
	CleartextDevicePath string
	HintEncryptionType  string
	MetadataSize        uint64
}

type Filesystem struct {
	MountPoints []string
	Size        uint64
}

func NewClient() (*Client, error) {
	c := &Client{}
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	c.conn = conn

	return c, nil
}

// BlockDevices returns the list of all block devices known to UDisks
func (c *Client) BlockDevices() ([]*BlockDevice, error) {
	conn := c.conn
	var list []string
	var filter map[string]interface{}
	obj := conn.Object("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2/Manager")
	err := obj.Call("org.freedesktop.UDisks2.Manager.GetBlockDevices", 0, &filter).Store(&list)
	if err != nil {
		return nil, err
	}

	bdevs := []*BlockDevice{}
	for _, bd := range list {
		dev := &BlockDevice{}
		bdevs = append(bdevs, dev)
		obj = conn.Object("org.freedesktop.UDisks2", dbus.ObjectPath(bd))
		dev.Device = bd
		stringProperty("org.freedesktop.UDisks2.Block.IdUUID", obj, &dev.UUID)
		stringProperty("org.freedesktop.UDisks2.Block.Id", obj, &dev.Id)

		var props map[string]dbus.Variant
		cbd, err := objGet(conn, "org.freedesktop.UDisks2.Block.CryptoBackingDevice", obj)
		if err == nil {
			cbd.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.freedesktop.UDisks2.Encrypted").Store(&props)
			if len(props) != 0 {
				dev.CryptoBackingDevice = &CryptoBackingDevice{}
				dev.CryptoBackingDevice.Path = string(cbd.Path())
				dev.CryptoBackingDevice.HintEncryptionType = props["HintEncryptionType"].Value().(string)
				dev.CryptoBackingDevice.MetadataSize = props["MetadataSize"].Value().(uint64)
			}
		}

		dev.Drive, err = c.getDrive(obj)

		obj.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.freedesktop.UDisks2.Filesystem").Store(&props)

		if len(props) != 0 {
			fs := Filesystem{}
			va := props["MountPoints"].Value()
			if va != nil {
				arr := va.([][]byte)
				for i := 0; i < len(arr); i++ {
					mpsa := arr[i]
					mpa := string(mpsa[0 : len(mpsa)-1])

					fs.MountPoints = append(fs.MountPoints, mpa)
				}
			}
			va = props["Size"].Value()
			if va != nil {
				fs.Size = va.(uint64)
			}

			dev.Filesystems = append(dev.Filesystems, fs)
		}
	}

	return bdevs, nil
}

// Drives returns the list of all block devices known to UDisks
func (c *Client) Drives() ([]*Drive, error) {
	drives := []*Drive{}
	node, err := introspect.Call(c.conn.Object("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2/drives"))
	if err != nil {
		return drives, err
	}

	for _, ch := range node.Children {
		path := "/org/freedesktop/UDisks2/drives/" + ch.Name
		obj := c.conn.Object("org.freedesktop.UDisks2", dbus.ObjectPath(path))
		drv, err := c.buildDrive(obj)
		if err != nil {
			return drives, err
		}
		drives = append(drives, drv)
	}

	return drives, nil
}
