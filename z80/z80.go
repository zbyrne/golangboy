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

func addSignedByteToU16(a uint16, b byte) uint16 {
	var upper byte
	if (b & 0x80 == 0x80) {
		upper = 0xFF
	}
	word:= []byte{b, upper}
	return a + binary.LittleEndian.Uint16(word)
}

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

func (z *Z80) regDecode(op byte) *byte {
	if op < 0x08 {
		return &z.B
	}
	if op < 0x10 {
		return &z.C
	}
	if op < 0x18 {
		return &z.D
	}
	if op < 0x20 {
		return &z.E
	}
	if op < 0x28 {
		return &z.H
	}
	if op < 0x30 {
		return &z.L
	}
	if op < 0x40 {
		return &z.A
	}
	// LD and ALU instructions follow this pattern for input operands
	switch op & 0x7 {
	case 0:
		return &z.B
	case 1:
		return &z.C
	case 2:
		return &z.D
	case 3:
		return &z.E
	case 4:
		return &z.H
	case 5:
		return &z.L
	case 6:
		panic("HL indirect instruction.")
	case 7:
		return &z.A
	}
	return nil
}

func (z *Z80) r16GetSetDecode(op byte) (func() uint16, func(uint16)) {
	if op < 0x10 {
		return func() uint16 {return z.getBC()}, func(x uint16) {z.setBC(x)}
	}
	if op < 0x20 {
		return func() uint16 {return z.getDE()}, func(x uint16) {z.setDE(x)}
	}
	if op < 0x30 {
		return func() uint16 {return z.getHL()}, func(x uint16) {z.setHL(x)}
	}
	if op < 0x40 {
		return func() uint16 {return z.SP}, func(x uint16) {z.SP = x}
	}
	return nil, nil
}

func (z *Z80) ldDestRegDecode(op byte) *byte {
	switch op & 0x78 {
	case 0x40:
		return &z.B
	case 0x48:
		return &z.C
	case 0x50:
		return &z.D
	case 0x58:
		return &z.E
	case 0x60:
		return &z.H
	case 0x68:
		return &z.L
	case 0x70:
		panic("HL indirect instruction.")
	case 0x78:
		return &z.A
	}
	return nil
}

