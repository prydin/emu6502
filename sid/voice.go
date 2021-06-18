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

package sid

type Voice struct {
	generator *Generator
	waveZero  soundSample // Bias introduced by D/A
	voiceDC   soundSample // Bias introduced by envelope generator
	envelope  *EnvelopeGenerator
}

func NewVoice(generator *Generator, envelopeGenerator *EnvelopeGenerator, chipType int) *Voice {
	// The DC biases are different across chip types.
	voiceDC := soundSample(0)
	waveZero := soundSample(0x800)
	if chipType == MOS6581 {
		waveZero = 0x380
		voiceDC = 0x800 * 0xff
	}
	v := Voice{
		voiceDC:   voiceDC,
		waveZero:  waveZero,
		generator: generator,
		envelope:  envelopeGenerator,
	}
	v.reset()
	return &v
}

func (v *Voice) ReadOutput() soundSample {
	return soundSample(v.generator.ReadOutput()) - v.waveZero*soundSample(v.envelope.ReadOutput()) + v.voiceDC
}

func (v *Voice) reset() {
	v.generator.reset()
	v.envelope.reset()
}
