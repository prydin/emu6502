package core

type Bank struct {
	selector int
	devices []AddressSpace
}

type BankSwitcher struct {
	banks []Bank
}

func (b *Bank) ReadByte(addr uint16) uint8 {
	return b.devices[b.selector].ReadByte(addr)
}

func (b *Bank) WriteByte(addr uint16, data uint8) {
	b.devices[b.selector].WriteByte(addr, data)
}

func (bs* BankSwitcher) Switch(selector int) {
	for i := range bs.banks {
		bs.banks[i].selector = selector
	}
}

func (bs *BankSwitcher) GetBank(index int) AddressSpace{
	return &bs.banks[index]
}

func NewBankSwitcher(devices [][]AddressSpace) *BankSwitcher {
	bs := BankSwitcher{ banks: make([]Bank, len(devices)) }
	for i, d := range devices {
		bs.banks[i] = Bank{ devices: d}
	}
	return &bs
}




