package core

import (
	"fmt"
	"log"
)

// Instructions
const (
	BRK = 0x00

	// Accumulator instructions
	STA_A    = 0x8d
	STA_Z    = 0x85
	LDA_I    = 0xa9
	LDA_Z    = 0xa5
	LDA_A    = 0xad
	LDA_ZX   = 0xb5
	STA_ZX   = 0x95
	LDA_AX   = 0xbd
	LDA_AY   = 0xb9
	LDA_INDX = 0xa1
	LDA_INDY = 0xb1
	STA_AX   = 0x9d
	STA_AY   = 0x99
	STA_INDX = 0x81
	STA_INDY = 0x91

	// X-index instructions
	STX_A  = 0x8e
	STX_Z  = 0x86
	LDX_I  = 0xa2
	LDX_Z  = 0xa6
	LDX_A  = 0xae
	LDX_ZY = 0xb6
	STX_ZY = 0x96
	INX    = 0xe8
	DEX    = 0xca

	// Y-index instructions
	LDY_I  = 0xa0
	LDY_Z  = 0xa4
	LDY_ZX = 0xb4
	LDY_A  = 0xac
	STY_A  = 0x8c
	STY_ZX = 0x94
	STY_Z  = 0x84
	INY    = 0xc8
	DEY    = 0x88

	// Arithmetic instructions
	INC_A  = 0xee
	INC_Z  = 0xe6
	INC_ZX = 0xf6
	INC_AX = 0xfe

	DEC_A  = 0Xce
	DEC_Z  = 0xc6
	DEC_ZX = 0xd6
	DEC_AX = 0xde

	// Jumps
	JMP     = 0x4c
	JMP_IND = 0x6c

	// Branches
	BCC = 0x90
	BCS = 0xb0
	BEQ = 0xf0
	BNE = 0xd0
	BMI = 0x30
	BPL = 0x10
	BVC = 0x50
	BVS = 0x70

	// Flag manipulation
	CLC = 0x18
	CLD = 0xd8
	CLI = 0x58
	CLV = 0xb8
	SEC = 0x38
	SED = 0xf8
	SEI = 0x78

	// Transfer instructions
	TAX = 0xaa
	TAY = 0xa8
	TXA = 0x8a
	TYA = 0x98
	TSX = 0xba
	TXS = 0x9a
)

