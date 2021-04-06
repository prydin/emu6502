package core

type Memory interface {
	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)
}

type LinearMemory struct {
	bytes []uint8
}

func (l *LinearMemory) ReadByte(addr uint16) uint8 {
	if int(addr) < len(l.bytes) {
		return l.bytes[int(addr)]
	} else {
		return 0
	}
}

func (l *LinearMemory) WriteByte(addr uint16, data uint8) {
	if int(addr) < len(l.bytes) {
		l.bytes[int(addr)] = data
	}
}
