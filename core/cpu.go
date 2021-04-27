package core

import (
	"fmt"
	"log"
)

// Instructions
const (
	BRK = 0x00
	NOP = 0xea

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
	LDX_AY = 0xbe
	LDX_ZY = 0xb6
	STX_ZY = 0x96
	INX    = 0xe8
	DEX    = 0xca

	// Y-index instructions
	LDY_I  = 0xa0
	LDY_Z  = 0xa4
	LDY_ZX = 0xb4
	LDY_A  = 0xac
	LDY_AX = 0xbc
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
	JSR_A   = 0x20
	RTS     = 0x60
	RTI     = 0x40

	// Branches
	BCC_R = 0x90
	BCS_R = 0xb0
	BEQ_R = 0xf0
	BNE_R = 0xd0
	BMI_R = 0x30
	BPL_R = 0x10
	BVC_R = 0x50
	BVS_R = 0x70

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

	// Stack instructions
	PLA = 0x68
	PHA = 0x48
	PLP = 0x28
	PHP = 0x08

	// Arithmetic
	ADC_I    = 0x69
	ADC_Z    = 0x65
	ADC_ZX   = 0x75
	ADC_A    = 0x6d
	ADC_AX   = 0x7d
	ADC_AY   = 0x79
	ADC_INDX = 0x61
	ADC_INDY = 0x71
	SBC_I    = 0xe9
	SBC_Z    = 0xe5
	SBC_ZX   = 0xf5
	SBC_A    = 0xed
	SBC_AX   = 0xfd
	SBC_AY   = 0xf9
	SBC_INDX = 0xe1
	SBC_INDY = 0xf1

	// Logic
	AND_I    = 0x29
	AND_Z    = 0x25
	AND_ZX   = 0x35
	AND_A    = 0x2d
	AND_AX   = 0x3d
	AND_AY   = 0x39
	AND_INDX = 0x21
	AND_INDY = 0x31
	ORA_I    = 0x09
	ORA_Z    = 0x05
	ORA_ZX   = 0x15
	ORA_A    = 0x0d
	ORA_AX   = 0x1d
	ORA_AY   = 0x19
	ORA_INDX = 0x01
	ORA_INDY = 0x11
	EOR_I    = 0x49
	EOR_Z    = 0x45
	EOR_ZX   = 0x55
	EOR_A    = 0x4d
	EOR_AX   = 0x5d
	EOR_AY   = 0x59
	EOR_INDX = 0x41
	EOR_INDY = 0x51
	ASL_ACC  = 0x0a
	ASL_Z    = 0x06
	ASL_ZX   = 0x16
	ASL_A    = 0x0e
	ASL_AX   = 0x1e
	LSR_ACC  = 0x4a
	LSR_Z    = 0x46
	LSR_ZX   = 0x56
	LSR_A    = 0x4e
	LSR_AX   = 0x5e
	ROL_ACC  = 0x2a
	ROL_Z    = 0x26
	ROL_ZX   = 0x36
	ROL_A    = 0x2e
	ROL_AX   = 0x3e
	ROR_ACC  = 0x6a
	ROR_Z    = 0x66
	ROR_ZX   = 0x76
	ROR_A    = 0x6e
	ROR_AX   = 0x7e
	BIT_Z    = 0x24
	BIT_A    = 0x2c

	// Comparisons
	CMP_I    = 0xc9
	CMP_Z    = 0xc5
	CMP_ZX   = 0xd5
	CMP_A    = 0xcd
	CMP_AX   = 0xdd
	CMP_AY   = 0xd9
	CMP_INDX = 0xc1
	CMP_INDY = 0xd1
	CPX_I    = 0xe0
	CPX_Z    = 0xe4
	CPX_A    = 0xec
	CPY_I    = 0xc0
	CPY_Z    = 0xc4
	CPY_A    = 0xcc
)

// CPU Status flags
const (
	FLAG_C = uint8(0x01)
	FLAG_Z = uint8(0x02)
	FLAG_I = uint8(0x04)
	FLAG_D = uint8(0x08)
	FLAG_B = uint8(0x10)
	FLAG_U = uint8(0x20) // Unused, always 1
	FLAG_V = uint8(0x40)
	FLAG_N = uint8(0x80)
)

const NMI_VEC = 0xfffa
const RST_VEC = 0xfffc
const IRQ_VEC = 0xfffe

type CPU struct {
	// User accessible registers
	pc    uint16
	sp    uint8
	a     uint8
	x     uint8
	y     uint8
	flags uint8

	// Internal registers
	operand            uint16       // Current operand address
	address            uint8        // Intermediate address storage during indirect addressing op
	alu                uint8        // ALU internal accumulator
	halted             bool         // Halt CPU. Used for debugging
	irqPending         bool         // Handle IRQ after current instruction
	nmiPending         bool         // Handle NMI after current instruction
	CrashOnInvalidInst bool         // Used for debugging
	HaltOnBRK          bool         // Used for debugging
	Trace              bool         // Trace each instruction to stdout
	instruction        *Instruction // Current instruction

	// AddressSpace abstraction
	bus *Bus

	// Pseudo-instructionSet
	instructionSet []Instruction // Pseudo-instructionSet
	microPc        int           // Microprogram counter

	// Interrupt pseudo instructions
	irqPI Instruction
	nmiPI Instruction
	brkPI Instruction
	rstPI Instruction
}

