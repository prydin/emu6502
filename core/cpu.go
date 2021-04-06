package core

import "log"

const BRK = 0x00

// Accumulator instructions
const STA_A = 0x8d
const STA_Z = 0x85
const LDA_I = 0xa9
const LDA_Z = 0xa5
const LDA_A = 0xad
const LDA_ZX = 0xb5
const STA_ZX = 0x95
const LDA_AX = 0xbd
const LDA_AY = 0xb9
const LDA_INDX = 0xa1
const STA_AX = 0x9d
const STA_AY = 0x99
const STA_INDX = 0x81

// X-index instructions
const STX_A = 0x8e
const STX_Z = 0x86
const LDX_I = 0xa2
const LDX_Z = 0xa6
const LDX_A = 0xae
const LDX_ZY = 0xb6
const STX_ZY = 0x96

// Y-index instructions
const LDY_I = 0xa0
const LDY_Z = 0xa4
const LDY_ZX = 0xb4
const LDY_A = 0xac
const STY_A = 0x8c
const STY_ZX = 0x94
const STY_Z = 0x84

// CPU Status flags
const FLAG_C = uint8(0x01)
const FLAG_Z = uint8(0x02)
const FLAG_I = uint8(0x04)
const FLAG_D = uint8(0x08)
const FLAG_N = uint8(0x10)

const RST_VEC = 0xfffc

type CPU struct {
	// User accessible registers
	pc    uint16
	sp    uint8
	a     uint8
	x     uint8
	y     uint8
	flags uint8

	// Internal registers
	opcode             uint8  // Current instruction opcode
	operand            uint16 // Current operand address
	address            uint8 // Intermediate address storage during indirect addressing op
	halted             bool   // Halt CPU. Used for debugging
	irqPending         bool   // Handle IRQ after current instruction
	nmiPending         bool   // Handle NMO after current instruction
	CrashOnInvalidInst bool   // Used for debugging

	// Memory abstraction
	mem Memory

	// Pseudo-microcode
	microcode [][]func() // Pseudo-microcode
	microPc   int        // Microprogram counter
	fetchNext bool       // Fetch next instruction, please
}

func (c *CPU) Init(mem Memory) {
	c.microcode = make([][]func(), 256)

	// Basic memory access microcode
	fetch16Bits := []func(){c.fetchOperandLow, c.fetchOperandHigh}
	fetch8Bits := []func(){c.fetchOperandLow}

	// Addressing modes
	zeroPageX := []func(){c.fetchOperandLow, c.addXToLowOperand}
	zeroPageY := []func(){c.fetchOperandLow, c.addYToLowOperand}
	absXOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddX, c.nop}
	absYOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddY, c.nop}
	absX := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addXToOperand}
	absY := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addYToOperand}
	indirectX := []func(){c.fetchAddress, c.addXToAddress, c.fetchIndirectLow, c.fetchIndirectHigh}

	// Processor control instructions
	c.microcode[BRK] = []func(){c.brk}

	// Accumulator load/store
	c.microcode[LDA_A] = append(fetch16Bits, c.lda)
	c.microcode[LDA_I] = []func(){c.lda_i}
	c.microcode[LDA_ZX] = append(zeroPageX, c.lda)
	c.microcode[LDA_AX] = append(absXOverlap, c.lda)
	c.microcode[LDA_AY] = append(absYOverlap, c.lda)
	c.microcode[LDA_INDX] = append(indirectX, c.lda)
	c.microcode[LDA_Z] = append(fetch8Bits, c.lda)
	c.microcode[STA_A] = append(fetch16Bits, c.sta)
	c.microcode[STA_Z] = append(fetch8Bits, c.sta)
	c.microcode[STA_ZX] = append(zeroPageX, c.sta)
	c.microcode[STA_AX] = append(absX, c.sta)
	c.microcode[STA_AY] = append(absY, c.sta)
	c.microcode[STA_INDX] = append(indirectX, c.sta)


	// Index X load/store
	c.microcode[LDX_A] = append(fetch16Bits, c.ldx)
	c.microcode[LDX_I] = []func(){c.ldx_i}
	c.microcode[LDX_Z] = append(fetch8Bits, c.ldx)
	c.microcode[LDX_ZY] = append(zeroPageY, c.ldx)
	c.microcode[STX_A] = append(fetch16Bits, c.stx)
	c.microcode[STX_Z] = append(fetch8Bits, c.stx)
	c.microcode[STX_ZY] = append(zeroPageY, c.stx)

	// Index Y load/store
	c.microcode[LDY_I] = []func(){c.ldy_i}
	c.microcode[LDY_ZX] = append(zeroPageX, c.ldy)
	c.microcode[LDY_Z] = append(fetch8Bits, c.ldy)
	c.microcode[STY_A] = append(fetch16Bits, c.sty)
	c.microcode[STY_ZX] = append(zeroPageY, c.sty)
	c.microcode[STY_Z] = append(fetch8Bits, c.sty)
	c.microcode[LDY_A] = append(fetch16Bits, c.ldy)
	c.mem = mem
}

