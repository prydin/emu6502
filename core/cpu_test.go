package core

import (
	"fmt"
	"github.com/beevik/go6502/asm"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

type IRQGenerator struct {
	cpu *CPU
}

type NMIGenerator struct {
	cpu *CPU
}

func (i *IRQGenerator) WriteByte(addr uint16, data uint8) {
	i.cpu.IRQ()
}

func (i *IRQGenerator) ReadByte(addr uint16) uint8 {
	return 0
}

func (n *NMIGenerator) WriteByte(addr uint16, data uint8) {
	n.cpu.NMI()
}

func (n *NMIGenerator) ReadByte(addr uint16) uint8 {
	return 0
}

func assemble(program string) ([]byte, error) {
	assy, _, err := asm.Assemble(strings.NewReader(program), "test.assm", os.Stderr, 0)
	for _, e := range assy.Errors {
		println(e)
	}
	return assy.Code, err
}

func loadProgram(source string) (*CPU, Bus) {
	bytes := make([]byte, 0x8000)
	romBytes := make([]byte, 0x1000)
	romBytes[RST_VEC-0xf000] = 0
	romBytes[RST_VEC+1-0xf000] = 0x10
	program, err := assemble(source)
	if err != nil {
		panic(err)
	}
	copy(bytes[0x1000:], program)
	mem := RAM{ bytes: bytes }
	bus := Bus{}
	cpu := CPU{}
	bus.Connect(&mem, 0x0000, 0x7fff)
	bus.Connect(&IRQGenerator{&cpu}, 0x8000,0x80ff)
	bus.Connect(&NMIGenerator{&cpu}, 0x8100,0x81ff)
	bus.Connect(&RAM{romBytes}, 0xf000,0xffff)
	cpu.Trace = true
	cpu.HaltOnBRK = true
	cpu.CrashOnInvalidInst = true
	cpu.Init(&bus)
	cpu.Reset()
	return &cpu, bus
}

func runProgram(source string) Bus {
	cpu, bus := loadProgram(source)
	for !cpu.IsHalted() {
		cpu.Clock()
	}
	return bus
}

func TestAcc(t *testing.T) {
	memory := runProgram(`
	; Basic addressing modes
	LDA #$42
	STA $2000
	LDA $2000
	STA $2001
	STA $10
	LDA $10
	STA $11
	; X zero page indexed
	LDX #$42
	LDA #$42
	STA	$01,x
	LDA $01,x
	STA $FF,x ; Should end up in $41
	; X absolute indexed
	LDX #$42
	STA $4200,X
	LDA	$4200,x
	STA $4201,x
	; Y absolute indexed
	LDY #$42
	STA $4300,Y
	LDA	$4300,Y
	STA $4301,Y
	; X indirect
	LDA #$42
	STA $4343
	LDX #$01
	LDA #$43
	STA $40
	LDA #$42
	STA $41
	LDA ($40,X)
	; Y indirect
	LDA	#$42
	STA $4445
	LDA #$44
	STA $22
	STA $23
	LDY #$01
	LDA ($22),Y
	LDY #$02
	STA ($22),Y
	BRK
`)

	require.Equal(t, uint8(0x42), memory.ReadByte(0x2000), "STA $2000 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x2001), "STA $2001 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0010), "STA $10 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0011), "STA $11 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0043), "STA $01,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0041), "STA $0ff,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4242), "STA $4200,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4242), "STA $4201,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4342), "STA $4300,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4342), "STA $4301,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4343), "STA $4343,x failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x4446), "STA $(11),Y failed")

}

func TestIndexX(t *testing.T) {
	memory := runProgram(`
	LDX #$42
	STX $2000
	LDX $2000
	STX $2001
	STX $10
	LDX $10
	STX $11
	; X zero page indexed
	LDY #$42
	STX	$01,Y
	LDX $01,Y
	STX $FF,Y ; Should end up in $41
	BRK
`)

	require.Equal(t, uint8(0x42), memory.ReadByte(0x2000), "STX $2000 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x2001), "STX $2001 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0010), "STX $10 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0011), "STX $11 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0043), "STX $01,y failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0041), "STX $ff,y failed")
}