func (c *CPU) Init(bus *Bus) {
	c.instructionSet = make([]Instruction, 256)
	c.bus = bus

	// Basic memory access instructionSet
	fetch16Bits := []func(){c.fetchOperandLow, c.fetchOperandHigh}
	fetch8Bits := []func(){c.fetchOperandLow}

	// Addressing mode instructionSet
	zeroPageX := []func(){c.fetchOperandLow, c.addXToLowOperand}
	zeroPageY := []func(){c.fetchOperandLow, c.addYToLowOperand}
	absXOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddX, c.nop}
	absYOverlap := []func(){c.fetchOperandLow, c.fetchOperandHighAndAddY, c.nop}
	absX := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addXToOperand}
	absY := []func(){c.fetchOperandLow, c.fetchOperandHigh, c.addYToOperand}
	indirectX := []func(){c.fetchAddressLow, c.addXToAddress, c.fetchIndirectLow, c.fetchIndirectHigh}
	indirectY := []func(){c.fetchAddressLow, c.fetchIndirectLow, c.fetchIndirectHighAndAddY, c.nop}

	// Processor control instructions
	if c.HaltOnBRK {
		c.instructionSet[BRK] = MkInstr("BRK", []func(){func() { c.halted = true }})
	} else {
		c.instructionSet[BRK] = MkInstr("BRK", []func(){
			func() {
				c.operand = uint16(c.bus.ReadByte(IRQ_VEC))
			},
			func() {
				c.operand |= uint16(c.bus.ReadByte(IRQ_VEC+1)) << 8
				c.pc++ // BRK is a two-byte instruction
			},
			c.pushInterruptReturnAddressHigh,
			c.pushInterruptReturnAddressLow,
			c.php_brk,
			c.jump,
		})
	}
	c.instructionSet[NOP] = MkInstr("NOP", []func(){c.nop})

	// Accumulator load/store
	c.instructionSet[LDA_A] = MkInstr("LDA_A", append(fetch16Bits, c.lda))
	c.instructionSet[LDA_I] = MkInstr("LDA_I", []func(){c.lda_i})
	c.instructionSet[LDA_ZX] = MkInstr("LDA_ZX", append(zeroPageX, c.lda))
	c.instructionSet[LDA_AX] = MkInstr("LDA_AX", append(absXOverlap, c.lda))
	c.instructionSet[LDA_AY] = MkInstr("LDA_AY", append(absYOverlap, c.lda))
	c.instructionSet[LDA_INDX] = MkInstr("LDA_INDX", append(indirectX, c.lda))
	c.instructionSet[LDA_INDY] = MkInstr("LDA_INDY", append(indirectY, c.lda))
	c.instructionSet[LDA_Z] = MkInstr("LDA_Z", append(fetch8Bits, c.lda))
	c.instructionSet[STA_A] = MkInstr("STA_A", append(fetch16Bits, c.sta))
	c.instructionSet[STA_Z] = MkInstr("STA_Z", append(fetch8Bits, c.sta))
	c.instructionSet[STA_ZX] = MkInstr("STA_ZX", append(zeroPageX, c.sta))
	c.instructionSet[STA_AX] = MkInstr("STA_AX", append(absX, c.sta))
	c.instructionSet[STA_AY] = MkInstr("STA_AY", append(absY, c.sta))
	c.instructionSet[STA_INDX] = MkInstr("STA_INDX", append(indirectX, c.sta))
	c.instructionSet[STA_INDY] = MkInstr("STA_INDY", append(indirectY, c.sta))

	// Index X load/store
	c.instructionSet[LDX_A] = MkInstr("LDX_A", append(fetch16Bits, c.ldx))
	c.instructionSet[LDX_AY] = MkInstr("LDX_AY", append(absYOverlap, c.ldx))
	c.instructionSet[LDX_I] = MkInstr("LDX_I", []func(){c.ldx_i})
	c.instructionSet[LDX_Z] = MkInstr("LDX_Z", append(fetch8Bits, c.ldx))
	c.instructionSet[LDX_ZY] = MkInstr("LDX_ZY", append(zeroPageY, c.ldx))
	c.instructionSet[STX_A] = MkInstr("STX_A", append(fetch16Bits, c.stx))
	c.instructionSet[STX_Z] = MkInstr("STX_Z", append(fetch8Bits, c.stx))
	c.instructionSet[STX_ZY] = MkInstr("STX_ZY", append(zeroPageY, c.stx))

	// Index Y load/store
	c.instructionSet[LDY_I] = MkInstr("LDY_I", []func(){c.ldy_i})
	c.instructionSet[LDY_ZX] = MkInstr("LDY_ZX", append(zeroPageX, c.ldy))
	c.instructionSet[LDY_Z] = MkInstr("LDY_Z", append(fetch8Bits, c.ldy))
	c.instructionSet[STY_A] = MkInstr("STY_A", append(fetch16Bits, c.sty))
	c.instructionSet[STY_ZX] = MkInstr("STY_ZX", append(zeroPageX, c.sty))
	c.instructionSet[STY_Z] = MkInstr("STY_Z", append(fetch8Bits, c.sty))
	c.instructionSet[LDY_A] = MkInstr("LDY_A", append(fetch16Bits, c.ldy))
	c.instructionSet[LDY_AX] = MkInstr("LDY_AX", append(absXOverlap, c.ldy))

	// Inc/dec register
	c.instructionSet[INX] = MkInstr("INX", []func(){c.inx})
	c.instructionSet[INY] = MkInstr("INY", []func(){c.iny})
	c.instructionSet[DEX] = MkInstr("DEX", []func(){c.dex})
	c.instructionSet[DEY] = MkInstr("DEY", []func(){c.dey})

	// INC
	c.instructionSet[INC_Z] = MkInstr("INC_Z", append(fetch8Bits, c.loadALU, c.inc, c.storeALU))
	c.instructionSet[INC_ZX] = MkInstr("INC_ZX", append(zeroPageX, c.loadALU, c.inc, c.storeALU))
	c.instructionSet[INC_A] = MkInstr("INC_A", append(fetch16Bits, c.loadALU, c.inc, c.storeALU))
	c.instructionSet[INC_AX] = MkInstr("INC_AX", append(absX, c.loadALU, c.inc, c.storeALU))

	// DEC
	c.instructionSet[DEC_Z] = MkInstr("DEC_Z", append(fetch8Bits, c.loadALU, c.dec, c.storeALU))
	c.instructionSet[DEC_ZX] = MkInstr("DEC_ZX", append(zeroPageX, c.loadALU, c.dec, c.storeALU))
	c.instructionSet[DEC_A] = MkInstr("DEC_A", append(fetch16Bits, c.loadALU, c.dec, c.storeALU))
	c.instructionSet[DEC_AX] = MkInstr("DEC_AX", append(absX, c.loadALU, c.dec, c.storeALU))

	// JMP
	c.instructionSet[JMP] = MkInstr("JMP", []func(){c.fetchOperandLow, c.fetchHighAndJump})
	c.instructionSet[JMP_IND] = MkInstr("JMP_IND", []func(){c.fetchOperandLow, c.fetchOperandHigh, c.loadPCLow, c.loadPCHigh})

	// JSR/RTS
	c.instructionSet[JSR_A] = MkInstr("JSR_A", append(fetch16Bits, c.pushReturnAddressHigh, c.pushReturnAddressLow, c.jump))
	c.instructionSet[RTS] = MkInstr("RTS", []func(){c.nop, c.pullOperandLow, c.pullOperandHigh, c.jump, c.incPc})
	c.instructionSet[RTI] = MkInstr("RTI", []func(){c.plp, c.pullOperandLow, c.pullOperandHigh, c.jump, c.nop})

	// Branching
	c.instructionSet[BCC_R] = MkInstr("BCC_R", []func(){c.bcc, c.doBranch, c.nop})
	c.instructionSet[BCS_R] = MkInstr("BCS_R", []func(){c.bcs, c.doBranch, c.nop})
	c.instructionSet[BEQ_R] = MkInstr("BEQ_R", []func(){c.beq, c.doBranch, c.nop})
	c.instructionSet[BNE_R] = MkInstr("BNE_R", []func(){c.bne, c.doBranch, c.nop})
	c.instructionSet[BMI_R] = MkInstr("BMI_R", []func(){c.bmi, c.doBranch, c.nop})
	c.instructionSet[BPL_R] = MkInstr("BPL_R", []func(){c.bpl, c.doBranch, c.nop})
	c.instructionSet[BVC_R] = MkInstr("BVC_R", []func(){c.bvc, c.doBranch, c.nop})
	c.instructionSet[BVS_R] = MkInstr("BVS_R", []func(){c.bvs, c.doBranch, c.nop})

	// Flag manipulations
	c.instructionSet[CLC] = MkInstr("CLC", []func(){c.clc})
	c.instructionSet[CLD] = MkInstr("CLD", []func(){c.cld})
	c.instructionSet[CLV] = MkInstr("CLV", []func(){c.clv})
	c.instructionSet[CLI] = MkInstr("CLI", []func(){c.cli})
	c.instructionSet[SEC] = MkInstr("SEC", []func(){c.sec})
	c.instructionSet[SED] = MkInstr("SED", []func(){c.sed})
	c.instructionSet[SEI] = MkInstr("SEI", []func(){c.sei})

	// Transfer instructions.
	c.instructionSet[TAX] = MkInstr("TAX", []func(){c.tax})
	c.instructionSet[TAY] = MkInstr("TAY", []func(){c.tay})
	c.instructionSet[TSX] = MkInstr("TSX", []func(){c.tsx})
	c.instructionSet[TXA] = MkInstr("TXA", []func(){c.txa})
	c.instructionSet[TYA] = MkInstr("TYA", []func(){c.tya})
	c.instructionSet[TXS] = MkInstr("TXS", []func(){c.txs})

	// Stack instructions
	// The NOPs are a bit of a cheat to get the instruction timing right.
	// The bus timing is still correct.
	c.instructionSet[PHA] = MkInstr("PHA", []func(){c.nop, c.pha})
	c.instructionSet[PHP] = MkInstr("PHP", []func(){c.nop, c.php})
	c.instructionSet[PLA] = MkInstr("PLA", []func(){c.nop, c.pla})
	c.instructionSet[PLP] = MkInstr("PLP", []func(){c.nop, c.plp})

	// Arithmetic
	c.instructionSet[ADC_A] = MkInstr("ADC_A", append(fetch16Bits, c.adc))
	c.instructionSet[ADC_I] = MkInstr("ADC_I", []func(){c.adc_i})
	c.instructionSet[ADC_ZX] = MkInstr("ADC_ZX", append(zeroPageX, c.adc))
	c.instructionSet[ADC_AX] = MkInstr("ADC_AX", append(absXOverlap, c.adc))
	c.instructionSet[ADC_AY] = MkInstr("ADC_AY", append(absYOverlap, c.adc))
	c.instructionSet[ADC_INDX] = MkInstr("ADC_INDX", append(indirectX, c.adc))
	c.instructionSet[ADC_INDY] = MkInstr("ADC_INDY", append(indirectY, c.adc))
	c.instructionSet[ADC_Z] = MkInstr("ADC_Z", append(fetch8Bits, c.adc))

	c.instructionSet[SBC_A] = MkInstr("SBC_A", append(fetch16Bits, c.sbc))
	c.instructionSet[SBC_I] = MkInstr("SBC_I", []func(){c.sbc_i})
	c.instructionSet[SBC_ZX] = MkInstr("SBC_ZX", append(zeroPageX, c.sbc))
	c.instructionSet[SBC_AX] = MkInstr("SBC_AX", append(absXOverlap, c.sbc))
	c.instructionSet[SBC_AY] = MkInstr("SBC_AY", append(absYOverlap, c.sbc))
	c.instructionSet[SBC_INDX] = MkInstr("SBC_INDX", append(indirectX, c.sbc))
	c.instructionSet[SBC_INDY] = MkInstr("SBC_INDY", append(indirectY, c.sbc))
	c.instructionSet[SBC_Z] = MkInstr("SBC_Z", append(fetch8Bits, c.sbc))

	// Logic
	c.instructionSet[AND_A] = MkInstr("AND_A", append(fetch16Bits, c.and))
	c.instructionSet[AND_I] = MkInstr("AND_I", []func(){c.and_i})
	c.instructionSet[AND_ZX] = MkInstr("AND_ZX", append(zeroPageX, c.and))
	c.instructionSet[AND_AX] = MkInstr("AND_AX", append(absXOverlap, c.and))
	c.instructionSet[AND_AY] = MkInstr("AND_AY", append(absYOverlap, c.and))
	c.instructionSet[AND_INDX] = MkInstr("AND_INDX", append(indirectX, c.and))
	c.instructionSet[AND_INDY] = MkInstr("AND_INDY", append(indirectY, c.and))
	c.instructionSet[AND_Z] = MkInstr("AND_Z", append(fetch8Bits, c.and))
	c.instructionSet[ORA_A] = MkInstr("ORA_A", append(fetch16Bits, c.ora))
	c.instructionSet[ORA_I] = MkInstr("ORA_I", []func(){c.ora_i})
	c.instructionSet[ORA_ZX] = MkInstr("ORA_ZX", append(zeroPageX, c.ora))
	c.instructionSet[ORA_AX] = MkInstr("ORA_AX", append(absXOverlap, c.ora))
	c.instructionSet[ORA_AY] = MkInstr("ORA_AY", append(absYOverlap, c.ora))
	c.instructionSet[ORA_INDX] = MkInstr("ORA_INDX", append(indirectX, c.ora))
	c.instructionSet[ORA_INDY] = MkInstr("ORA_INDY", append(indirectY, c.ora))
	c.instructionSet[ORA_Z] = MkInstr("ORA_Z", append(fetch8Bits, c.ora))
	c.instructionSet[EOR_A] = MkInstr("EOR_A", append(fetch16Bits, c.eor))
	c.instructionSet[EOR_I] = MkInstr("EOR_I", []func(){c.eor_i})
	c.instructionSet[EOR_ZX] = MkInstr("EOR_ZX", append(zeroPageX, c.eor))
	c.instructionSet[EOR_AX] = MkInstr("EOR_AX", append(absXOverlap, c.eor))
	c.instructionSet[EOR_AY] = MkInstr("EOR_AY", append(absYOverlap, c.eor))
	c.instructionSet[EOR_INDX] = MkInstr("EOR_INDX", append(indirectX, c.eor))
	c.instructionSet[EOR_INDY] = MkInstr("EOR_INDY", append(indirectY, c.eor))
	c.instructionSet[EOR_Z] = MkInstr("EOR_Z", append(fetch8Bits, c.eor))
	c.instructionSet[BIT_Z] = MkInstr("BIT_Z", append(fetch8Bits, c.bit))
	c.instructionSet[BIT_A] = MkInstr("BIT_A", append(fetch16Bits, c.bit))
	c.instructionSet[ASL_ACC] = MkInstr("ASL_ACC", []func(){c.asl_acc})
	c.instructionSet[ASL_Z] = MkInstr("ASL_Z", append(fetch8Bits, c.loadALU, c.asl_alu, c.storeALU))
	c.instructionSet[ASL_ZX] = MkInstr("ASL_ZX", append(zeroPageX, c.loadALU, c.asl_alu, c.storeALU))
	c.instructionSet[ASL_A] = MkInstr("ASL_A", append(fetch16Bits, c.loadALU, c.asl_alu, c.storeALU))
	c.instructionSet[ASL_AX] = MkInstr("ASL_AX", append(absX, c.loadALU, c.asl_alu, c.storeALU))
	c.instructionSet[LSR_ACC] = MkInstr("LSR_ACC", []func(){c.lsr_acc})
	c.instructionSet[LSR_Z] = MkInstr("LSR_Z", append(fetch8Bits, c.loadALU, c.lsr_alu, c.storeALU))
	c.instructionSet[LSR_ZX] = MkInstr("LSR_ZX", append(zeroPageX, c.loadALU, c.lsr_alu, c.storeALU))
	c.instructionSet[LSR_A] = MkInstr("LSR_A", append(fetch16Bits, c.loadALU, c.lsr_alu, c.storeALU))
	c.instructionSet[LSR_AX] = MkInstr("LSR_AX", append(absX, c.loadALU, c.lsr_alu, c.storeALU))
	c.instructionSet[ROL_ACC] = MkInstr("ROL_ACC", []func(){c.rol_acc})
	c.instructionSet[ROL_Z] = MkInstr("ROL_Z", append(fetch8Bits, c.loadALU, c.rol_alu, c.storeALU))
	c.instructionSet[ROL_ZX] = MkInstr("ROL_ZX", append(zeroPageX, c.loadALU, c.rol_alu, c.storeALU))
	c.instructionSet[ROL_A] = MkInstr("ROL_A", append(fetch16Bits, c.loadALU, c.rol_alu, c.storeALU))
	c.instructionSet[ROL_AX] = MkInstr("ROL_AX", append(absX, c.loadALU, c.rol_alu, c.storeALU))
	c.instructionSet[ROR_ACC] = MkInstr("ROR_ACC", []func(){c.ror_acc})
	c.instructionSet[ROR_Z] = MkInstr("ROR_Z", append(fetch8Bits, c.loadALU, c.ror_alu, c.storeALU))
	c.instructionSet[ROR_ZX] = MkInstr("ROR_ZX", append(zeroPageX, c.loadALU, c.ror_alu, c.storeALU))
	c.instructionSet[ROR_A] = MkInstr("ROR_A", append(fetch16Bits, c.loadALU, c.ror_alu, c.storeALU))
	c.instructionSet[ROR_AX] = MkInstr("ROR_AX", append(absX, c.loadALU, c.ror_alu, c.storeALU))

	// Comparisons
	c.instructionSet[CMP_A] = MkInstr("CMP_A", append(fetch16Bits, c.cmp))
	c.instructionSet[CMP_I] = MkInstr("CMP_I", []func(){c.cmp_i})
	c.instructionSet[CMP_ZX] = MkInstr("CMP_ZX", append(zeroPageX, c.cmp))
	c.instructionSet[CMP_AX] = MkInstr("CMP_AX", append(absXOverlap, c.cmp))
	c.instructionSet[CMP_AY] = MkInstr("CMP_AY", append(absYOverlap, c.cmp))
	c.instructionSet[CMP_INDX] = MkInstr("CMP_INDX", append(indirectX, c.cmp))
	c.instructionSet[CMP_INDY] = MkInstr("CMP_INDY", append(indirectY, c.cmp))
	c.instructionSet[CMP_Z] = MkInstr("CMP_Z", append(fetch8Bits, c.cmp))
	c.instructionSet[CPX_I] = MkInstr("CPX_I", []func(){c.cpx_i})
	c.instructionSet[CPX_Z] = MkInstr("CPX_Z", append(fetch8Bits, c.cpx))
	c.instructionSet[CPX_A] = MkInstr("CPX_A", append(fetch16Bits, c.cpx))
	c.instructionSet[CPY_I] = MkInstr("CPY_I", []func(){c.cpy_i})
	c.instructionSet[CPY_Z] = MkInstr("CPY_Z", append(fetch8Bits, c.cpy))
	c.instructionSet[CPY_A] = MkInstr("CPY_A", append(fetch16Bits, c.cpy))

	interruptTail := []func(){
		c.pushInterruptReturnAddressHigh,
		c.pushInterruptReturnAddressLow,
		c.php_hw,
		c.nop,
		c.jump,
	}

	// Interrupt pseudo instructions
	c.irqPI = MkInstr("[IRQ]", append([]func(){
		func() {
			c.operand = uint16(c.bus.ReadByte(IRQ_VEC))
		},
		func() {
			c.operand |= uint16(c.bus.ReadByte(IRQ_VEC+1)) << 8
		}}, interruptTail...))
	c.rstPI = MkInstr("[RST]", append([]func(){
		func() {
			c.operand = uint16(c.bus.ReadByte(RST_VEC))
		},
		func() {
			c.operand |= uint16(c.bus.ReadByte(RST_VEC+1)) << 8
		}}, c.jump))
	c.nmiPI = MkInstr("[NMI]", append([]func(){
		func() {
			c.operand = uint16(c.bus.ReadByte(NMI_VEC))
		},
		func() {
			c.operand |= uint16(c.bus.ReadByte(NMI_VEC+1)) << 8
		}}, interruptTail...))
}

