# Gommodore 64

 ![Splash Screen](/doc/pacman.gif)

This is an attempt at emulating the iconic Commodore 64 in Go. It is 
currently very much a work in progress and contributions are appreciated.

## Design

### CPU
The most notable feature of this emulator is that it attempts to be
clock cycle correct. That is, instead of emulating an instruction as 
a single unit and then delaying execution for the corresponding number of 
cycles, we emulate every clock cycle separately. While this approach makes 
the CPU emulation a bit (but not much) more complex, it makes the 
emulation of e.g. the VIC-II chip *vastly* more straightforward.

Although the 6510 wasn't a microcode-based CPU, we have implemented the 
emulation as a microcode machine as this turned out to be a lot easier
than a state machine-approach.

### Performance
Great care has been taken to make the emulation as allocation-less as 
possible. Unfortunately, due to the design of the OpenGL libraries we're
using, we are currently not able to keep it 100% allocation-less, so you 
may experience occasional dips in FPS and (very rare) dropped frames when 
garbage collection occurs. 

## What works
* CPU emulation passes Klaus' test suite
* Boots BASIC without problems and seems to run BASIC programs just fine
* Correct(?) timing of bad lines etc.
* Most of VIC-II seems to work (including sprites)

## Left to do
* Emulate the SID chip
* More flexible (and usable) keyboard mapping
* Serial ports
* Tape emulation
* Memory image loading and saving
* NTSC mode

## Known bugs
* Colors of bitmap graphics seem messed up
* 38 column mode is broken