func TestIndexY(t *testing.T) {
	memory := runProgram(`
	LDY #$42
	STY $2000
	LDY $2000
	STY $2001
	STY $10
	LDY $10
	STY $11
	; X zero page indexed
	LDX #$42
	STY	$01,X
	LDY $01,X
	STY $FF,X ; Should end up in $41
	BRK
`)

	require.Equal(t, uint8(0x42), memory.ReadByte(0x2000), "STY $2000 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x2001), "STY $2001 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0010), "STY $10 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0011), "STY $11 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0043), "STY $01,X failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0041), "STY $ff,X failed")
}

func TestINX_DEX(t *testing.T) {
	memory := runProgram(`
	LDX #$00
	INX
	STX $00 ; 1
	INX
	STX $01 ; 2
	DEX
	DEX
	STX $02 ; 0
	DEX 
	STX $03 ; FF
	INX
	STX $04 ; 0
	; TODO: Test flags!
`)
	require.Equal(t, uint8(1), memory.ReadByte(0x0000))
	require.Equal(t, uint8(2), memory.ReadByte(0x0001))
	require.Equal(t, uint8(0), memory.ReadByte(0x0002))
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0003))
	require.Equal(t, uint8(0), memory.ReadByte(0x0004))
}

func TestINY_DEY(t *testing.T) {
	memory := runProgram(`
	LDY #$00
	INY
	STY $00 ; 1
	INY
	STY $01 ; 2
	DEY
	DEY
	STY $02 ; 0
	DEY 
	STY $03 ; FF
	INY
	STY $04 ; 0
	; TODO: Test flags!
`)
	require.Equal(t, uint8(1), memory.ReadByte(0x0000))
	require.Equal(t, uint8(2), memory.ReadByte(0x0001))
	require.Equal(t, uint8(0), memory.ReadByte(0x0002))
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0003))
	require.Equal(t, uint8(0), memory.ReadByte(0x0004))
}

func TestINC(t *testing.T) {
	memory := runProgram(`
	; Zero page
	LDA #$42
	STA $10
	INC $10
	; Absolute
	STA $2010
	INC $2010
	; Zero page indexed
	STA $11
	LDX #$01
	INC $10,X
	; Absolute indexed
	STA $2111
	INC $2110,X
	; Wraparound
	LDA #$FF
	STA $12
	INC $12
`)
	require.Equal(t, uint8(0x43), memory.ReadByte(0x0010), "INC Zero page failed")
	require.Equal(t, uint8(0x43), memory.ReadByte(0x2010), "INC Absolute failed")
	require.Equal(t, uint8(0x43), memory.ReadByte(0x0011), "INC Zero page indexed failed")
	require.Equal(t, uint8(0x43), memory.ReadByte(0x2111), "INC absolute indexed failed")
	require.Equal(t, uint8(0), memory.ReadByte(0x0012), "INC Wraparound failed")

}

func TestDEC(t *testing.T) {
	memory := runProgram(`
	; Zero page
	LDA #$42
	STA $10
	DEC $10
	; Absolute
	STA $2010
	DEC $2010
	; Zero page indexed
	STA $11
	LDX #$01
	DEC $10,X
	; Absolute indexed
	STA $2111
	DEC $2110,X
	; Wraparound
	LDA #$00
	STA $12
	DEC $12
`)
	require.Equal(t, uint8(0x41), memory.ReadByte(0x0010), "DEC Zero page failed")
	require.Equal(t, uint8(0x41), memory.ReadByte(0x2010), "DEC Absolute failed")
	require.Equal(t, uint8(0x41), memory.ReadByte(0x0011), "DEC Zero page indexed failed")
	require.Equal(t, uint8(0x41), memory.ReadByte(0x2111), "DEC absolute indexed failed")
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0012), "DEC Wraparound failed")
}

