package chip8

// OneJumpTo is the 1NNN opcode which jump to the NNN address
func OneJumpTo(m *Memory, opcode uint16) {
	m.PC = opcode & 0xFFFF
}

// ThreeEqSkip is the 3XNN opcode
// which skip the next instruction if VX equals NN
func ThreeEqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>16]
	if vx == byte(opcode&0x00FF) {
		m.PC += 2
	}
}

// FourNeqSkip is the 4XNN opcode
// which skip the next instruction if VX not equals NN
func FourNeqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>16]
	if vx != byte(opcode&0x00FF) {
		m.PC += 2
	}
}

// FiveEqSkip is the 5XY0 opcode
// which skip the next instruction if VX not equals NN
func FiveEqSkip(m *Memory, opcode uint16) {
	vx := m.V[(opcode&0x0F00)>>16]
	vy := m.V[(opcode&0x00F0)>>8]
	if vx == vy {
		m.PC += 2
	}
}

// SixSetRegister is the 6XNN opcode
// which set VX to NN
func SixSetRegister(m *Memory, opcode uint16) {
	m.V[(opcode&0x0F00)>>16] = byte(opcode & 0x00FF)
}

// SevenAddToRegister is the 7XNN opcode
// which add NN to VX
func SevenAddToRegister(m *Memory, opcode uint16) {
	sum := uint16(m.V[(opcode&0x0F00)>>16]) + opcode&0x00FF
	if sum > 0x00FF {
		m.VF = true
		//TODO is this the standard behaviour ?
		m.V[(opcode & 0x0F00)] = 0
	} else {
		m.V[(opcode&0x0F00)>>16] = byte(sum)
	}
}

// ASetAddressRegister is the ANNN opcode
// which set the Address register I to NNN
func ASetAddressRegister(m *Memory, opcode uint16) {
	m.I = opcode & 0x0FFF
}

// BJumpToV0 is the BNNN
// which jump to the address V0 + NNN
func BJumpToV0(m *Memory, opcode uint16) {
	m.PC = uint16(m.V[0]) + (opcode & 0x0FFF)
}
