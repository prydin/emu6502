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
	vbus := core.Bus{}
	c.bus.ConnectClockablePh1(&c.cpu)
	colorRam := core.MakeRAM(1024)
	c.cpu.Init(&c.bus)
	c.vic.Init(&c.bus, &vbus, colorRam, screen, dimensions)
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

	ram0 := core.MakeRAM(40960) // Main RAM
	ram1 := core.MakeRAM(8192) // High RAM
	ram2 := core.MakeRAM(4096) // Banked RAM
	ram3 := core.MakeRAM(4096) // Banked RAM
	ram4 := core.MakeRAM(8192) // Banked RAM

	io := core.NewPagedSpace([]core.AddressSpace{
		&c.vic,
	})

	// Set up the main system bus
	switcher := core.NewBankSwitcher([][]core.AddressSpace {
		{ ram2, ram2, ram2, basic, ram2, ram2, ram2, basic },
		{ ram3, &charset.CharacterROM, &charset.CharacterROM, &charset.CharacterROM, ram3, io, io, io},
		{ ram4, ram4, kernal, kernal, ram4, ram4, kernal, kernal},
	})
	switcher.Switch(7)

	c.bus.Connect(ram0, 0x0000, 0x9fff) // Main RAM
	c.bus.Connect(ram1, 0xc000, 0xcfff) // High RAM
	c.bus.Connect(switcher.GetBank(0), 0xa000, 0xbfff)
	c.bus.Connect(switcher.GetBank(1), 0xd000, 0xdfff)
	c.bus.Connect(switcher.GetBank(2), 0xe000, 0xffff)
	c.bus.Connect(colorRam, 0xd800, 0xdbff) // Color RAM

	// Set up the Vic-II bus
	vbus.Connect(ram0, 0x0000, 0x9fff)
	vbus.Connect(ram1, 0xc000, 0xcfff)
	vbus.Connect(ram2, 0xa000, 0xbfff)
	vbus.Connect(ram3, 0xd000, 0xdfff)
	vbus.Connect(ram4, 0xe000, 0xffff)

	// Overlay character ROM on top of RAM
	vbus.Connect(&charset.CharacterROM, 0x1000, 0x1fff)
	vbus.Connect(&charset.CharacterROM, 0x9000, 0x9fff)
	return nil
}