func TestJMP(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Absolute
		LDA #$42
		STA $00
		JMP L1
		LDA #$00
		STA $00
L1		JMP (ADDR)
		BRK
L2		LDA #$42
		STA $01
		BRK
ADDR	.DW	L2
`)
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "Absolute JMP failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0001), "Indirect JMP failed")
}

func TestBranchOnCarry(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDA #$42
		CLC
		BCC L1
		BRK
L1		STA $00
		BCS DONE
		STA $01
		SEC
		BCC DONE
		STA $02
		BCS L2
		BRK
L2		STA $03
DONE	BRK
`)
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "Branch 1 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0001), "Branch 2 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0002), "Branch 3 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0003), "Branch 4 failed")
}

func TestBranchOnZero(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDA #$42
		LDX #$01 ; Clear zero flag
		BNE L1
		BRK
L1		STA $00
		BEQ DONE
		STA $01
		LDX #$00 ; Set zero flag
		BNE DONE
		STA $02
		LDX #$00 ; Clear zero flag
		BEQ L2
		BRK
L2		STA $03
DONE	BRK
`)
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "Branch 1 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0001), "Branch 2 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0002), "Branch 3 failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0003), "Branch 4 failed")
}

func TestCountDown(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDX #$10
LOOP	TXA
		STA $2000,X
		DEX
		BNE LOOP
		BRK
`)
	for i := 0; i < 0x10; i++ {
		require.Equal(t, uint8(i), memory.ReadByte(0x2000+uint16(i)), fmt.Sprintf("Failed at index $%02x", i))
	}
}

func TestCountUp(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDX #$f0
LOOP	TXA
		STA $2000,X
		INX
		BNE LOOP
		BRK
`)
	for i := 0xf0; i < 0x100; i++ {
		require.Equal(t, uint8(i), memory.ReadByte(0x2000+uint16(i)), fmt.Sprintf("Failed at index $%02x", i))
	}
}

func TestStack(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDA #$42
		LDX #$ff
		TXS
		PHA
		TSX
		STX $00
		PLA
		STA $01
		TSX
		STX $02
		LDA #$00 ; Set Zero flag
		SEC
		PHP
		LDA $01ff
		STA $03
		TSX
		STX $04
		BRK
`)

	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0000), "PHA failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0001), "PLA failed")
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0002), "PLA sp failed")
	require.Equal(t, uint8(0x03), memory.ReadByte(0x0003), "PHP failed")
	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0004), "PLP failed")
}

func TestSubroutine(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		LDX #$FF
		TXS
		LDA #$00
		JSR L1
		STA $00
		BRK
L1		LDA #$42
		RTS
`)
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "JSR failed")
}

func TestAdd(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		CLC
		; Basic additions
		LDA #$11
		ADC #$11
		BCS DONE
		STA $00 ; $22
		LDA #$0A
		ADC #$0A
		BCS DONE
		STA $01 ; $14
		LDA #$FF
		ADC #$FF
		BCC DONE
		STA $02 ; $FE
		CLC
		LDA #$01
		ADC #$FF 
		BCC DONE
		STA $03	; $00
		CLC
		; Decimal mode
		SED
		LDA #$11
		ADC #$11
		STA $04 ; $22
		LDA #$08
		ADC #$08
		STA $05 ; $16
		; Addressing modes
		CLD
		CLC
		LDA #$00
		ADC DATA ; Absolute
		STA $06 ; $01
		LDX #$01
		ADC DATA,X ; Indexed X
		STA $07 ; $03
		LDA	#DATA & $ff
		STA $20
		LDA #DATA >> 8
		STA $21
		CLC
		LDA #$03
		ADC ($1f,X) ; Indirect X
		STA $08 ; $04
		LDY #$01
		LDA	#(DATA-1) & $ff
		STA $20
		LDA #(DATA-1) >> 8
		STA $21
		LDA #$10
		CLC 
		ADC	($20),Y
		STA $09
DONE	BRK
DATA	.DB $01,$02
ADDR	.DW DATA-1
`)
	require.Equal(t, uint8(0x22), memory.ReadByte(0x0000), "$11+$11 failed")
	require.Equal(t, uint8(0x14), memory.ReadByte(0x0001), "$0a+$0a failed")
	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0002), "$ff+$ff failed")
	require.Equal(t, uint8(0x00), memory.ReadByte(0x0003), "$01+$ff failed")
	require.Equal(t, uint8(0x22), memory.ReadByte(0x0004), "11+11 failed")
	require.Equal(t, uint8(0x16), memory.ReadByte(0x0005), "8+8 failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0006), "Absolute failed")
	require.Equal(t, uint8(0x03), memory.ReadByte(0x0007), "Indexed X failed")
	require.Equal(t, uint8(0x04), memory.ReadByte(0x0008), "Indirect X failed")
	require.Equal(t, uint8(0x11), memory.ReadByte(0x0009), "Indirect Y failed")
}

