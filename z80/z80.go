package z80

import (
	"encoding/binary"
)

const (
	Z_FLAG byte = 1 << 7
	N_FLAG byte = 1 << 6
	H_FLAG byte = 1 << 5
	C_FLAG byte = 1 << 4
)

type Memory interface {
	ReadByte(uint16) byte
	WriteByte(uint16, byte)
	ReadWord(uint16) uint16
	WriteWord(uint16, uint16)
}

type Z80 struct {
	A, F, B, C, D, E, H, L byte

	PC, SP uint16
	mem Memory
}

type ClockTicks int

func New(m Memory) Z80 {
	return Z80{mem: m}
}

func (z Z80) getBC() uint16 {
	bc := []byte{z.C, z.B}
	return binary.LittleEndian.Uint16(bc)
}

func (z Z80) getDE() uint16 {
	de := []byte{z.E, z.D}
	return binary.LittleEndian.Uint16(de)
}

func (z Z80) getHL() uint16 {
	hl := []byte{z.L, z.H}
	return binary.LittleEndian.Uint16(hl)
}

func (z Z80) getAF() uint16 {
	af := []byte{z.F, z.A}
	return binary.LittleEndian.Uint16(af)
}

func (z *Z80) setBC(bc uint16) {
	bytes := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bytes, bc)
	z.C = bytes[0]
	z.B = bytes[1]
}

func (z *Z80) setDE(bc uint16) {
	bytes := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bytes, bc)
	z.E = bytes[0]
	z.D = bytes[1]
}

func (z *Z80) setHL(bc uint16) {
	bytes := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bytes, bc)
	z.L = bytes[0]
	z.H = bytes[1]
}

func (z *Z80) setAF(bc uint16) {
	bytes := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(bytes, bc)
	z.F = bytes[0]
	z.A = bytes[1]
}

func (z Z80) getZFlag() bool {
	return z.F & Z_FLAG != 0
}

func (z Z80) getNFlag() bool {
	return z.F & N_FLAG != 0
}

func (z Z80) getHFlag() bool {
	return z.F & H_FLAG != 0
}

func (z Z80) getCFlag() bool {
	return z.F & C_FLAG != 0
}

func (z *Z80) setZFlag(f bool) {
	if f {
		z.F |= Z_FLAG
	} else {
		z.F &= ^Z_FLAG
	}
}

func (z *Z80) setNFlag(f bool) {
	if f {
		z.F |= N_FLAG
	} else {
		z.F &= ^N_FLAG
	}
}

func (z *Z80) setHFlag(f bool) {
	if f {
		z.F |= H_FLAG
	} else {
		z.F &= ^H_FLAG
	}
}

func (z *Z80) setCFlag(f bool) {
	if f {
		z.F |= C_FLAG
	} else {
		z.F &= ^C_FLAG
	}
}
