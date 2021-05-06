package charset

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	f, err := os.Open("charset.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	for i, b := range buf {
		fmt.Printf("0x%02x,", b)
		if i % 16 == 15 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}
}