func TestSubtract(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		CLC
		; Basic subtractions 
		LDA #$11
		SBC #$11
		BCS DONE
		STA $00 ; $00
		LDA #$11
		SBC #$0A
		BCS DONE
		STA $01 ; $07
		LDA #$FF
		SBC #$FF
		BCS DONE
		STA $02 ; $00
		CLC
		LDA #$01
		SBC #$FF 
		BCC DONE
		STA $03	; $02
		CLC
		; Decimal mode
		SED
		LDA #$11
		SBC #$11
		STA $04 ; $22
		LDA #$18
		SBC #$09
		STA $05 ; $09
		; Addressing modes
		CLD
		CLC
		LDA #$00
		SBC DATA ; Absolute
		STA $06 ; $01
		LDA #$10
		LDX #$01
		SBC DATA,X ; Indexed X
		STA $07 ; $03
		LDA	#DATA & $0ff
		STA $20
		LDA #DATA >> 8
		STA $21
		CLC
		LDA #$03
		SBC ($1f,X) ; Indirect X
		STA $08 ; $02
		LDY #$01
		LDA	#DONE & $ff
		STA $20
		LDA #DONE >> 8
		STA $21
		LDA #$10
		CLC 
		SBC	($20),Y
		STA $09
DONE	BRK
DATA	.DB $01,$02
ADDR	.DW DATA-1
`)
	require.Equal(t, uint8(0x00), memory.ReadByte(0x0000), "$11-$11 failed")
	require.Equal(t, uint8(0x07), memory.ReadByte(0x0001), "$11-$0a failed")
	require.Equal(t, uint8(0x00), memory.ReadByte(0x0002), "$ff-$ff failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0003), "$01-$ff failed")
	require.Equal(t, uint8(0x00), memory.ReadByte(0x0004), "11-11 failed")
	require.Equal(t, uint8(0x09), memory.ReadByte(0x0005), "18-9 failed")
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0006), "Absolute failed")
	require.Equal(t, uint8(0x0d), memory.ReadByte(0x0007), "Indexed X failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0008), "Indirect X failed")
	require.Equal(t, uint8(0x0f), memory.ReadByte(0x0009), "Indirect Y failed")
}

func TestOr(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Basic ORA
		LDA #$aa
		ORA #$55
		BEQ DONE
		STA $00 ; $ff
		LDA #$01
		ORA #$AA
		BEQ DONE
		STA $01 ; $AB
		LDA #$F0
		ORA #$0F
		BEQ DONE
		STA $02 ; $FF
		LDA #$00
		ORA #$00
		BNE DONE
		LDA #$01
		STA $03	; $01
		; Addressing modes
		LDA #$00
		ORA DATA ; Absolute
		STA $06 ; $01
		LDA #$00
		LDX #$01
		ORA DATA,X ; Indexed X
		STA $07 ; $02
		LDA	#DATA & $0ff
		STA $20
		LDA #DATA >> 8
		STA $21
		LDA #$00
		ORA ($1f,X) ; Indirect X
		STA $08 ; $02
		LDY #$01
		LDA	#DONE & $ff
		STA $20
		LDA #DONE >> 8
		STA $21
		LDA #$00
		ORA	($20),Y
		STA $09 ; $02
DONE	BRK
DATA	.DB $01,$02
ADDR	.DW DATA-1
`)
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0000), "$0a | $05 failed")
	require.Equal(t, uint8(0xab), memory.ReadByte(0x0001), "$0a & $01 failed")
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0002), "$00 & $00 failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0003), "$f0 & $0f failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0006), "Absolute failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0007), "Indexed X failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0008), "Indirect X failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0009), "Indirect Y failed")
}