// CPU Status flags
const (
	FLAG_C = uint8(0x01)
	FLAG_Z = uint8(0x02)
	FLAG_I = uint8(0x04)
	FLAG_D = uint8(0x08)
	FLAG_N = uint8(0x10)
	FLAG_V = uint8(0x20)
)

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
	address            uint8  // Intermediate address storage during indirect addressing op
	alu                uint8  // ALU internal accumulator
	halted             bool   // Halt CPU. Used for debugging
	irqPending         bool   // Handle IRQ after current instruction
	nmiPending         bool   // Handle NMO after current instruction
	CrashOnInvalidInst bool   // Used for debugging
	Trace              bool   // Trace each instruction to stdout

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

	// Addressing mode microcode
	zeroPageX := []func(){c.fetchOperandLow, c.addXToLowOperand}
	zeroPageY := []func(){c.fetchOperandLow, c.addYToLowOperand}
	absXOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddX, c.nop}
	absYOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddY, c.nop}
	absX := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addXToOperand}
	absY := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addYToOperand}
	indirectX := []func(){c.fetchAddressLow, c.addXToAddress, c.fetchIndirectLow, c.fetchIndirectHigh}
	indirectY := []func(){c.fetchAddressLow, c.fetchIndirectLow, c.fetchIndirectHighAndAddY, c.nop}

	// Processor control instructions
	c.microcode[BRK] = []func(){c.brk}

	// Accumulator load/store
	c.microcode[LDA_A] = append(fetch16Bits, c.lda)
	c.microcode[LDA_I] = []func(){c.lda_i}
	c.microcode[LDA_ZX] = append(zeroPageX, c.lda)
	c.microcode[LDA_AX] = append(absXOverlap, c.lda)
	c.microcode[LDA_AY] = append(absYOverlap, c.lda)
	c.microcode[LDA_INDX] = append(indirectX, c.lda)
	c.microcode[LDA_INDY] = append(indirectY, c.lda)
	c.microcode[LDA_Z] = append(fetch8Bits, c.lda)
	c.microcode[STA_A] = append(fetch16Bits, c.sta)
	c.microcode[STA_Z] = append(fetch8Bits, c.sta)
	c.microcode[STA_ZX] = append(zeroPageX, c.sta)
	c.microcode[STA_AX] = append(absX, c.sta)
	c.microcode[STA_AY] = append(absY, c.sta)
	c.microcode[STA_INDX] = append(indirectX, c.sta)
	c.microcode[STA_INDY] = append(indirectY, c.sta)

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

	// Inc/dec register
	c.microcode[INX] = []func(){c.inx}
	c.microcode[INY] = []func(){c.iny}
	c.microcode[DEX] = []func(){c.dex}
	c.microcode[DEY] = []func(){c.dey}

	// INC
	c.microcode[INC_Z] = append(fetch8Bits, c.loadALU, c.inc, c.storeALU)
	c.microcode[INC_ZX] = append(zeroPageX, c.loadALU, c.inc, c.storeALU)
	c.microcode[INC_A] = append(fetch16Bits, c.loadALU, c.inc, c.storeALU)
	c.microcode[INC_AX] = append(absX, c.loadALU, c.inc, c.storeALU)

	// DEC
	c.microcode[DEC_Z] = append(fetch8Bits, c.loadALU, c.dec, c.storeALU)
	c.microcode[DEC_ZX] = append(zeroPageX, c.loadALU, c.dec, c.storeALU)
	c.microcode[DEC_A] = append(fetch16Bits, c.loadALU, c.dec, c.storeALU)
	c.microcode[DEC_AX] = append(absX, c.loadALU, c.dec, c.storeALU)

	// JMP
	c.microcode[JMP] = []func(){c.fetchOperandLow, c.fetchHighAndJump}
	c.microcode[JMP_IND] = []func(){c.fetchOperandLow, c.fetchOperandHigh, c.loadPCLow, c.loadPCHigh}

	// Branching
	c.microcode[BCC] = append(fetch8Bits, c.bcc)
	c.microcode[BCS] = append(fetch8Bits, c.bcs)
	c.microcode[BEQ] = append(fetch8Bits, c.beq)
	c.microcode[BNE] = append(fetch8Bits, c.bne)
	c.microcode[BMI] = append(fetch8Bits, c.bmi)
	c.microcode[BPL] = append(fetch8Bits, c.bpl)
	c.microcode[BVC] = append(fetch8Bits, c.bvc)
	c.microcode[BVS] = append(fetch8Bits, c.bvs)

	// Flag manipulations
	c.microcode[CLC] = []func(){c.clc}
	c.microcode[CLD] = []func(){c.cld}
	c.microcode[CLV] = []func(){c.clv}
	c.microcode[CLI] = []func(){c.clc}
	c.microcode[SEC] = []func(){c.sec}
	c.microcode[SED] = []func(){c.sed}
	c.microcode[SEI] = []func(){c.sei}

	// Transfer instructions
	c.microcode[TAX] = []func(){c.tax}
	c.microcode[TAY] = []func(){c.tay}
	c.microcode[TSX] = []func(){c.tsx}
	c.microcode[TXA] = []func(){c.txa}
	c.microcode[TYA] = []func(){c.tya}
	c.microcode[TXS] = []func(){c.txs}
	c.mem = mem
}

func (c *CPU) Reset() {
	c.flags = 0
	c.halted = false
	// TODO: What about SP?
	c.pc = (uint16(c.mem.ReadByte(RST_VEC))) + uint16(c.mem.ReadByte(RST_VEC+1))<<8
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
	if c.Trace {
		fmt.Printf("PC=%04x [PC]=%02x MPC=%02x OP=%02x SP=%04x A=%02x X=%02x Y=%02x Flags=%02x Oper=%04x\n",
			c.pc, c.mem.ReadByte(c.pc), c.microPc, c.opcode, c.sp, c.a, c.x, c.y, c.flags, c.operand)
	}
}

