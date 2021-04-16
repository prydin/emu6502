package core

type AddressSpace interface {
	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)
}

type ROM struct {
	bytes []uint8
}

type RAM struct {
	bytes []uint8
}

type Device struct {
	start uint16
	end uint16
	device AddressSpace
}

type Bus struct {
	devices []Device
}

func (b *Bus) Connect(device AddressSpace, start, end uint16) {
	b.devices = append(b.devices, Device{ start, end, device })
}

func (b *Bus) ReadByte(addr uint16) uint8 {
	for _, d := range b.devices {
		if addr >= d.start && addr <= d.end {
			return d.device.ReadByte(addr - d.start)
		}
	}
	return 0
}

func (b *Bus) WriteByte(addr uint16, data uint8) {
	for _, d := range b.devices {
		if addr >= d.start && addr <= d.end {
			d.device.WriteByte(addr - d.start, data)
		}
	}
}

func (r *ROM) ReadByte(addr uint16) uint8 {
	if int(addr) < len(r.bytes) {
		return r.bytes[int(addr)]
	} else {
		return 0
	}
}

func (r *ROM) WriteByte(addr uint16, data uint8) {
}

func (r *RAM) ReadByte(addr uint16) uint8 {
	if int(addr) < len(r.bytes) {
		return r.bytes[int(addr)]
	} else {
		return 0
	}
}

func (l *RAM) WriteByte(addr uint16, data uint8) {
	if int(addr) < len(l.bytes) {
		l.bytes[int(addr)] = data
	}
}