func TestAnd(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Basic AND
		LDA #$11
		AND #$11
		BEQ DONE
		STA $00 ; $11
		LDA #$FF
		AND #$AA
		BEQ DONE
		STA $01 ; $AA
		LDA #$F0
		AND #$0F
		BNE DONE
		ADC #$01 ; 0 is the default memory content, so make sure we make our mark
		STA $02 ; $01
		LDA #$01
		SBC #$FF
		BEQ DONE
		STA $03	; $01
		; Addressing modes
		LDA #$FF
		AND DATA ; Absolute
		STA $06 ; $01
		LDA #$FF
		LDX #$01
		AND DATA,X ; Indexed X
		STA $07 ; $02
		LDA	#DATA & $0ff
		STA $20
		LDA #DATA >> 8
		STA $21
		LDA #$ff
		AND ($1f,X) ; Indirect X
		STA $08 ; $02
		LDY #$01
		LDA	#DONE & $ff
		STA $20
		LDA #DONE >> 8
		STA $21
		LDA #$ff
		AND	($20),Y
		STA $09 ; $02
DONE	BRK
DATA	.DB $01,$02
ADDR	.DW DATA-1
`)
	require.Equal(t, uint8(0x11), memory.ReadByte(0x0000), "$11 & $11 failed")
	require.Equal(t, uint8(0xaa), memory.ReadByte(0x0001), "$ff & aa failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0002), "$f0 & $0f failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0006), "Absolute failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0007), "Indexed X failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0008), "Indirect X failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0009), "Indirect Y failed")
}

func TestEor(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Basic EOR
		LDA #$11
		EOR #$11
		BNE DONE
		STA $00 ; $00
		LDA #$FF
		EOR #$AA
		BEQ DONE
		STA $01 ; $55
		LDA #$F0
		EOR #$0F
		BEQ DONE
		STA $02 ; $ff
		LDA #$01
		EOR #$FF
		BEQ DONE
		STA $03	; $FE
		; Addressing modes
		LDA #$FF
		EOR DATA ; Absolute
		STA $06 ; $01
		LDA #$FF
		LDX #$01
		EOR DATA,X ; Indexed X
		STA $07 ; $02
		LDA	#DATA & $0ff
		STA $20
		LDA #DATA >> 8
		STA $21
		LDA #$ff
		EOR ($1f,X) ; Indirect X
		STA $08 ; $fe
		LDY #$01
		LDA	#DONE & $ff
		STA $20
		LDA #DONE >> 8
		STA $21
		LDA #$ff
		EOR	($20),Y
		STA $09 ; $fe
DONE	BRK
DATA	.DB $01,$02
ADDR	.DW DATA-1
`)
	require.Equal(t, uint8(0x00), memory.ReadByte(0x0000), "$11 ^ $11 failed")
	require.Equal(t, uint8(0x55), memory.ReadByte(0x0001), "$ff ^ $aa failed")
	require.Equal(t, uint8(0xff), memory.ReadByte(0x0002), "$f0 ^ $0f failed")
	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0006), "Absolute failed")
	require.Equal(t, uint8(0xfd), memory.ReadByte(0x0007), "Indexed X failed")
	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0008), "Indirect X failed")
	require.Equal(t, uint8(0xfe), memory.ReadByte(0x0009), "Indirect Y failed")
}