func (c *CPU) Reset() {
	c.flags = 0
	c.halted = false
	c.sp = 0xfd
	c.instruction = &c.rstPI // Load RST pseudo instruction
	c.microPc = 0
}

func (c *CPU) IRQ() {
	if c.flags&FLAG_I == 0 {
		c.irqPending = true
	}
}

func (c *CPU) NMI() {
	c.nmiPending = true
}

func (c *CPU) Clock() {
	if c.instruction == nil || c.microPc >= len(c.instruction.Microcode) {
		if c.nmiPending {
			c.instruction = &c.nmiPI
			c.nmiPending = false
			c.irqPending = false // Interrupt hijacking: NMI during while waiting for IRQ to kick in cancels the IRQ
		} else if c.irqPending {
			c.instruction = &c.irqPI
			c.irqPending = false
		} else {
			c.fetchOpcode()
		}
		c.microPc = 0
	} else {
		c.instruction.Microcode[c.microPc]()
		c.microPc++
	}
	if c.Trace {
		fmt.Println(c.StateAsString())
	}
}

func (c *CPU) StateAsString() string {
	code := ""
	if c.microPc == 0 {
		code = c.instruction.Dissasemble(c.bus, c.pc)
	}
	return fmt.Sprintf("PC=%04x [PC]=%02x MPC=%02x SP=%04x A=%02x X=%02x Y=%02x Flags=%02x Oper=%04x, [Oper]=%02x ALU=%02x Addr=%02x %s",
		c.pc, c.bus.ReadByte(c.pc), c.microPc, c.sp, c.a, c.x, c.y, c.flags, c.operand, c.bus.ReadByte(c.operand), c.alu, c.address, code)
}

