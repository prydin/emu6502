package core

import (
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