func TestAsl(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Basic shift
		LDA #$01
		CLC
		ASL
		BCS DONE
		STA $00 ; $02
		; Carry in
		SEC
		LDA #$01
		ASL
		BCS DONE
		STA $01 ; $03
		; Carry out 
		SEC
		LDA #$80
		ASL
		BCC DONE
		STA $02 ; $01
		; Zero page
		LDA #$01
		STA $03
		CLC
		ASL $03
		; Zero page + 1
		LDA #$01
		STA $04
		LDX #$01
		CLC
		ASL $03,X
		; Absolute
		ASL DATA
		LDA DATA
		STA $05 ; $02
		; Absolute + X
		LDX #$01
		ASL DATA-1,X
		LDA DATA
		STA $06
		
DONE	BRK
DATA	.DB $01
		
`)
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0000), "1 << 1 failed")
	require.Equal(t, uint8(0x03), memory.ReadByte(0x0001), "1 << 1 with carry in failed")
	require.Equal(t, uint8(0x01), memory.ReadByte(0x0002), "$80 << 1 with carry out failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0003), "Zero page failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0004), "Zero page + X failed")
	require.Equal(t, uint8(0x02), memory.ReadByte(0x0005), "Absolute failed")
	require.Equal(t, uint8(0x04), memory.ReadByte(0x0006), "Absolute + X failed")
}

func TestFibonacci(t *testing.T) {
	memory := runProgram(`
		.ORG $1000

		; Initialize
		LDA #01
		STA $00 ; F(N-1) sum low
		STA $02 ; F(N) sum low
		LDA #00
		STA $01 ; F(N-1) sum high
		STA $03 ; F(N) sum high
		LDA #18 ; 
		STA $04 ; Number of iterations

		; 16 bit addition X,Y = F(N) + F(N-1)
LOOP	CLC
		LDA $00
		ADC $02
		TAX
		LDA $01
		ADC $03
		TAY
		
		; X and Y hold the low and high bits of the new sum.
		; Shift old F(N) to new F(N-1)
		LDA $02
		STA $00
		LDA $03
		STA $01
		
		; Store new F(N)
		STX $02
		STY $03
		
		; Keep looping?
		DEC $04
		BNE LOOP
		BRK
`)
	require.Equal(t, 6765, int(memory.ReadByte(0x0002))+int(memory.ReadByte(0x0003))<<8)
}

func TestMultiply(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; 2*3
		LDA #$02
		LDX #$03
		JSR MULT
		STA $10
		STX $11
		; 10*11
		LDA #$0a
		LDX #$0b
		JSR MULT
		STA $12
		STX $13
		; 100*2
		LDA #100
		LDX #2
		JSR MULT
		STA $14
		STX	$15
		; 2*100
		LDA #2
		LDX #100
		JSR MULT
		STA $16
		STX	$17
		;200*200
		LDA #200
		LDX #200
		JSR MULT
		STA $18
		STX	$19
		BRK

LOTERM	.EQ $00
HITERM	.EQ $01
LOSUM	.EQ $02
HISUM	.EQ $03
MULT	STX LOTERM	; Second term
		LDY #$00
		STY LOSUM	; Running sum
		STY HISUM
		STY HITERM
LOOP	CMP #$00	; Any set bits left in first term?
		BEQ END
		TAY	
		AND #$01	; Bit set in first term? Add shifted second term to running sum
		BEQ SHIFT
		CLC
		LDA LOTERM	; 16 bit addition
		ADC LOSUM
		STA LOSUM
		LDA HITERM
		ADC HISUM
		STA HISUM
SHIFT	TYA
		LSR
		CLC
		ASL LOTERM	; Shift the second term and continue
		ASL HITERM
		CLC
		BCC LOOP
END		LDA LOSUM
		LDX HISUM
		RTS
`)
	require.Equal(t, uint16(6), uint16(memory.ReadByte(0x0010)) + uint16(memory.ReadByte(0x0011)) << 8, "2*3 failed")
	require.Equal(t, uint16(110),uint16(memory.ReadByte(0x0012)) + uint16(memory.ReadByte(0x0013)) << 8, "10*11 failed")
	require.Equal(t, uint16(200),uint16(memory.ReadByte(0x0014)) + uint16(memory.ReadByte(0x0015)) << 8, "100*2 failed")
	require.Equal(t, uint16(200),uint16(memory.ReadByte(0x0016)) + uint16(memory.ReadByte(0x0017)) << 8, "2*100 failed")
	require.Equal(t, uint16(40000),uint16(memory.ReadByte(0x0018)) + uint16(memory.ReadByte(0x0019)) << 8, "200*2000 failed")
}