func (c *CPU) fetchOpcode() {
	c.opcode = c.mem.ReadByte(c.pc)
	if c.CrashOnInvalidInst && len(c.microcode[c.opcode]) == 0 {
		log.Fatalf("Unknown opcode: %2x at address %4x", c.opcode, c.pc)
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

func (c *CPU) fetchHighAndJump() {
	c.fetchOperandHigh()
	c.pc = c.operand
	c.fetchNext = true
}

func (c *CPU) fetchOperandLow() {
	c.fetchLow(&c.operand)
}

func (c *CPU) fetchOperandHigh() {
	c.fetchHigh(&c.operand)
}

func (c *CPU) fetchAddressLow() {
	c.address = c.mem.ReadByte(c.pc)
	c.pc++
}

func (c *CPU) fetchAddressHigh() {
	c.address = c.mem.ReadByte(c.pc)
	c.pc++
}

func (c *CPU) loadPCLow() {
	c.pc = uint16(c.mem.ReadByte(c.operand))
}

func (c *CPU) loadPCHigh() {
	c.pc |= uint16(c.mem.ReadByte(c.operand+1)) << 8
	c.fetchNext = true
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
	c.operand |= uint16(c.mem.ReadByte(uint16(c.address+1))) << 8
}

func (c *CPU) fetchIndirectHighAndAddY() {
	c.operand |= uint16(c.mem.ReadByte(uint16(c.address+1))) << 8
	t := c.operand + uint16(c.y)
	if t&0xf0 == c.operand&0xf0 {
		c.microPc++ // Skip extra clock cycle
	}
	c.operand = t
}

func (c *CPU) branchIf(mask, wanted uint8) {
	if c.flags&mask == wanted {
		if c.operand >= 0x80 {
			c.pc -= uint16(^uint8(c.operand) + 1) // 2s complement
		} else {
			c.pc += c.operand
		}
	}
	c.fetchNext = true
}

func (c *CPU) bne() {
	c.branchIf(FLAG_Z, 0)
}

func (c *CPU) beq() {
	c.branchIf(FLAG_Z, FLAG_Z)
}

func (c *CPU) bcc() {
	c.branchIf(FLAG_C, 0)
}

func (c *CPU) bcs() {
	c.branchIf(FLAG_C, FLAG_C)
}

func (c *CPU) bvc() {
	c.branchIf(FLAG_V, 0)
}

func (c *CPU) bvs() {
	c.branchIf(FLAG_V, FLAG_V)
}

func (c *CPU) bpl() {
	c.branchIf(FLAG_N, 0)
}

func (c *CPU) bmi() {
	c.branchIf(FLAG_N, FLAG_N)
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

func (c *CPU) inx() {
	c.x++
	c.updateNZ(c.x)
	c.fetchNext = true
}

func (c *CPU) iny() {
	c.y++
	c.updateNZ(c.y)
	c.fetchNext = true
}

func (c *CPU) dex() {
	c.x--
	c.updateNZ(c.x)
	c.fetchNext = true
}

func (c *CPU) dey() {
	c.y--
	c.updateNZ(c.y)
	c.fetchNext = true
}

func (c *CPU) inc() {
	c.alu++
	c.updateNZ(c.alu)
}

func (c *CPU) dec() {
	c.alu--
	c.updateNZ(c.alu)
}

func (c *CPU) clc() {
	c.flags &= ^FLAG_C
	c.fetchNext = true
}

func (c *CPU) clv() {
	c.flags &= ^FLAG_V
	c.fetchNext = true
}

func (c *CPU) cld() {
	c.flags &= ^FLAG_D
	c.fetchNext = true
}

func (c *CPU) cli() {
	c.flags &= ^FLAG_I
	c.fetchNext = true
}

func (c *CPU) sec() {
	c.flags |= FLAG_C
	c.fetchNext = true
}

func (c *CPU) sed() {
	c.flags |= FLAG_D
	c.fetchNext = true
}

func (c *CPU) sei() {
	c.flags |= FLAG_I
	c.fetchNext = true
}

func (c *CPU) tax() {
	c.a = c.x
	c.updateNZ(c.x)
	c.fetchNext = true
}

func (c *CPU) tay() {
	c.y = c.a
	c.updateNZ(c.y)
	c.fetchNext = true
}

func (c *CPU) tsx() {
	c.x = c.sp
	c.updateNZ(c.x)
	c.fetchNext = true
}

func (c *CPU) txa() {
	c.a = c.x
	c.updateNZ(c.a)
	c.fetchNext = true
}

func (c *CPU) txs() {
	c.sp = c.x
	c.updateNZ(c.sp)
	c.fetchNext = true
}

func (c *CPU) tya() {
	c.a = c.y
	c.updateNZ(c.y)
	c.fetchNext = true
}

func (c *CPU) nop() {
}

func (c *CPU) updateNZ(b uint8) {
	if b == 0 {
		c.flags |= FLAG_Z
	} else {
		c.flags &= ^FLAG_Z
	}
	if b&0x80 != 0 {
		c.flags |= FLAG_N
	} else {
		c.flags &= ^FLAG_N
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

func (c *CPU) loadALU() {
	c.alu = c.mem.ReadByte(c.operand)
}

func (c *CPU) storeALU() {
	c.mem.WriteByte(c.operand, c.alu)
	c.fetchNext = true // This is always the last microinstruction in a sequence!
}

func (c *CPU) brk() {
	c.halted = true
}

func (c *CPU) IsHalted() bool {
	return c.halted
}
