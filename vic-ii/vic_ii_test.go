package vic_ii

import (
	"github.com/prydin/emu6502/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVicII_Registers(t *testing.T) {
	v := VicII{}
	v.Init(&core.Bus{}, nil, core.MakeRAM(1024), nil, PALDimensions)
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
			if addr >= REG_BORDER {
				mask = 0x0f
			} else {
				mask = 0xff
			}
		}
		for data := 0; data <= 0xff; data++ {
			v.WriteByte(addr, uint8(data))
			for i:= uint16(0); i < 0x2e; i++ {
				if i != addr {
					v.WriteByte(i, uint8(data + 1))
				}
			}
			require.Equal(t, uint8(data) & mask, v.ReadByte(addr) & mask, "Mismatch at address %04x", addr)
		}
	}
}

func TestInterruptRegister(t *testing.T) {
	v := VicII{}
	v.Init(&core.Bus{}, nil, core.MakeRAM(1024), nil, PALDimensions)
	for i := 0; i < 16; i++ {
		v.WriteByte(REG_IRQ, uint8(i))
		if i == 0 {
			require.Equalf(t, uint8(0), v.ReadByte(REG_IRQ) & 0x80, "IRQ flag should be cleared for 0x%02x", i)
		} else {
			require.Equalf(t, uint8(0x80), v.ReadByte(REG_IRQ) & 0x80, "IRQ flag should be set for 0x%02x", i)
		}
		if i & IRQ_LP != 0 {
			require.True(t, v.irqLp, "IRQ_LP should be set")
		} else {
			require.False(t, v.irqLp, "IRQ_LP should be cleared")
		}
		if i & IRQ_SPRITE_SPRITE_COLL != 0 {
			require.True(t, v.irqSpriteSprite, "IRQ_SPRITE_SPRITE_COLL should be set")
		} else {
			require.False(t, v.irqSpriteSprite, "IRQ_SPRITE_SPRITE_COLL should be cleared")
		}
		if i & IRQ_SPRITE_BG_COLL != 0 {
			require.True(t, v.irqSpriteBg, "IRQ_SPRITE_BG_COLL should be set")
		} else {
			require.False(t, v.irqSpriteBg, "IRQ_SPRITE_BG_COLL should be cleared")
		}
		if i & IRQ_RASTER != 0 {
			require.True(t, v.irqRaster, "IRQ_RASTER should be set")
		} else {
			require.False(t, v.irqRaster, "IRQ_RASTER should be cleared")
		}
	}
}

func TestInterruptEnabledRegister(t *testing.T) {
	v := VicII{}
	v.Init(&core.Bus{}, nil, core.MakeRAM(1024), nil, PALDimensions)
	for i := 0; i < 16; i++ {
		v.WriteByte(REG_IRQ_ENABLE, uint8(i))
		if i & IRQ_LP != 0 {
			require.True(t, v.irqLpEnabled, "IRQ_LP should be set")
		} else {
			require.False(t, v.irqLpEnabled, "IRQ_LP should be cleared")
		}
		if i & IRQ_SPRITE_SPRITE_COLL != 0 {
			require.True(t, v.irqSpriteSpriteEnabled, "IRQ_SPRITE_SPRITE_COLL should be set")
		} else {
			require.False(t, v.irqSpriteSpriteEnabled, "IRQ_SPRITE_SPRITE_COLL should be cleared")
		}
		if i & IRQ_SPRITE_BG_COLL != 0 {
			require.True(t, v.irqSpriteBgEnabled, "IRQ_SPRITE_BG_COLL should be set")
		} else {
			require.False(t, v.irqSpriteBgEnabled, "IRQ_SPRITE_BG_COLL should be cleared")
		}
		if i & IRQ_RASTER != 0 {
			require.True(t, v.irqRasterEnabled, "IRQ_RASTER should be set")
		} else {
			require.False(t, v.irqRasterEnabled, "IRQ_RASTER should be cleared")
		}
	}
}