func TestCmp(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Basic operations
		LDX #$FF
		TXS
		LDA #00
		CMP #00
		PHP
		LDA #$AA
		CMP #$AA
		PHP
		LDA #$01
		CMP #$02
		PHP
		LDA #$02
		CMP #$01
		PHP
		LDA #$FE
		CMP #$FF
		PHP
		LDA #$FF
		CMP #$FE
		PHP

		; Addressing modes
		LDA #$AA
		STA ZP ; Zero page
		CLC
		CMP ZP
		PHP
		LDX #$01
		CMP ZP-1,X ; Zero page, X
		PHP
		CMP DATA-1,X ; Abs X
		PHP
		LDY #$01
		CMP DATA-1,Y ; Abs Y
		PHP
		LDY #$30
		STY ADDR
		LDY #$00
		STY ADDR+1
		CMP (ADDR-1,X) ; Ind X
		PHP
		DEC ADDR
		LDY #$01
		CMP (ADDR),Y  ; Ind Y
		PHP
		BRK

ZP		.EQ $30
DATA	.DB $AA
ADDR	.EQ $31

`)
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01ff), "CMP $00,$00 failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01fe), "CMP $AA,$AA failed")
	require.Equal(t, FLAG_N, memory.ReadByte(0x01fd), "CMP $01,$02 failed")
	require.Equal(t, FLAG_C, memory.ReadByte(0x01fc), "CMP $02,$01 failed")

	require.Equal(t, FLAG_N, memory.ReadByte(0x01fb), "CMP $fe,$ff failed")
	require.Equal(t, FLAG_C, memory.ReadByte(0x01fa), "CMP $ff,$fe failed")

	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f9), "CMP Zero page failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f8), "CMP Zero page X failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f7), "CMP Abs X failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f6), "CMP Abs Y failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f5), "CMP Ind X failed")
	require.Equal(t, FLAG_Z|FLAG_C, memory.ReadByte(0x01f4), "CMP Ind Y failed")

}

func TestPerformance(t *testing.T) {
	cpu, _ := loadProgram(`
		.ORG $1000
		LDY #00
		STY $20
		STY $21
OUTER	LDY #$00
		STY $20
INNER	TYA
		LDA $20
		LDX $21
		JSR MULT
		DEC $20 
		BNE INNER
		DEC $21 
		BNE OUTER
		BRK

LOTERM	.EQ $00
HITERM	.EQ $01
LOSUM	.EQ $02
HISUM	.EQ $03
MULT	STX LOTERM	; Second term
		LDY #$00
		STY LOSUM	; Running sum
		STY HISUM
		STY HITERM
LOOP	CMP #$00	; Any set bits left in first term?
		BEQ END
		TAY	
		AND #$01	; Bit set in first term? Add shifted second term to running sum
		BEQ SHIFT
		CLC
		LDA LOTERM	; 16 bit addition
		ADC LOSUM
		STA LOSUM
		LDA HITERM
		ADC HISUM
		STA HISUM
SHIFT	TYA
		LSR
		CLC
		ASL LOTERM	; Shift the second term and continue
		ASL HITERM
		CLC
		BCC LOOP
END		LDA LOSUM
		LDX HISUM
		RTS
`)
	cycles := 0
	cpu.Trace = false
	start := time.Now()
	for !cpu.IsHalted() {
		cpu.Clock()
		cycles++
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Elapsed time: %s, cycles: %d, speed: %d\n", elapsed, cycles, (cycles*1000) / int(elapsed))

	// Run at 1MHz
	start = time.Now()
	cycles = 0
	cpu.Reset()
	for !cpu.IsHalted() {
		now := time.Now()
		sleepTime := time.Duration(now.UnixNano() % 1000)
		endTime := now.Add(sleepTime)
		for time.Now().Before(endTime) {
			time.Sleep(0)
		}
		cpu.Clock()
		cycles++
	}
	elapsed = time.Now().Sub(start)
	fmt.Printf("Elapsed time: %s, cycles: %d, speed: %d\n", elapsed, cycles, (cycles*1000) / int(elapsed))
}

func TestIRQ(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
TRIGGER	.EQ $8000
		LDA #IRQ & $FF
		STA $FFFE
		LDA #IRQ >> 8
		STA $FFFF
LOOP	STA $8000 ; Triggers interrupt
		LDA $00
		CMP #$42
		BNE LOOP
		TSX
		STX $01
		; Try with interrupt disable bit set
		LDA #IRQ2 & $FF
		STA $FFFE
		LDA #IRQ2 >> 8
		STA $FFFF
		SEI
		STA $8000 ; Triggers interrupt (but won't, since I bit is set)
		BRK
IRQ		PHA
		LDA #$42
		STA $00
		PLA
		RTI
IRQ2	PHA
		LDA #$43
		STA $00
		PLA
		RTI
`)
	require.Equal(t, uint8(0x2c), memory.ReadByte(0xfffe), "IRQ setup lo failed")
	require.Equal(t, uint8(0x10), memory.ReadByte(0xffff), "IRQ setup hi failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "IRQ failed")
	require.Equal(t, uint8(0xfd), memory.ReadByte(0x0001), "Stack pointer incorrect")

}

