package z80

import (
	"encoding/binary"
	"testing"
)

type mockMemory struct {
	buff []byte
}

func newMockMemory(len uint16) Memory {
	var m mockMemory
	m.buff = make([]byte, len, len)
	return &m
}

func (m *mockMemory) ReadByte(addr uint16) byte {
	return m.buff[addr]
}

func (m *mockMemory) WriteByte(addr uint16, val byte) {
	m.buff[addr] = val
}

func (m *mockMemory) ReadWord(addr uint16) uint16 {
	return binary.LittleEndian.Uint16(m.buff[addr:addr+2])
}

func (m *mockMemory) WriteWord(addr uint16, val uint16) {
	binary.LittleEndian.PutUint16(m.buff[addr:addr+2], val)
}

func TestNew(t *testing.T) {
	z := New(nil)
	if z.A != 0 {
		t.Errorf("A = %02X", z.A)
	}
	if z.F != 0 {
		t.Errorf("F = %02X", z.F)
	}
	if z.B != 0 {
		t.Errorf("B = %02X", z.B)
	}
	if z.C != 0 {
		t.Errorf("C = %02X", z.C)
	}
	if z.D != 0 {
		t.Errorf("D = %02X", z.D)
	}
	if z.E != 0 {
		t.Errorf("E = %02X", z.E)
	}
	if z.H != 0 {
		t.Errorf("H = %02X", z.H)
	}
	if z.L != 0 {
		t.Errorf("L = %02X", z.L)
	}
	if z.PC != 0 {
		t.Errorf("PC = %02X", z.PC)
	}
	if z.SP != 0 {
		t.Errorf("SP = %02X", z.SP)
	}
}

func TestGetBC(t *testing.T) {
	z := New(nil)
	z.B = 0xAA
	z.C = 0x55
	bc := z.getBC()
	if bc != 0xAA55 {
		t.Errorf("BC = %04X", bc)
	}
}

func TestSetBC(t *testing.T) {
	z := New(nil)
	z.setBC(0xAA55)
	if z.B != 0xAA {
		t.Errorf("B = %02X", z.B)
	}
	if z.C != 0x55 {
		t.Errorf("C = %02X", z.C)
	}
}

func TestGetDE(t *testing.T) {
	z := New(nil)
	z.D = 0xAA
	z.E = 0x55
	bc := z.getDE()
	if bc != 0xAA55 {
		t.Errorf("DE = %04X", bc)
	}
}

func TestSetDE(t *testing.T) {
	z := New(nil)
	z.setDE(0xAA55)
	if z.D != 0xAA {
		t.Errorf("D = %02X", z.D)
	}
	if z.E != 0x55 {
		t.Errorf("E = %02X", z.E)
	}
}

func TestGetHL(t *testing.T) {
	z := New(nil)
	z.H = 0xAA
	z.L = 0x55
	bc := z.getHL()
	if bc != 0xAA55 {
		t.Errorf("HL = %04X", bc)
	}
}

func TestSetHL(t *testing.T) {
	z := New(nil)
	z.setHL(0xAA55)
	if z.H != 0xAA {
		t.Errorf("H = %02X", z.H)
	}
	if z.L != 0x55 {
		t.Errorf("L = %02X", z.L)
	}
}

func TestGetAF(t *testing.T) {
	z := New(nil)
	z.A = 0xAA
	z.F = 0x55
	bc := z.getAF()
	if bc != 0xAA55 {
		t.Errorf("AF = %04X", bc)
	}
}

func TestSetAF(t *testing.T) {
	z := New(nil)
	z.setAF(0xAA55)
	if z.A != 0xAA {
		t.Errorf("A = %02X", z.A)
	}
	if z.F != 0x55 {
		t.Errorf("F = %02X", z.F)
	}
}

func TestGetZFlag(t *testing.T) {
	z := New(nil)
	z.F = 0x80
	if !z.getZFlag() {
		t.Error("Z flag is false.")
	}
	z.F = 0
	if z.getZFlag() {
		t.Error("Z flag is true.")
	}
}

func TestGetNFlag(t *testing.T) {
	z := New(nil)
	z.F = 0x40
	if !z.getNFlag() {
		t.Error("N flag is false.")
	}
	z.F = 0
	if z.getNFlag() {
		t.Error("N flag is true.")
	}
}

func TestGetHFlag(t *testing.T) {
	z := New(nil)
	z.F = 0x20
	if !z.getHFlag() {
		t.Error("H flag is false.")
	}
	z.F = 0
	if z.getHFlag() {
		t.Error("H flag is true.")
	}
}

func TestGetCFlag(t *testing.T) {
	z := New(nil)
	z.F = 0x10
	if !z.getCFlag() {
		t.Error("C flag is false.")
	}
	z.F = 0
	if z.getCFlag() {
		t.Error("C flag is true.")
	}
}

func TestSetZFlag(t *testing.T) {
	z := New(nil)
	z.setZFlag(true)
	if z.F != 0x80 {
		t.Error("Failed to set Z flag.")
	}
	z.setZFlag(false)
	if z.F != 0 {
		t.Error("Failed to clear Z flag.")
	}
}

func TestSetNFlag(t *testing.T) {
	z := New(nil)
	z.setNFlag(true)
	if z.F != 0x40 {
		t.Error("Failed to set N flag.")
	}
	z.setNFlag(false)
	if z.F != 0 {
		t.Error("Failed to clear N flag.")
	}
}

func TestSetHFlag(t *testing.T) {
	z := New(nil)
	z.setHFlag(true)
	if z.F != 0x20 {
		t.Error("Failed to set H flag.")
	}
	z.setHFlag(false)
	if z.F != 0 {
		t.Error("Failed to clear H flag.")
	}
}

func TestSetCFlag(t *testing.T) {
	z := New(nil)
	z.setCFlag(true)
	if z.F != 0x10 {
		t.Error("Failed to set C flag.")
	}
	z.setCFlag(false)
	if z.F != 0 {
		t.Error("Failed to clear C flag.")
	}
}
