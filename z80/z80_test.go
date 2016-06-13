package z80

import (
	"testing"
)

func TestNew(t *testing.T) {
	z := New()
	if z.A != 0 {
		t.Errorf("A = %02x", z.A)
	}
	if z.F != 0 {
		t.Errorf("F = %02x", z.F)
	}
	if z.B != 0 {
		t.Errorf("B = %02x", z.B)
	}
	if z.C != 0 {
		t.Errorf("C = %02x", z.C)
	}
	if z.D != 0 {
		t.Errorf("D = %02x", z.D)
	}
	if z.E != 0 {
		t.Errorf("E = %02x", z.E)
	}
	if z.H != 0 {
		t.Errorf("H = %02x", z.H)
	}
	if z.L != 0 {
		t.Errorf("L = %02x", z.L)
	}
	if z.PC != 0 {
		t.Errorf("PC = %02x", z.PC)
	}
	if z.SP != 0 {
		t.Errorf("SP = %02x", z.SP)
	}
}

func TestGetBC(t *testing.T) {
	z := New()
	z.B = 0xAA
	z.C = 0x55
	bc := z.getBC()
	if bc != 0xAA55 {
		t.Errorf("BC = %04x", bc)
	}
}

func TestSetBC(t *testing.T) {
	z := New()
	z.setBC(0xAA55)
	if z.B != 0xAA {
		t.Errorf("B = %02x", z.B)
	}
	if z.C != 0x55 {
		t.Errorf("C = %02x", z.C)
	}
}

func TestGetDE(t *testing.T) {
	z := New()
	z.D = 0xAA
	z.E = 0x55
	bc := z.getDE()
	if bc != 0xAA55 {
		t.Errorf("DE = %04x", bc)
	}
}

func TestSetDE(t *testing.T) {
	z := New()
	z.setDE(0xAA55)
	if z.D != 0xAA {
		t.Errorf("D = %02x", z.D)
	}
	if z.E != 0x55 {
		t.Errorf("E = %02x", z.E)
	}
}

func TestGetHL(t *testing.T) {
	z := New()
	z.H = 0xAA
	z.L = 0x55
	bc := z.getHL()
	if bc != 0xAA55 {
		t.Errorf("HL = %04x", bc)
	}
}

func TestSetHL(t *testing.T) {
	z := New()
	z.setHL(0xAA55)
	if z.H != 0xAA {
		t.Errorf("H = %02x", z.H)
	}
	if z.L != 0x55 {
		t.Errorf("L = %02x", z.L)
	}
}

func TestGetAF(t *testing.T) {
	z := New()
	z.A = 0xAA
	z.F = 0x55
	bc := z.getAF()
	if bc != 0xAA55 {
		t.Errorf("AF = %04x", bc)
	}
}

func TestSetAF(t *testing.T) {
	z := New()
	z.setAF(0xAA55)
	if z.A != 0xAA {
		t.Errorf("A = %02x", z.A)
	}
	if z.F != 0x55 {
		t.Errorf("F = %02x", z.F)
	}
}

func TestGetZFlag(t *testing.T) {
	z := New()
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
	z := New()
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
	z := New()
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
	z := New()
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
	z := New()
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
	z := New()
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
	z := New()
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
	z := New()
	z.setCFlag(true)
	if z.F != 0x10 {
		t.Error("Failed to set C flag.")
	}
	z.setCFlag(false)
	if z.F != 0 {
		t.Error("Failed to clear C flag.")
	}
}