func TestNMI(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
TRIGGER	.EQ $8000
		LDA #NMI & $FF
		STA $FFFA
		LDA #NMI >> 8
		STA $FFFB
LOOP	STA $8100 ; Triggers interrupt
		LDA $00
		CMP #$42
		BNE LOOP
		TSX
		STX $01
		BRK
NMI		PHA
		LDA #$42
		STA $00
		PLA
		RTI
`)
	require.Equal(t, uint8(0x17), memory.ReadByte(0xfffa), "NMI setup lo failed")
	require.Equal(t, uint8(0x10), memory.ReadByte(0xfffb), "NMI setup hi failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "NMI failed")
	require.Equal(t, uint8(0xfd), memory.ReadByte(0x0001), "Stack pointer incorrect")
}

func TestNMIDuringIRQ(t *testing.T) {
	memory := runProgram(`
		.ORG $1000
		; Set up vectors
		LDA #NMI & $FF
		STA $FFFA
		LDA #NMI >> 8
		STA $FFFB
		LDA #IRQ & $FF
		STA $FFFE
		LDA #IRQ >> 8
		STA $FFFF
LOOP	STA $8000 ; Triggers IRQ
		LDA $00
		CMP #$42
		BNE LOOP
		TSX
		STX $02
		BRK
NMI		PHA
		LDA #$42
		STA $00
		PLA
		RTI
IRQ		PHA
		LDA #$42
		STA $8100 ; Triggers NMI
		STA $01
		PLA
		RTI
`)

	require.Equal(t, uint8(0x42), memory.ReadByte(0x0000), "NMI failed")
	require.Equal(t, uint8(0x42), memory.ReadByte(0x0001), "IRQ failed")
	require.Equal(t, uint8(0xfd), memory.ReadByte(0x0002), "Stack pointer incorrect")
}

func TestKlaus(t *testing.T) {
	bytes := make([]byte, 65546)
	mem := RAM{ bytes: bytes }
	bus := Bus{}
	cpu := CPU{}
	bus.Connect(&mem, 0x0000, 0xffff)
	cpu.CrashOnInvalidInst = false
	cpu.Init(&bus)
	cpu.pc = 0x0400
	err := Load("../testsuite/6502_functional_test.bin", &bus, 0)
	if err != nil {
		panic(err)
	}
	cycles := 100000000000
	for !cpu.IsHalted() {
		cpu.Clock()
		cycles--
		if cycles == 0 {
			break
		}
		if cpu.pc == 0x346f {
			cpu.Trace = true
		}
		if cpu.pc == 0x3490 {
			fmt.Println("Hit BP")
		}
	}
}

