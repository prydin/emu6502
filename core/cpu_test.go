package core

import (
	"fmt"
	"github.com/beevik/go6502/asm"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func assemble(program string) ([]byte, error) {
	assy, _, err := asm.Assemble(strings.NewReader(program), "test.assm", os.Stderr, 0)
	for _, e := range assy.Errors {
		println(e)
	}
	return assy.Code, err
}

func runProgram(source string) []byte {
	bytes := make([]byte, 65536)
	bytes[RST_VEC] = 0
	bytes[RST_VEC+1] = 0x10

	program, err := assemble(source)
	if err != nil {
		panic(err)
	}
	copy(bytes[0x1000:], program)
	mem := LinearMemory{bytes: bytes}
	cpu := CPU{}
	cpu.Trace = true
	cpu.CrashOnInvalidInst = true
	cpu.Init(&mem)
	cpu.Reset()
	for !cpu.IsHalted() {
		cpu.Clock()
	}
	return bytes
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

	require.Equal(t, uint8(0x42), memory[0x2000], "STA $2000 failed")
	require.Equal(t, uint8(0x42), memory[0x2001], "STA $2001 failed")
	require.Equal(t, uint8(0x42), memory[0x0010], "STA $10 failed")
	require.Equal(t, uint8(0x42), memory[0x0011], "STA $11 failed")
	require.Equal(t, uint8(0x42), memory[0x0043], "STA $01,x failed")
	require.Equal(t, uint8(0x42), memory[0x0041], "STA $0ff,x failed")
	require.Equal(t, uint8(0x42), memory[0x4242], "STA $4200,x failed")
	require.Equal(t, uint8(0x42), memory[0x4242], "STA $4201,x failed")
	require.Equal(t, uint8(0x42), memory[0x4342], "STA $4300,x failed")
	require.Equal(t, uint8(0x42), memory[0x4342], "STA $4301,x failed")
	require.Equal(t, uint8(0x42), memory[0x4343], "STA $4343,x failed")
	require.Equal(t, uint8(0x42), memory[0x4446], "STA $(11),Y failed")

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

	require.Equal(t, uint8(0x42), memory[0x2000], "STX $2000 failed")
	require.Equal(t, uint8(0x42), memory[0x2001], "STX $2001 failed")
	require.Equal(t, uint8(0x42), memory[0x0010], "STX $10 failed")
	require.Equal(t, uint8(0x42), memory[0x0011], "STX $11 failed")
	require.Equal(t, uint8(0x42), memory[0x0043], "STX $01,y failed")
	require.Equal(t, uint8(0x42), memory[0x0041], "STX $ff,y failed")
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

	require.Equal(t, uint8(0x42), memory[0x2000], "STY $2000 failed")
	require.Equal(t, uint8(0x42), memory[0x2001], "STY $2001 failed")
	require.Equal(t, uint8(0x42), memory[0x0010], "STY $10 failed")
	require.Equal(t, uint8(0x42), memory[0x0011], "STY $11 failed")
	require.Equal(t, uint8(0x42), memory[0x0043], "STY $01,X failed")
	require.Equal(t, uint8(0x42), memory[0x0041], "STY $ff,X failed")
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
	require.Equal(t, uint8(1), memory[0x0000])
	require.Equal(t, uint8(2), memory[0x0001])
	require.Equal(t, uint8(0), memory[0x0002])
	require.Equal(t, uint8(0xff), memory[0x0003])
	require.Equal(t, uint8(0), memory[0x0004])
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
	require.Equal(t, uint8(1), memory[0x0000])
	require.Equal(t, uint8(2), memory[0x0001])
	require.Equal(t, uint8(0), memory[0x0002])
	require.Equal(t, uint8(0xff), memory[0x0003])
	require.Equal(t, uint8(0), memory[0x0004])
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
	require.Equal(t, uint8(0x43), memory[0x0010], "INC Zero page failed")
	require.Equal(t, uint8(0x43), memory[0x2010], "INC Absolute failed")
	require.Equal(t, uint8(0x43), memory[0x0011], "INC Zero page indexed failed")
	require.Equal(t, uint8(0x43), memory[0x2111], "INC absolute indexed failed")
	require.Equal(t, uint8(0), memory[0x0012], "INC Wraparound failed")

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
	require.Equal(t, uint8(0x41), memory[0x0010], "DEC Zero page failed")
	require.Equal(t, uint8(0x41), memory[0x2010], "DEC Absolute failed")
	require.Equal(t, uint8(0x41), memory[0x0011], "DEC Zero page indexed failed")
	require.Equal(t, uint8(0x41), memory[0x2111], "DEC absolute indexed failed")
	require.Equal(t, uint8(0xff), memory[0x0012], "DEC Wraparound failed")
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
	require.Equal(t, uint8(0x42), memory[0x0000], "Absolute JMP failed")
	require.Equal(t, uint8(0x42), memory[0x0001], "Indirect JMP failed")
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
	require.Equal(t, uint8(0x42), memory[0x0000], "Branch 1 failed")
	require.Equal(t, uint8(0x42), memory[0x0001], "Branch 2 failed")
	require.Equal(t, uint8(0x42), memory[0x0002], "Branch 3 failed")
	require.Equal(t, uint8(0x42), memory[0x0003], "Branch 4 failed")
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
	require.Equal(t, uint8(0x42), memory[0x0000], "Branch 1 failed")
	require.Equal(t, uint8(0x42), memory[0x0001], "Branch 2 failed")
	require.Equal(t, uint8(0x42), memory[0x0002], "Branch 3 failed")
	require.Equal(t, uint8(0x42), memory[0x0003], "Branch 4 failed")
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
		require.Equal(t, uint8(i), memory[0x2000+i], fmt.Sprintf("Failed at index $%02x", i))
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
		require.Equal(t, uint8(i), memory[0x2000+i], fmt.Sprintf("Failed at index $%02x", i))
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

	require.Equal(t, uint8(0xfe), memory[0x0000], "PHA failed")
	require.Equal(t, uint8(0x42), memory[0x0001], "PLA failed")
	require.Equal(t, uint8(0xff), memory[0x0002], "PLA sp failed")
	require.Equal(t, uint8(0x03), memory[0x0003], "PHP failed")
	require.Equal(t, uint8(0xfe), memory[0x0004], "PLP failed")
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
	require.Equal(t, uint8(0x42), memory[0x0000], "JSR failed")
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
		CLC
		LDA #$00
		ADC DATA ; Absolute
		STA $06 ; $01
		LDX #$01
		ADC DATA,X ; Indexed X
		STA $07 ; $03
		LDA	#DATA & $0f
		STA $20
		LDA #DATA >> 8
		STA $21
		CLC
		LDA #$03
		ADC ($20,X) ; Indirect X
		STA $08 ; $04
		LDY #$01
		LDA	#(DATA-1) & $0f
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
	require.Equal(t, uint8(0x22), memory[0x0000], "$11+$11 failed")
	require.Equal(t, uint8(0x14), memory[0x0001], "$0a+$0a failed")
	require.Equal(t, uint8(0xfe), memory[0x0002], "$ff+$ff failed")
	require.Equal(t, uint8(0x00), memory[0x0003], "$01+$ff failed")
	require.Equal(t, uint8(0x22), memory[0x0004], "11+11 failed")
	require.Equal(t, uint8(0x16), memory[0x0005], "8+8 failed")
	require.Equal(t, uint8(0x01), memory[0x0006], "Absolute failed")
	require.Equal(t, uint8(0x03), memory[0x0007], "Indexed X failed")
	require.Equal(t, uint8(0x03), memory[0x0008], "Indirect X failed")
	require.Equal(t, uint8(0x11), memory[0x0009], "Indirect Y failed")
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
	require.Equal(t, 6765, int(memory[0x0002]) + int(memory[0x0003]) << 8)
}