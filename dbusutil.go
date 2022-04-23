package udisks

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

func stringProperty(path string, obj dbus.BusObject, p *string) error {
	v, err := obj.GetProperty(path)
	if err != nil {
		return err
	}

	var ok bool
	*p, ok = v.Value().(string)
	if !ok {
		return err
	}

	return nil
}

func boolProperty(path string, obj dbus.BusObject, p *bool) error {
	v, err := obj.GetProperty(path)
	if err != nil {
		return err
	}

	var ok bool
	*p, ok = v.Value().(bool)
	if !ok {
		return err
	}

	return nil
}

func uint64Property(path string, obj dbus.BusObject, p *uint64) error {
	v, err := obj.GetProperty(path)
	if err != nil {
		return err
	}

	var ok bool
	*p, ok = v.Value().(uint64)
	if !ok {
		return err
	}

	return nil
}

func objGet(conn *dbus.Conn, p string, obj dbus.BusObject) (dbus.BusObject, error) {
	v, err := obj.GetProperty(p)
	if err != nil {
		return nil, err
	}

	var ok bool
	path, ok := v.Value().(dbus.ObjectPath)
	if !ok {
		return nil, fmt.Errorf("invalid object property")
	}
	driveObj := conn.Object("org.freedesktop.UDisks2", path)

	return driveObj, nil
}

func (c *Client) getDrive(blkobj dbus.BusObject) (*Drive, error) {
	objDrv, err := objGet(c.conn, "org.freedesktop.UDisks2.Block.Drive", blkobj)
	if err != nil || objDrv.Path() == dbus.ObjectPath("/") {
		return nil, ErrInvalidDrive
	}

	return c.buildDrive(objDrv)
}

func (c *Client) buildDrive(objDrv dbus.BusObject) (*Drive, error) {
	drv := &Drive{
		Vendor:         "",
		Model:          "",
		Serial:         "",
		Id:             "",
		MediaRemovable: false,
		Ejectable:      false,
		MediaAvailable: false,
	}
	stringProperty("org.freedesktop.UDisks2.Drive.Vendor", objDrv, &drv.Vendor)
	stringProperty("org.freedesktop.UDisks2.Drive.Serial", objDrv, &drv.Serial)
	stringProperty("org.freedesktop.UDisks2.Drive.Model", objDrv, &drv.Model)
	stringProperty("org.freedesktop.UDisks2.Drive.Id", objDrv, &drv.Id)
	stringProperty("org.freedesktop.UDisks2.Drive.ConnectionBus", objDrv, &drv.ConnectionBus)
	stringProperty("org.freedesktop.UDisks2.Drive.Seat", objDrv, &drv.Seat)
	stringProperty("org.freedesktop.UDisks2.Drive.SiblingId", objDrv, &drv.SiblingId)
	boolProperty("org.freedesktop.UDisks2.Drive.MediaRemovable", objDrv, &drv.MediaRemovable)
	boolProperty("org.freedesktop.UDisks2.Drive.MediaAvailable", objDrv, &drv.MediaAvailable)
	boolProperty("org.freedesktop.UDisks2.Drive.Ejectable", objDrv, &drv.Ejectable)
	boolProperty("org.freedesktop.UDisks2.Drive.Removable", objDrv, &drv.Removable)
	uint64Property("org.freedesktop.UDisks2.Drive.Size", objDrv, &drv.Size)

	return drv, nil
}
