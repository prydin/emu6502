/*
 * Copyright (c) 2021 Pontus Rydin
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the Software
 * without restriction, including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
 * to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package computer

import (
	"github.com/prydin/emu6502/charset"
	"github.com/prydin/emu6502/cia"
	"github.com/prydin/emu6502/core"
	"github.com/prydin/emu6502/keyboard"
	vic_ii "github.com/prydin/emu6502/vic-ii"
)

type Commodore64 struct {
	Cpu      core.CPU
	Vic      vic_ii.VicII
	Bus      core.Bus
	ram      core.RAM
	Keyboard *keyboard.Keyboard
}

func (c *Commodore64) Clock() {
	c.Vic.Clock()
}

func (c *Commodore64) Init(screen vic_ii.Raster, dimensions vic_ii.ScreenDimensions) error {
	// Create components
	c.Cpu = core.CPU{}
	c.Vic = vic_ii.VicII{}
	c.Bus = core.Bus{}
	vbus := core.Bus{}
	c.Bus.ConnectClockablePh1(&c.Cpu)
	colorRam := core.MakeRAM(1024)

	// Initialize CPI
	c.Cpu.Init(&c.Bus)

	// Initialize CIAs
	cia1 := cia.CIA{}
	cia1.Init(&c.Bus)
	c.Bus.ConnectClockablePh1(&cia1)
	cia2 := cia.CIA{}
	cia2.Init(&c.Bus)
	c.Bus.ConnectClockablePh1(&cia2)

	// Initialize VIC
	c.Vic.Init(&vbus, &c.Bus, &cia2, colorRam, screen, dimensions)
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
	ram1 := core.MakeRAM(8192)  // High RAM
	ram2 := core.MakeRAM(8192)  // Banked RAM
	ram3 := core.MakeRAM(4096)  // Banked RAM
	ram4 := core.MakeRAM(8192)  // Banked RAM

	io := core.NewPagedSpace([]core.AddressSpace{
		&c.Vic,            // D000
		&c.Vic,            // D100
		&c.Vic,            // D200
		&c.Vic,            // D300
		core.MakeRAM(256), // D400 TODO: Placeholder for SID
		core.MakeRAM(256), // D500 TODO: Placeholder for SID
		core.MakeRAM(256), // D600 TODO: Placeholder for SID
		core.MakeRAM(256), // D700 TODO: Placeholder for SID
		colorRam.Page(0),  // D800
		colorRam.Page(1),  // D900
		colorRam.Page(2),  // DA00
		colorRam.Page(3),  // DB00
		&cia1,             // DC00
		&cia2,             // DD00
	}, true)

	// Set up the main system Bus
	switcher := core.NewBankSwitcher([][]core.AddressSpace{
		{ram2, ram2, ram2, basic, ram2, ram2, ram2, basic},
		{ram3, ram3, &charset.CharacterROM, &charset.CharacterROM, ram3, io, io, io},
		{ram4, ram4, kernal, kernal, ram4, ram4, kernal, kernal},
	})
	switcher.Switch(7)
	c.Bus.SetSwitcher(switcher)

	c.Bus.Connect(ram0, 0x0000, 0x9fff) // Main RAM
	c.Bus.Connect(ram1, 0xc000, 0xcfff) // High RAM
	c.Bus.Connect(switcher.GetBank(0), 0xa000, 0xbfff)
	c.Bus.Connect(switcher.GetBank(1), 0xd000, 0xdfff)
	c.Bus.Connect(switcher.GetBank(2), 0xe000, 0xffff)
	//c.Bus.Connect(colorRam, 0xd800, 0xdbff) // Color RAM

	// Connect peripherals
	c.Keyboard = &keyboard.Keyboard{}
	c.Keyboard.Init(&cia1)
	c.Bus.ConnectClockablePh1(c.Keyboard)

	// Set up the Vic-II Bus
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