func (z *Z80) Dispatch() ClockTicks {
	var op byte
	var reg *byte
	var getReg16 func() uint16
	var setReg16 func(uint16)
	op = z.mem.ReadByte(z.PC)
	z.PC++
	switch op {
	case 0x00:
		// NOP
		return 4
	case 0x01, 0x11, 0x21, 0x31:
		// LD R16 nn
		_, setReg16 = z.r16GetSetDecode(op)
		setReg16(z.mem.ReadWord(z.PC))
		z.PC += 2
		return 12
	case 0x02, 0x12:
		// LD (R16) A
		getReg16, _ = z.r16GetSetDecode(op)
		z.mem.WriteByte(getReg16(), z.A)
		return 8
	case 0x03, 0x13, 0x23, 0x33:
		// INC R16
		getReg16, setReg16 = z.r16GetSetDecode(op)
		setReg16(getReg16() + 1)
		return 8
	case 0x04, 0x0C, 0x14, 0x1C, 0x24, 0x2C, 0x3C:
		// INC R8
		reg = z.regDecode(op)
		res := *reg + 1
		z.setZFlag(res == 0)
		z.setNFlag(false)
		z.setHFlag(((*reg & 0xF) + 1) >= 0x10)
		*reg = res
		return 4
	case 0x05, 0x0D, 0x15, 0x1D, 0x25, 0x2D, 0x3D:
		// DEC R8
		reg = z.regDecode(op)
		res := *reg - 1
		z.setZFlag(res == 0)
		z.setNFlag(true)
		z.setHFlag(((*reg & 0xF) + 0xF) >= 0x10)
		*reg = res
		return 4
	case 0x06, 0x0E, 0x16, 0x1E, 0x26, 0x2E, 0x3E:
		// LD R8 n
		reg = z.regDecode(op)
		*reg = z.mem.ReadByte(z.PC)
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
	case 0x09, 0x19, 0x29, 0x39:
		// ADD HL R16
		getReg16, _ = z.r16GetSetDecode(op)
		hl := z.getHL()
		r16 := getReg16()
		z.setCFlag(hl > 0xFFFF - r16)
		z.setNFlag(false)
		z.setHFlag(hl&0xFFF + r16&0xFFF >= 0x1000)
		z.setHL(hl + r16)
		return 8
	case 0x0A, 0x1A:
		// LD A (R16)
		getReg16, _ = z.r16GetSetDecode(op)
		z.A = z.mem.ReadByte(getReg16())
		return 8
	case 0x0B, 0x1B, 0x2B, 0x3B:
		// DEC R16
		getReg16, setReg16 = z.r16GetSetDecode(op)
		setReg16(getReg16() - 1)
		return 8
	case 0x0F:
		// RRC A
		val := z.A >> 1
		val |= (z.A & 1) << 7
		z.setCFlag(z.A & 1 != 0)
		z.setNFlag(false)
		z.setHFlag(false)
		z.setZFlag(val == 0)
		z.A = val
		return 4
	case 0x10:
		// STOP
		// TODO: fill this in once there's more system around.
		z.PC++
		return 4
	case 0x17:
		// RL A
		var carry uint8 = 0
		val := z.A << 1
		if z.getCFlag() {
			carry = 1
		}
		val |= carry
		z.setCFlag(z.A & 0x80 != 0)
		z.setNFlag(false)
		z.setHFlag(false)
		z.setZFlag(val == 0)
		z.A = val
		return 4
	case 0x18:
		// JR n
		offset := z.mem.ReadByte(z.PC)
		z.PC++
		z.PC = addSignedByteToU16(z.PC, offset)
		return 12
	case 0x1F:
		// RRA
	case 0x20:
		// JR NZ n
	case 0x22:
		// LD (HL+) A
		z.mem.WriteByte(z.getHL(), z.A)
		z.setHL(z.getHL() + 1)
		return 8
	case 0x27:
		// DAA
	case 0x28:
		// JR Z n
	case 0x2A:
		// LD A (HL+)
		z.A = z.mem.ReadByte(z.getHL())
		z.setHL(z.getHL() + 1)
		return 8
	case 0x2F:
		// CPL
	case 0x30:
		// JR NC n
	case 0x32:
		// LD (HL-) A
		z.mem.WriteByte(z.getHL(), z.A)
		z.setHL(z.getHL() - 1)
		return 8
	case 0x34:
		// INC (HL)
	case 0x35:
		// DEC (HL)
	case 0x36:
		// LD (HL) n
	case 0x37:
		// SCF
		z.setCFlag(true)
		z.setNFlag(false)
		z.setHFlag(false)
		return 4
	case 0x38:
		// JR C n
	case 0x3A:
		// LD A (HL-)
		z.A = z.mem.ReadByte(z.getHL())
		z.setHL(z.getHL() - 1)
	case 0x3F:
		// CCF
		z.setCFlag(false)
		z.setNFlag(false)
		z.setHFlag(false)
		return 4
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x47,
		0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4F,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x57,
		0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5F,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x67,
		0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6F,
		0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7F:
		// LD R8 R8
		*z.ldDestRegDecode(op) = *z.regDecode(op)
		return 4
	case 0x46, 0x4E, 0x56, 0x5E, 0x66, 0x6E, 0x7E:
		// LD R8 (HL)
		*z.ldDestRegDecode(op) = z.mem.ReadByte(z.getHL())
		return 8
	case 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x77:
		// LD (HL) R8
		z.mem.WriteByte(z.getHL(), *z.regDecode(op))
		return 8
	case 0x76:
		// HALT
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x87:
		// ADD A R8
		reg = z.regDecode(op)
		z.setZFlag(z.A + *reg == 0)
		z.setNFlag(false)
		z.setHFlag((z.A & 0xF) > 0xF - (*reg) & 0xF)
		z.setCFlag(z.A > 0xFF - *reg)
		z.A += *reg
		return 4
	case 0x86:
		// ADD A (HL)
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