func (c *CPU) fetchOpcode() {
	opcode := c.bus.ReadByte(c.pc)
	if c.CrashOnInvalidInst && len(c.instructionSet[opcode].Microcode) == 0 {
		log.Fatalf("Unknown opcode: %2x at address %4x", opcode, c.pc)
	}
	c.instruction = &c.instructionSet[opcode]
	c.pc++
}

func (c *CPU) fetchLow(target *uint16) {
	*target = uint16(c.bus.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) fetchHigh(target *uint16) {
	*target |= uint16(c.bus.ReadByte(c.pc)) << 8
	c.pc++
}

func (c *CPU) fetchHighAndJump() {
	c.fetchOperandHigh()
	c.pc = c.operand
}

func (c *CPU) jump() {
	c.pc = c.operand
}

func (c *CPU) fetchOperandLow() {
	c.fetchLow(&c.operand)
}

func (c *CPU) fetchOperandHigh() {
	c.fetchHigh(&c.operand)
}

func (c *CPU) fetchAddressLow() {
	c.address = c.bus.ReadByte(c.pc)
	c.pc++
}

func (c *CPU) fetchAddressHigh() {
	c.address = c.bus.ReadByte(c.pc)
	c.pc++
}

func (c *CPU) loadPCLow() {
	c.pc = uint16(c.bus.ReadByte(c.operand))
}

func (c *CPU) loadPCHigh() {
	// Implements 6502 indirect jump bug. MSB is not incremented when crossing page boundaries
	if c.operand&0x00ff == 0x00ff {
		c.pc |= uint16(c.bus.ReadByte(c.operand&0xff00)) << 8
	} else {
		c.pc |= uint16(c.bus.ReadByte(c.operand+1)) << 8
	}
}

func (c *CPU) fetchOperandHighAndAdd(reg *uint8) {
	c.operand |= uint16(c.bus.ReadByte(c.pc)) << 8
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
	*reg = c.bus.ReadByte(c.operand)
	c.updateNZ(*reg)
}

func (c *CPU) loadRegisterImmed(reg *uint8) {
	*reg = c.bus.ReadByte(c.pc)
	c.pc++
	c.updateNZ(*reg)
}

func (c *CPU) storeRegister(reg *uint8) {
	c.bus.WriteByte(c.operand, *reg)
}

func (c *CPU) fetchIndirectLow() {
	c.operand = uint16(c.bus.ReadByte(uint16(c.address)))
}

func (c *CPU) fetchIndirectHigh() {
	c.operand |= uint16(c.bus.ReadByte(uint16(c.address+1))) << 8
}

func (c *CPU) fetchIndirectHighAndAddY() {
	c.operand |= uint16(c.bus.ReadByte(uint16(c.address+1))) << 8
	t := c.operand + uint16(c.y)
	if t&0xf0 == c.operand&0xf0 {
		c.microPc++ // Skip extra clock cycle
	}
	c.operand = t
}

func (c *CPU) branchIf(mask, wanted uint8) {
	if c.flags&mask != wanted {
		c.microPc += 2 // Skip jump and optional nop
	}
	c.operand = uint16(c.bus.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) doBranch() {
	oldPc := c.pc
	if c.operand >= 0x80 {
		c.pc -= uint16(^uint8(c.operand) + 1) // 2s complement
	} else {
		c.pc += c.operand
	}
	if oldPc&0xff00 == c.pc&0xff00 {
		// Same page. Skip extra cycle
		c.microPc++
	}
}

func (c *CPU) push(v uint8) {
	c.bus.WriteByte(uint16(c.sp)+0x0100, v)
	c.sp--
}

func (c *CPU) pull() uint8 {
	c.sp++
	return c.bus.ReadByte(uint16(c.sp) + 0x0100)
}

func (c *CPU) pushReturnAddressLow() {
	c.push(uint8((c.pc - 1) & 0xff))
}

func (c *CPU) pushReturnAddressHigh() {
	c.push(uint8((c.pc - 1) >> 8))
}

func (c *CPU) pushInterruptReturnAddressLow() {
	c.push(uint8(c.pc & 0xff))
}

func (c *CPU) pushInterruptReturnAddressHigh() {
	c.push(uint8(c.pc >> 8))
}

func (c *CPU) incPc() {
	c.pc++
}

func (c *CPU) pullOperandHigh() {
	c.operand |= uint16(c.pull()) << 8
}

func (c *CPU) pullOperandLow() {
	c.operand = uint16(c.pull())
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
}

func (c *CPU) iny() {
	c.y++
	c.updateNZ(c.y)
}

func (c *CPU) dex() {
	c.x--
	c.updateNZ(c.x)
}

func (c *CPU) dey() {
	c.y--
	c.updateNZ(c.y)
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
}

func (c *CPU) clv() {
	c.flags &= ^FLAG_V
}

func (c *CPU) cld() {
	c.flags &= ^FLAG_D
}

func (c *CPU) cli() {
	c.flags &= ^FLAG_I
}

func (c *CPU) sec() {
	c.flags |= FLAG_C
}

func (c *CPU) sed() {
	c.flags |= FLAG_D
}

func (c *CPU) sei() {
	c.flags |= FLAG_I
}

func (c *CPU) tax() {
	c.x = c.a
	c.updateNZ(c.x)
}

func (c *CPU) tay() {
	c.y = c.a
	c.updateNZ(c.y)
}

func (c *CPU) tsx() {
	c.x = c.sp
	c.updateNZ(c.x)
}

func (c *CPU) txa() {
	c.a = c.x
	c.updateNZ(c.a)
	// c.fetchNext = true
}

func (c *CPU) txs() {
	c.sp = c.x
}

func (c *CPU) tya() {
	c.a = c.y
	c.updateNZ(c.a)
}

func (c *CPU) pha() {
	c.push(c.a)
}

func (c *CPU) php() {
	c.push(c.flags | FLAG_B | FLAG_U)
}

// Hardware interrupts must clear the B flag
func (c *CPU) php_hw() {
	c.push((c.flags & ^(FLAG_B)) | FLAG_U)
	c.flags |= FLAG_I
}

func (c *CPU) php_brk() {
	c.push(c.flags | FLAG_B | FLAG_U)
	c.flags |= FLAG_I
}

func (c *CPU) pla() {
	c.a = c.pull()
	c.updateNZ(c.a)
}

func (c *CPU) plp() {
	c.flags = c.pull() & ^(FLAG_B | FLAG_U)
}

func (c *CPU) nop() {
}

func (c *CPU) updateNZ(b uint8) {
	c.updateFlag(FLAG_Z, b == 0)
	c.updateFlag(FLAG_N, b&0x80 != 0)
}

func (c *CPU) updateFlag(flag uint8, value bool) {
	if value {
		c.flags |= flag
	} else {
		c.flags &= ^flag
	}
}

func (c *CPU) adc_i() {
	t := c.bus.ReadByte(c.pc)
	c.pc++
	c.add(t)
}

func (c *CPU) adc() {
	c.add(c.bus.ReadByte(c.operand))
}

func (c *CPU) and_i() {
	c.a &= c.bus.ReadByte(c.pc)
	c.pc++
	c.updateNZ(c.a)
}

func (c *CPU) and() {
	c.a &= c.bus.ReadByte(c.operand)
	c.updateNZ(c.a)
}

func (c *CPU) bit() {
	v := c.bus.ReadByte(c.operand)
	c.updateFlag(FLAG_Z, v&c.a == 0)
	c.flags = (c.flags & ^uint8(0xc0)) | v&0xc0
}

func (c *CPU) ora_i() {
	c.a |= c.bus.ReadByte(c.pc)
	c.pc++
	c.updateNZ(c.a)
}

func (c *CPU) ora() {
	c.a |= c.bus.ReadByte(c.operand)
	c.updateNZ(c.a)
}

func (c *CPU) eor_i() {
	c.a ^= c.bus.ReadByte(c.pc)
	c.pc++
	c.updateNZ(c.a)
}

func (c *CPU) eor() {
	c.a ^= c.bus.ReadByte(c.operand)
	c.updateNZ(c.a)
}

func (c *CPU) asl(v *uint8) {
	c.updateFlag(FLAG_C, *v&0x80 != 0)
	*v <<= 1
	c.updateNZ(*v)
}

func (c *CPU) asl_acc() {
	c.asl(&c.a)
}

func (c *CPU) asl_alu() {
	c.asl(&c.alu)
}

func (c *CPU) rol(v *uint8) {
	carry := c.flags & FLAG_C
	c.updateFlag(FLAG_C, *v&0x80 != 0)
	*v <<= 1
	*v |= carry
	c.updateNZ(*v)
}

func (c *CPU) rol_acc() {
	c.rol(&c.a)
}

func (c *CPU) rol_alu() {
	c.rol(&c.alu)
}

func (c *CPU) ror(v *uint8) {
	carry := (c.flags & FLAG_C) << 7
	c.updateFlag(FLAG_C, *v&0x01 != 0)
	*v >>= 1
	*v |= carry
	c.updateNZ(*v)
}

func (c *CPU) ror_acc() {
	c.ror(&c.a)
}

func (c *CPU) ror_alu() {
	c.ror(&c.alu)
}

func (c *CPU) lsr(v *uint8) {
	c.updateFlag(FLAG_C, *v&0x01 != 0)
	*v >>= 1
	c.updateNZ(*v)
}

func (c *CPU) lsr_acc() {
	c.lsr(&c.a)
}

func (c *CPU) lsr_alu() {
	c.lsr(&c.alu)
}

func (c *CPU) compareTwo(a, b uint8) {
	v := uint16(a) - uint16(b)
	c.updateNZ(uint8(v & 0xff))
	c.updateFlag(FLAG_C, v < 0x100)
}

func (c *CPU) cmp_i() {
	c.compareTwo(c.a, c.bus.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) cpx_i() {
	c.compareTwo(c.x, c.bus.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) cpy_i() {
	c.compareTwo(c.y, c.bus.ReadByte(c.pc))
	c.pc++
}

func (c *CPU) cmp() {
	c.compareTwo(c.a, c.bus.ReadByte(c.operand))
}

func (c *CPU) cpx() {
	c.compareTwo(c.x, c.bus.ReadByte(c.operand))
}

func (c *CPU) cpy() {
	c.compareTwo(c.y, c.bus.ReadByte(c.operand))
}

func (c *CPU) add(addend uint8) {
	acc := uint16(c.a)
	add := uint16(addend)
	carryIn := c.flags & FLAG_C
	var v uint16

	if c.flags&FLAG_D != 0 {
		lo := acc&0x0f + add&0x0f + uint16(carryIn)
		var carrylo uint16
		if lo >= 0x0a {
			carrylo = 0x10
			lo -= 0x0a
		}
		hi := (acc & 0xf0) + (add & 0xf0) + carrylo
		if hi >= 0xa0 {
			c.flags |= FLAG_C
			hi -= 0xa0
		} else {
			c.flags &= ^FLAG_C
		}
		v = hi | lo
		c.updateFlag(FLAG_V, ((acc^v)&0x80) != 0 && ((acc^add)&0x80) == 0)
	} else {
		v = acc + add + uint16(carryIn)
		c.updateFlag(FLAG_C, v >= 0x100)
		c.updateFlag(FLAG_V, ((acc&0x80) == (add&0x80)) && ((acc&0x80) != (v&0x80)))
	}

	c.a = uint8(v)
	c.updateNZ(c.a)
}

func (c *CPU) sbc() {
	c.subtract(c.bus.ReadByte(c.operand))
}

func (c *CPU) sbc_i() {
	t := c.bus.ReadByte(c.pc)
	c.pc++
	c.subtract(t)
}

func (c *CPU) subtract(addend uint8) {
	acc := uint16(c.a)
	sub := uint16(addend)
	carryIn := c.flags & FLAG_C
	var v uint16

	if c.flags&FLAG_D != 0 {
		v = (acc & 0x0f) - (sub & 0x0f) - (1 - uint16(carryIn))
		if v&0x10 != 0 {
			v = (v - 0x06) & 0x0f | ((acc & 0xf0) - (sub & 0xf0) - 0x10)
		} else {
			v = (v & 0x0f) | ((acc & 0xf0) - (sub & 0xf0))
		}
		if v&0x100 != 0 {
			v -= 0x60
		}
	} else {
		v = acc - sub - (1 - uint16(carryIn))
	}
	c.updateFlag(FLAG_C, v < 0x100)
	c.updateFlag(FLAG_V, ((acc&0x80) != (sub&0x80)) && ((acc&0x80) != (v&0x80)))
	c.a = uint8(v)
	c.updateNZ(c.a)
}

func (c *CPU) addXToLowOperand() {
	c.operand = uint16(uint8(c.operand&0xff) + c.x)
}

func (c *CPU) addYToLowOperand() {
	c.operand = uint16(uint8(c.operand&0xff) + c.y)
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
	c.alu = c.bus.ReadByte(c.operand)
}

func (c *CPU) storeALU() {
	c.bus.WriteByte(c.operand, c.alu)
}

func (c *CPU) brk() {
	c.halted = true
}

func (c *CPU) IsHalted() bool {
	return c.halted
}
