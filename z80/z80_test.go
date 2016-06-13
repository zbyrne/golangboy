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
