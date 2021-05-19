package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBankSwitcher(t *testing.T) {
	ram0_1 := MakeRAM(1024)
	ram0_2 := MakeRAM(1024)
	ram1_1 := MakeRAM(1024)
	ram1_2 := MakeRAM(1024)
	ram2_1 := MakeRAM(1024)
	ram2_2 := MakeRAM(1024)
	bs := NewBankSwitcher([][]AddressSpace{
		[]AddressSpace{ ram0_1, ram0_2 },
		[]AddressSpace{ ram1_1, ram1_2 },
		[]AddressSpace{ ram2_1, ram2_2 },
	})

	bs.Switch(0)
	bus := Bus{}
	bus.Connect(bs.GetBank(0), 0x0000, 0x00ff)
	bus.Connect(bs.GetBank(1), 0x1000, 0x10ff)
	bus.WriteByte(0x0000, 0)
	bus.WriteByte(0x1000, 1)
	bs.Switch(1)
	bus.WriteByte(0x0000, 2)
	bus.WriteByte(0x1000, 3)
	bs.Switch(0)
	require.Equal(t, uint8(0), bus.ReadByte(0x0000))
	require.Equal(t, uint8(1), bus.ReadByte(0x1000))
	bs.Switch(1)
	require.Equal(t, uint8(2), bus.ReadByte(0x0000))
	require.Equal(t, uint8(3), bus.ReadByte(0x1000))
}
