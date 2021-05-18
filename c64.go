package main

import (
	"github.com/prydin/emu6502/charset"
	"github.com/prydin/emu6502/core"
	vic_ii "github.com/prydin/emu6502/vic-ii"
)

type Commodore64 struct {
	cpu core.CPU
	vic vic_ii.VicII
	bus core.Bus
	ram core.RAM
}

func (c *Commodore64) Clock() {
	c.vic.Clock()
}

func (c *Commodore64) Init(screen vic_ii.Raster, dimensions vic_ii.ScreenDimensions) error {
	// Create components
	c.cpu = core.CPU{}
	c.vic = vic_ii.VicII{}
	c.bus = core.Bus{}
	c.bus.ConnectClockablePh1(&c.cpu)
	c.cpu.Init(&c.bus)
	c.vic.Init(&c.bus, screen, dimensions)
	c.ram = core.RAM{Bytes: make([]uint8, 10000)} // TODO: Size

	// Load ROMs
	var err error
	basic, err := core.LoadROM("roms/basic.901226-01.bin")
	if err != nil {
		return err
	}
	kernal, err := core.LoadROM("roms/kernal.901227-02.bin")
	if err != nil {
		return err
	}

	// Connect components to bus
	c.bus.Connect(&c.vic, 0xd000, 0xd02f)
	c.bus.Connect(basic, 0xa000, 0xbfff)
	c.bus.Connect(kernal, 0xe000, 0xffff)
	c.bus.Connect(&charset.CharacterROM, 0xd000, 0xdfff)
	c.bus.Connect(core.MakeRAM(0x9fff), 0x0000, 0x9fff) // Main RAM (including video RAM)
	c.bus.Connect(core.MakeRAM(4096), 0xc000, 0xcfff) // Upper RAM
	c.bus.Connect(core.MakeRAM(1024), 0xd800, 0xdbff) // Color RAM
	return nil
}