func (c *CPU) Reset() {
	c.flags = 0
	c.halted = false
	// TODO: What about SP?
	c.pc = (uint16(c.mem.ReadByte(RST_VEC))) + uint16(uint16(c.mem.ReadByte(RST_VEC+1))<<8)
	c.microPc = 0
	c.fetchNext = true
}

func (c *CPU) Clock() {
	if c.fetchNext {
		c.fetchOpcode()
		c.microPc = 0
		c.fetchNext = false
	} else {
		c.microcode[c.opcode][c.microPc]()
		c.microPc++
	}
}

func (c *CPU) fetchOpcode() {
	c.opcode = c.mem.ReadByte(c.pc)
	if c.CrashOnInvalidInst && len(c.microcode[c.opcode]) == 0 {
		log.Fatalf("Unknown opcode: %x", c.opcode)
	}
	c.pc++
}

func (c *CPU) fetchLow(target *uint16) {
	*target = uint16(c.mem.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) fetchHigh(target *uint16) {
	*target |= uint16(c.mem.ReadByte(c.pc)) << 8
	c.pc++
}

func (c *CPU) fetchOperandLow() {
	c.fetchLow(&c.operand)
}

func (c *CPU) fetchOperandHigh() {
	c.fetchHigh(&c.operand)
}

func (c *CPU) fetchAddress() {
	c.address = c.mem.ReadByte(c.pc)
	c.pc++
}

func (c *CPU) fetchOperandHighAndAdd(reg *uint8) {
	c.operand |= uint16(c.mem.ReadByte(c.pc)) << 8
	c.pc++
	t := c.operand + uint16(*reg)
	if t&0xff00 == c.operand&0xff00 {
		c.microPc++ // Skip extra clock cycle if it didn't cross page boundaries
	}
	c.operand = t
}

func (c *CPU) fetchOperandHighAndAddX() {
	c.fetchOperandHighAndAdd(&c.x)
}

func (c *CPU) fetchOperandHighAndAddY() {
	c.fetchOperandHighAndAdd(&c.y)
}

func (c *CPU) loadRegister(reg *uint8) {
	*reg = c.mem.ReadByte(c.operand)
	c.fetchNext = true
	c.updateNZ(*reg)
}

func (c *CPU) loadRegisterImmed(reg *uint8) {
	*reg = c.mem.ReadByte(c.pc)
	c.pc++
	c.fetchNext = true
	c.updateNZ(*reg)
}

func (c *CPU) storeRegister(reg *uint8) {
	c.mem.WriteByte(c.operand, *reg)
	c.fetchNext = true
}

func (c *CPU) fetchIndirectLow() {
	c.operand = uint16(c.mem.ReadByte(uint16(c.address)))
}

func (c *CPU) fetchIndirectHigh() {
	c.operand |= uint16(c.mem.ReadByte(uint16(c.address + 1)) << 8)
}

func (c *CPU) lda() {
	c.loadRegister(&c.a)
}

func (c *CPU) lda_i() {
	c.loadRegisterImmed(&c.a)
}

func (c *CPU) sta() {
	c.storeRegister(&c.a)
}

func (c *CPU) ldx() {
	c.loadRegister(&c.x)
}

func (c *CPU) ldx_i() {
	c.loadRegisterImmed(&c.x)
}

func (c *CPU) stx() {
	c.storeRegister(&c.x)
}

func (c *CPU) ldy() {
	c.loadRegister(&c.y)
}

func (c *CPU) ldy_i() {
	c.loadRegisterImmed(&c.y)
}

func (c *CPU) sty() {
	c.storeRegister(&c.y)
}

func (c *CPU) nop() {
}

func (c *CPU) updateNZ(b uint8) {
	if b == 0 {
		c.flags |= FLAG_Z
	} else {
		c.flags &= ^FLAG_Z
	}
	if b&0x80 == 0 {
		c.flags |= FLAG_C
	} else {
		c.flags &= ^FLAG_C
	}
}

func (c *CPU) addXToLowOperand() {
	c.operand = uint16(uint8(c.operand&0xff) + c.x)
}

func (c *CPU) addYToLowOperand() {
	c.operand = uint16(uint8(c.operand&0xff) + c.x)
}

func (c *CPU) addXToOperand() {
	c.operand = c.operand + uint16(c.x)
}

func (c *CPU) addYToOperand() {
	c.operand = c.operand + uint16(c.y)
}

func (c *CPU) addXToAddress() {
	c.address += c.x
}

func (c *CPU) brk() {
	c.halted = true
}

func (c *CPU) IsHalted() bool {
	return c.halted
}
