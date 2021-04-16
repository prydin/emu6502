package core

import (
	"fmt"
	"strings"
)

const (
	ModeImplied = iota
	ModeImmediate
	ModeZeroPage
	ModeZeroPageX
	ModeZeroPageY
	ModeAbsolute
	ModeAbsoluteX
	ModeAbsoluteY
	ModeIndirectX
	ModeIndirectY
	ModeRelative
	ModeIndirect
)

type Instruction struct {
	Mnemonic string
	AddrMode int
	Microcode []func()
}

func MkInstr(mnemonic string, microcode []func()) Instruction {
	parts := strings.Split(mnemonic, "_")
	mode := ModeImplied
	if len(parts) == 2 {
		switch parts[1] {
		case "I":
			mode = ModeImmediate
		case "Z":
			mode = ModeZeroPage
		case "ZX":
			mode = ModeZeroPageX
		case "ZY":
			mode = ModeZeroPageY
		case "A":
			mode = ModeAbsolute
		case "AX":
			mode = ModeAbsoluteX
		case "AY":
			mode = ModeAbsoluteY
		case "INDX":
			mode = ModeIndirectX
		case "INDY":
			mode = ModeIndirectY
		case "R":
			mode = ModeRelative
		case "IND":
			mode = ModeIndirect
		}
	}
	return Instruction{
		parts[0],
		mode,
		microcode,
	}
}

func (i *Instruction) Dissasemble(memory AddressSpace, pc uint16) string{
	s := i.Mnemonic + " "
	switch i.AddrMode {
	case ModeImmediate:
		s += fmt.Sprintf("#$%02x", memory.ReadByte(pc))
	case ModeZeroPage:
		s += fmt.Sprintf("$%02x", memory.ReadByte(pc))
	case ModeZeroPageX:
		s += fmt.Sprintf("$%02x,X", memory.ReadByte(pc))
	case ModeZeroPageY:
		s += fmt.Sprintf("$%02x,Y", memory.ReadByte(pc))
		case ModeAbsolute:
		s += fmt.Sprintf("$%02x,Y", memory.ReadByte(pc))
	case ModeAbsoluteX:
		s += fmt.Sprintf("$%04x,X", uint16(memory.ReadByte(pc)) + uint16(memory.ReadByte(pc+1) << 8))
	case ModeAbsoluteY:
		s += fmt.Sprintf("$%04x,Y", uint16(memory.ReadByte(pc)) + uint16(memory.ReadByte(pc+1) << 8))
	case ModeIndirectX:
		s += fmt.Sprintf("($%04x,X)", uint16(memory.ReadByte(pc)) + uint16(memory.ReadByte(pc+1) << 8))
	case ModeIndirectY:
		s += fmt.Sprintf("($%04x),Y", uint16(memory.ReadByte(pc)) + uint16(memory.ReadByte(pc+1) << 8))
	case ModeIndirect:
		s += fmt.Sprintf("($%04x)", uint16(memory.ReadByte(pc)) + uint16(memory.ReadByte(pc+1) << 8))
	case ModeRelative:
		s += fmt.Sprintf("%02x", memory.ReadByte(pc))
	}
	return s
}




