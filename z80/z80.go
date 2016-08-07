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

func (z *Z80) push(word uint16) {
	z.SP -= 2
	z.mem.WriteWord(z.SP, word)
}

func (z Z80) pop() uint16 {
	word := z.mem.ReadWord(z.SP)
	z.SP += 2
	return word
}

func (z *Z80) Dispatch() ClockTicks {
	var op byte
	op = z.mem.ReadByte(z.PC)
	z.PC++
	switch op {
	case 0x00:
		// NOP
		return 4
	case 0x01:
		// LD BC nn
		z.setBC(z.mem.ReadWord(z.PC))
		z.PC += 2
		return 12
	case 0x02:
		// LD (BC) A
		z.mem.WriteByte(z.getBC(), z.A)
		return 8
	case 0x03:
		// INC BC
		z.setBC(z.getBC() + 1)
		return 8
	case 0x04:
		// INC B
		res := z.B + 1
		z.setZFlag(res == 0)
		z.setNFlag(false)
		z.setHFlag(((z.B & 0xF) + 1) >= 0x10)
		z.B = res
		return 4
	case 0x05:
		// DEC B
		res := z.B - 1
		z.setZFlag(res == 0)
		z.setNFlag(true)
		z.setHFlag(((z.B & 0xF) + 0xF) >= 0x10)
		z.B = res
		return 4
	case 0x06:
		// LD B n
		z.B = z.mem.ReadByte(z.PC)
		z.PC++
		return 8
	case 0x07:
		// RLC A
		val := z.A << 1
		val |= z.A >> 7
		z.setCFlag(z.A & 0x80 != 0)
		z.setNFlag(false)
		z.setHFlag(false)
		z.setZFlag(val == 0)
		z.A = val
		return 4
	case 0x08:
		// LD (nn) SP
		z.mem.WriteWord(z.mem.ReadWord(z.PC), z.SP)
		z.PC += 2
		return 20
	case 0x09:
		// ADD HL BC
		hl := z.getHL()
		bc := z.getBC()
		z.setCFlag(hl > 0xFFFF - bc)
		z.setNFlag(false)
		z.setHFlag(hl&0xFFF + bc&0xFFF >= 0x1000)
		z.setHL(hl + bc)
		return 8
	case 0x0A:
		// LD A (BC)
		z.A = z.mem.ReadByte(z.getBC())
		return 8
	case 0x0B:
		// DEC BC
	case 0x0C:
		// INC C
	case 0x0D:
		// DEC C
	case 0x0E:
		// LD C n
	case 0x0F:
		// RRC A
	case 0x10:
	case 0x11:
	case 0x12:
	case 0x13:
	case 0x14:
	case 0x15:
	case 0x16:
	case 0x17:
	case 0x18:
	case 0x19:
	case 0x1A:
	case 0x1B:
	case 0x1C:
	case 0x1D:
	case 0x1E:
	case 0x1F:
	case 0x20:
	case 0x21:
	case 0x22:
	case 0x23:
	case 0x24:
	case 0x25:
	case 0x26:
	case 0x27:
	case 0x28:
	case 0x29:
	case 0x2A:
	case 0x2B:
	case 0x2C:
	case 0x2D:
	case 0x2E:
	case 0x2F:
	case 0x30:
	case 0x31:
	case 0x32:
	case 0x33:
	case 0x34:
	case 0x35:
	case 0x36:
	case 0x37:
	case 0x38:
	case 0x39:
	case 0x3A:
	case 0x3B:
	case 0x3C:
	case 0x3D:
	case 0x3E:
	case 0x3F:
	case 0x40:
	case 0x41:
	case 0x42:
	case 0x43:
	case 0x44:
	case 0x45:
	case 0x46:
	case 0x47:
	case 0x48:
	case 0x49:
	case 0x4A:
	case 0x4B:
	case 0x4C:
	case 0x4D:
	case 0x4E:
	case 0x4F:
	case 0x50:
	case 0x51:
	case 0x52:
	case 0x53:
	case 0x54:
	case 0x55:
	case 0x56:
	case 0x57:
	case 0x58:
	case 0x59:
	case 0x5A:
	case 0x5B:
	case 0x5C:
	case 0x5D:
	case 0x5E:
	case 0x5F:
	case 0x60:
	case 0x61:
	case 0x62:
	case 0x63:
	case 0x64:
	case 0x65:
	case 0x66:
	case 0x67:
	case 0x68:
	case 0x69:
	case 0x6A:
	case 0x6B:
	case 0x6C:
	case 0x6D:
	case 0x6E:
	case 0x6F:
	case 0x70:
	case 0x71:
	case 0x72:
	case 0x73:
	case 0x74:
	case 0x75:
	case 0x76:
	case 0x77:
	case 0x78:
	case 0x79:
	case 0x7A:
	case 0x7B:
	case 0x7C:
	case 0x7D:
	case 0x7E:
	case 0x7F:
	case 0x80:
	case 0x81:
	case 0x82:
	case 0x83:
	case 0x84:
	case 0x85:
	case 0x86:
	case 0x87:
	case 0x88:
	case 0x89:
	case 0x8A:
	case 0x8B:
	case 0x8C:
	case 0x8D:
	case 0x8E:
	case 0x8F:
	case 0x90:
	case 0x91:
	case 0x92:
	case 0x93:
	case 0x94:
	case 0x95:
	case 0x96:
	case 0x97:
	case 0x98:
	case 0x99:
	case 0x9A:
	case 0x9B:
	case 0x9C:
	case 0x9D:
	case 0x9E:
	case 0x9F:
	case 0xA0:
	case 0xA1:
	case 0xA2:
	case 0xA3:
	case 0xA4:
	case 0xA5:
	case 0xA6:
	case 0xA7:
	case 0xA8:
	case 0xA9:
	case 0xAA:
	case 0xAB:
	case 0xAC:
	case 0xAD:
	case 0xAE:
	case 0xAF:
	case 0xB0:
	case 0xB1:
	case 0xB2:
	case 0xB3:
	case 0xB4:
	case 0xB5:
	case 0xB6:
	case 0xB7:
	case 0xB8:
	case 0xB9:
	case 0xBA:
	case 0xBB:
	case 0xBC:
	case 0xBD:
	case 0xBE:
	case 0xBF:
	case 0xC0:
	case 0xC1:
	case 0xC2:
	case 0xC3:
	case 0xC4:
	case 0xC5:
	case 0xC6:
	case 0xC7:
	case 0xC8:
	case 0xC9:
	case 0xCA:
	case 0xCB:
	case 0xCC:
	case 0xCD:
	case 0xCE:
	case 0xCF:
	case 0xD0:
	case 0xD1:
	case 0xD2:
	case 0xD3:
	case 0xD4:
	case 0xD5:
	case 0xD6:
	case 0xD7:
	case 0xD8:
	case 0xD9:
	case 0xDA:
	case 0xDB:
	case 0xDC:
	case 0xDD:
	case 0xDE:
	case 0xDF:
	case 0xE0:
	case 0xE1:
	case 0xE2:
	case 0xE3:
	case 0xE4:
	case 0xE5:
	case 0xE6:
	case 0xE7:
	case 0xE8:
	case 0xE9:
	case 0xEA:
	case 0xEB:
	case 0xEC:
	case 0xED:
	case 0xEE:
	case 0xEF:
	case 0xF0:
	case 0xF1:
	case 0xF2:
	case 0xF3:
	case 0xF4:
	case 0xF5:
	case 0xF6:
	case 0xF7:
	case 0xF8:
	case 0xF9:
	case 0xFA:
	case 0xFB:
	case 0xFC:
	case 0xFD:
	case 0xFE:
	case 0xFF:
	}
	return 0
}
