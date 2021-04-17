package core

import (
	"io/ioutil"
	"os"
)

func Load(filename string, memory AddressSpace, start uint16) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	code, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	for i, b := range code {
		memory.WriteByte(uint16(i) + start, b)
	}
	return nil
}
