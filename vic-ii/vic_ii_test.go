package vic_ii

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVicII_Registers(t *testing.T) {
	v := VicII{}
	v.Init()
	var mask uint8
	for addr := uint16(0); addr < 0x2e; addr++ {
		switch addr {
		case REG_LPY:
			mask = 0
		case REG_LPX:
			mask = 0
		case REG_CTRL2:
			mask = 0x1f
		case REG_MEMPTR:
			mask = 0xfe
		case REG_IRQ:
			mask = 0x0f // TODO: Check the IRQ bit separately
		case REG_IRQ_ENABLE:
			mask = 0x0f
		case REG_DATA_COLL:
			mask = 0x00 // TODO: Check when we support it
		case REG_SPRITE_COLL:
			mask = 0x00 // TODO: Check when we support it
		default:
			mask = 0xff
		}
		for data := 0; data <= 0xff; data++ {
			v.WriteByte(addr, uint8(data))
			require.Equal(t, uint8(data) & mask, v.ReadByte(addr) & mask, "Mismatch at address %04x", addr)
		}
	}
}
