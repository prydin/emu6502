/*
 * Copyright (c) 2021 Pontus Rydin
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the Software
 * without restriction, including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
 * to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

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
	Mnemonic  string
	AddrMode  int
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

func (i *Instruction) Dissasemble(memory AddressSpace, pc uint16) string {
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
		s += fmt.Sprintf("$%04x",  uint16(memory.ReadByte(pc))+uint16(memory.ReadByte(pc+1))<<8)
	case ModeAbsoluteX:
		s += fmt.Sprintf("$%04x,X", uint16(memory.ReadByte(pc))+uint16(memory.ReadByte(pc+1))<<8)
	case ModeAbsoluteY:
		s += fmt.Sprintf("$%04x,Y", uint16(memory.ReadByte(pc))+uint16(memory.ReadByte(pc+1))<<8)
	case ModeIndirectX:
		s += fmt.Sprintf("($%02x,X)", uint16(memory.ReadByte(pc)))
	case ModeIndirectY:
		s += fmt.Sprintf("($%02x),Y", uint16(memory.ReadByte(pc)))
	case ModeIndirect:
		s += fmt.Sprintf("($%04x)", uint16(memory.ReadByte(pc))+uint16(memory.ReadByte(pc+1)<<8))
	case ModeRelative:
		s += fmt.Sprintf("%02x", memory.ReadByte(pc))
	}
	return s
}
