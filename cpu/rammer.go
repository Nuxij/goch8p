package cpu

import (
	"fmt"
)

type Device interface {
	UUID() string
	Read(addr uint16) (byte, error)
	Reads(addr uint16, size uint16) ([]byte, error)
	Write(addr uint16, value byte) error
	Writes(addr uint16, values []byte) error
}

type RegionDevice struct {
	ID     string
	Size   uint16
	Offset uint16
}

type Region struct {
	Start   uint16
	Devices []RegionDevice
}

type InvalidRegionAlignment struct {
	alignment uint16
	addr      uint16
}

func (ira InvalidRegionAlignment) Error() string {
	return fmt.Sprintf("Invalid address %X for given alignement: %X", ira.addr, ira.alignment)
}

type Rammer struct {
	uuid      string
	alignment uint16
	devices   map[string]Device
	regions   map[uint16]Region
}

func NewRammer(alignment uint16, devices []Device) *Rammer {
	r := &Rammer{
		uuid:      fmt.Sprintf("Rammer::%s", RandomStringUUID()),
		alignment: alignment,
		devices:   make(map[string]Device),
		regions:   make(map[uint16]Region),
	}
	for _, device := range devices {
		r.devices[device.UUID()] = device
	}
	return r
}

func (r *Rammer) UUID() string {
	return r.uuid
}

func (r *Rammer) getRegionID(addr uint16) uint16 {
	return addr / r.alignment
}

func (r *Rammer) GetRegion(addr uint16) *Region {
	id := r.getRegionID(addr)
	if region, ok := r.regions[id]; ok {
		return &region
	}
	return nil
}

func (r *Rammer) SetRegion(start uint16, size uint16, device Device, deviceOffset uint16) error {
	if start != 0 && start%r.alignment != 0 {
		return InvalidRegionAlignment{r.alignment, start}
	}
	for i := uint16(start); i < start+size; i += r.alignment {
		r.regions[r.getRegionID(i)] = Region{
			Start: start,
			Devices: []RegionDevice{
				{device.UUID(), size, deviceOffset},
			},
		}
	}
	return nil
}

type AddressInvalid struct {
	addr uint16
}

func (a AddressInvalid) Error() string {
	return fmt.Sprintf("Invalid address %X", a.addr)
}

func (r *Rammer) checkBounds(addr uint16) (*Region, error) {
	if region := r.GetRegion(addr); region != nil {
		return region, nil
	}
	return nil, AddressInvalid{addr}
}

func (r *Rammer) Read(addr uint16) (byte, error) {
	if _, err := r.checkBounds(addr); err != nil {
		return 0, fmt.Errorf("cannot read: %w", err)
	} else {
		return r.readThrough(addr)
	}
}

func (r *Rammer) readThrough(addr uint16) (byte, error) {
	region := r.GetRegion(addr)
	if region == nil {
		return 0, AddressInvalid{addr}
	}
	lastDevice := region.Devices[len(region.Devices)-1]
	return r.devices[lastDevice.ID].Read(lastDevice.Offset + (addr - region.Start))
}

func (r *Rammer) Reads(addr uint16, size uint16) ([]byte, error) {
	if region, err := r.checkBounds(addr); err != nil {
		return nil, fmt.Errorf("cannot read: %w", err)
	} else {
		deviceAddr := region.Devices[0].Offset + (addr - region.Start)
		if deviceAddr+size > region.Devices[0].Size {
			size = region.Devices[0].Size - deviceAddr
		}
		return r.devices[region.Devices[0].ID].Reads(deviceAddr, size)
	}
}

func (r *Rammer) Write(addr uint16, value byte) error {
	if region, err := r.checkBounds(addr); err != nil {
		return fmt.Errorf("cannot write: %w", err)
	} else {
		return r.devices[region.Devices[0].ID].Write(region.Devices[0].Offset+(addr-region.Start), value)
	}
}

// TODO: Fix this to traverse the list, checking R/W to continue, or something
// func (r *Rammer) writeThrough(addr uint16, value byte) error {
// 	region := r.GetRegion(addr)
// 	if region == nil {
// 		return AddressInvalid{addr}
// 	}
// 	lastDevice := region.Devices[len(region.Devices)-1]
// 	return r.devices[lastDevice.ID].Write(lastDevice.Offset+(addr-region.Start), value)
// }

func (r *Rammer) Writes(addr uint16, values []byte) error {
	if region, err := r.checkBounds(addr); err != nil {
		return fmt.Errorf("cannot write: %w", err)
	} else {
		deviceAddr := region.Devices[0].Offset + (addr - region.Start)
		if deviceAddr+uint16(len(values)) > region.Devices[0].Size {
			values = values[:region.Devices[0].Size-deviceAddr]
		}
		return r.devices[region.Devices[0].ID].Writes(deviceAddr, values)
	}
}
