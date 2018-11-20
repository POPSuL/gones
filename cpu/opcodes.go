package cpu

type OpCode struct {
	FullName string
	BaseName string
	Addressing uint
	Cycle uint
}

var Cycles = [...]uint{
	7, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 3, 2, 2, 2, 3, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	6, 6, 2, 8, 3, 3, 5, 5, 4, 2, 2, 2, 5, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 6, 7,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 6, 2, 6, 4, 4, 4, 4, 2, 4, 2, 5, 5, 4, 5, 5,
	2, 6, 2, 6, 3, 3, 3, 3, 2, 2, 2, 2, 4, 4, 4, 4,
	2, 5, 2, 5, 4, 4, 4, 4, 2, 4, 2, 4, 4, 4, 4, 4,
	2, 6, 2, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
	2, 6, 3, 8, 3, 3, 5, 5, 2, 2, 2, 2, 4, 4, 6, 6,
	2, 5, 2, 8, 4, 4, 6, 6, 2, 4, 2, 7, 4, 4, 7, 7,
}

var OpCodes = map[int]OpCode{
	0xA9: {"LDA_IMM", "LDA", Immediate, Cycles[0xA9]},
	0xA5: {"LDA_ZERO", "LDA", ZeroPage, Cycles[0xA5]},
	0xAD: {"LDA_ABS", "LDA", Absolute, Cycles[0xAD]},
	0xB5: {"LDA_ZEROX", "LDA", ZeroPageX, Cycles[0xB5]},
	0xBD: {"LDA_ABSX", "LDA", AbsoluteX, Cycles[0xBD]},
	0xB9: {"LDA_ABSY", "LDA", AbsoluteY, Cycles[0xB9]},
	0xA1: {"LDA_INDX", "LDA", PreIndexedIndirect, Cycles[0xA1]},
	0xB1: {"LDA_INDY", "LDA", PostIndexedIndirect, Cycles[0xB1]},
	0xA2: {"LDX_IMM", "LDX", Immediate, Cycles[0xA2]},
	0xA6: {"LDX_ZERO", "LDX", ZeroPage, Cycles[0xA6]},
	0xAE: {"LDX_ABS", "LDX", Absolute, Cycles[0xAE]},
	0xB6: {"LDX_ZEROY", "LDX", ZeroPageY, Cycles[0xB6]},
	0xBE: {"LDX_ABSY", "LDX", AbsoluteY, Cycles[0xBE]},
	0xA0: {"LDY_IMM", "LDY", Immediate, Cycles[0xA0]},
	0xA4: {"LDY_ZERO", "LDY", ZeroPage, Cycles[0xA4]},
	0xAC: {"LDY_ABS", "LDY", Absolute, Cycles[0xAC]},
	0xB4: {"LDY_ZEROX", "LDY", ZeroPageX, Cycles[0xB4]},
	0xBC: {"LDY_ABSX", "LDY", AbsoluteX, Cycles[0xBC]},
	0x85: {"STA_ZERO", "STA", ZeroPage, Cycles[0x85]},
	0x8D: {"STA_ABS", "STA", Absolute, Cycles[0x8D]},
	0x95: {"STA_ZEROX", "STA", ZeroPageX, Cycles[0x95]},
	0x9D: {"STA_ABSX", "STA", AbsoluteX, Cycles[0x9D]},
	0x99: {"STA_ABSY", "STA", AbsoluteY, Cycles[0x99]},
	0x81: {"STA_INDX", "STA", PreIndexedIndirect, Cycles[0x81]},
	0x91: {"STA_INDY", "STA", PostIndexedIndirect, Cycles[0x91]},
	0x86: {"STX_ZERO", "STX", ZeroPage, Cycles[0x86]},
	0x8E: {"STX_ABS", "STX", Absolute, Cycles[0x8E]},
	0x96: {"STX_ZEROY", "STX", ZeroPageY, Cycles[0x96]},
	0x84: {"STY_ZERO", "STY", ZeroPage, Cycles[0x84]},
	0x8C: {"STY_ABS", "STY", Absolute, Cycles[0x8C]},
	0x94: {"STY_ZEROX", "STY", ZeroPageX, Cycles[0x94]},
	0x8A: {"TXA", "TXA", Implied, Cycles[0x8A]},
	0x98: {"TYA", "TYA", Implied, Cycles[0x98]},
	0x9A: {"TXS", "TXS", Implied, Cycles[0x9A]},
	0xA8: {"TAY", "TAY", Implied, Cycles[0xA8]},
	0xAA: {"TAX", "TAX", Implied, Cycles[0xAA]},
	0xBA: {"TSX", "TSX", Implied, Cycles[0xBA]},
	0x8: {"PHP", "PHP", Implied, Cycles[0x08]},
	0x28: {"PLP", "PLP", Implied, Cycles[0x28]},
	0x48: {"PHA", "PHA", Implied, Cycles[0x48]},
	0x68: {"PLA", "PLA", Implied, Cycles[0x68]},
	0x69: {"ADC_IMM", "ADC", Immediate, Cycles[0x69]},
	0x65: {"ADC_ZERO", "ADC", ZeroPage, Cycles[0x65]},
	0x6D: {"ADC_ABS", "ADC", Absolute, Cycles[0x6D]},
	0x75: {"ADC_ZEROX", "ADC", ZeroPageX, Cycles[0x75]},
	0x7D: {"ADC_ABSX", "ADC", AbsoluteX, Cycles[0x7D]},
	0x79: {"ADC_ABSY", "ADC", AbsoluteY, Cycles[0x79]},
	0x61: {"ADC_INDX", "ADC", PreIndexedIndirect, Cycles[0x61]},
	0x71: {"ADC_INDY", "ADC", PostIndexedIndirect, Cycles[0x71]},
	0xE9: {"SBC_IMM", "SBC", Immediate, Cycles[0xE9]},
	0xE5: {"SBC_ZERO", "SBC", ZeroPage, Cycles[0xE5]},
	0xED: {"SBC_ABS", "SBC", Absolute, Cycles[0xED]},
	0xF5: {"SBC_ZEROX", "SBC", ZeroPageX, Cycles[0xF5]},
	0xFD: {"SBC_ABSX", "SBC", AbsoluteX, Cycles[0xFD]},
	0xF9: {"SBC_ABSY", "SBC", AbsoluteY, Cycles[0xF9]},
	0xE1: {"SBC_INDX", "SBC", PreIndexedIndirect, Cycles[0xE1]},
	0xF1: {"SBC_INDY", "SBC", PostIndexedIndirect, Cycles[0xF1]},
	0xE0: {"CPX_IMM", "CPX", Immediate, Cycles[0xE0]},
	0xE4: {"CPX_ZERO", "CPX", ZeroPage, Cycles[0xE4]},
	0xEC: {"CPX_ABS", "CPX", Absolute, Cycles[0xEC]},
	0xC0: {"CPY_IMM", "CPY", Immediate, Cycles[0xC0]},
	0xC4: {"CPY_ZERO", "CPY", ZeroPage, Cycles[0xC4]},
	0xCC: {"CPY_ABS", "CPY", Absolute, Cycles[0xCC]},
	0xC9: {"CMP_IMM", "CMP", Immediate, Cycles[0xC9]},
	0xC5: {"CMP_ZERO", "CMP", ZeroPage, Cycles[0xC5]},
	0xCD: {"CMP_ABS", "CMP", Absolute, Cycles[0xCD]},
	0xD5: {"CMP_ZEROX", "CMP", ZeroPageX, Cycles[0xD5]},
	0xDD: {"CMP_ABSX", "CMP", AbsoluteX, Cycles[0xDD]},
	0xD9: {"CMP_ABSY", "CMP", AbsoluteY, Cycles[0xD9]},
	0xC1: {"CMP_INDX", "CMP", PreIndexedIndirect, Cycles[0xC1]},
	0xD1: {"CMP_INDY", "CMP", PostIndexedIndirect, Cycles[0xD1]},
	0x29: {"AND_IMM", "AND", Immediate, Cycles[0x29]},
	0x25: {"AND_ZERO", "AND", ZeroPage, Cycles[0x25]},
	0x2D: {"AND_ABS", "AND", Absolute, Cycles[0x2D]},
	0x35: {"AND_ZEROX", "AND", ZeroPageX, Cycles[0x35]},
	0x3D: {"AND_ABSX", "AND", AbsoluteX, Cycles[0x3D]},
	0x39: {"AND_ABSY", "AND", AbsoluteY, Cycles[0x39]},
	0x21: {"AND_INDX", "AND", PreIndexedIndirect, Cycles[0x21]},
	0x31: {"AND_INDY", "AND", PostIndexedIndirect, Cycles[0x31]},
	0x49: {"EOR_IMM", "EOR", Immediate, Cycles[0x49]},
	0x45: {"EOR_ZERO", "EOR", ZeroPage, Cycles[0x45]},
	0x4D: {"EOR_ABS", "EOR", Absolute, Cycles[0x4D]},
	0x55: {"EOR_ZEROX", "EOR", ZeroPageX, Cycles[0x55]},
	0x5D: {"EOR_ABSX", "EOR", AbsoluteX, Cycles[0x5D]},
	0x59: {"EOR_ABSY", "EOR", AbsoluteY, Cycles[0x59]},
	0x41: {"EOR_INDX", "EOR", PreIndexedIndirect, Cycles[0x41]},
	0x51: {"EOR_INDY", "EOR", PostIndexedIndirect, Cycles[0x51]},
	0x9: {"ORA_IMM", "ORA", Immediate, Cycles[0x09]},
	0x5: {"ORA_ZERO", "ORA", ZeroPage, Cycles[0x05]},
	0xD: {"ORA_ABS", "ORA", Absolute, Cycles[0x0D]},
	0x15: {"ORA_ZEROX", "ORA", ZeroPageX, Cycles[0x15]},
	0x1D: {"ORA_ABSX", "ORA", AbsoluteX, Cycles[0x1D]},
	0x19: {"ORA_ABSY", "ORA", AbsoluteY, Cycles[0x19]},
	0x1: {"ORA_INDX", "ORA", PreIndexedIndirect, Cycles[0x01]},
	0x11: {"ORA_INDY", "ORA", PostIndexedIndirect, Cycles[0x11]},
	0x24: {"BIT_ZERO", "BIT", ZeroPage, Cycles[0x24]},
	0x2C: {"BIT_ABS", "BIT", Absolute, Cycles[0x2C]},
	0xA: {"ASL", "ASL", Accumulator, Cycles[0x0A]},
	0x6: {"ASL_ZERO", "ASL", ZeroPage, Cycles[0x06]},
	0xE: {"ASL_ABS", "ASL", Absolute, Cycles[0x0E]},
	0x16: {"ASL_ZEROX", "ASL", ZeroPageX, Cycles[0x16]},
	0x1E: {"ASL_ABSX", "ASL", AbsoluteX, Cycles[0x1E]},
	0x4A: {"LSR", "LSR", Accumulator, Cycles[0x4A]},
	0x46: {"LSR_ZERO", "LSR", ZeroPage, Cycles[0x46]},
	0x4E: {"LSR_ABS", "LSR", Absolute, Cycles[0x4E]},
	0x56: {"LSR_ZEROX", "LSR", ZeroPageX, Cycles[0x56]},
	0x5E: {"LSR_ABSX", "LSR", AbsoluteX, Cycles[0x5E]},
	0x2A: {"ROL", "ROL", Accumulator, Cycles[0x2A]},
	0x26: {"ROL_ZERO", "ROL", ZeroPage, Cycles[0x26]},
	0x2E: {"ROL_ABS", "ROL", Absolute, Cycles[0x2E]},
	0x36: {"ROL_ZEROX", "ROL", ZeroPageX, Cycles[0x36]},
	0x3E: {"ROL_ABSX", "ROL", AbsoluteX, Cycles[0x3E]},
	0x6A: {"ROR", "ROR", Accumulator, Cycles[0x6A]},
	0x66: {"ROR_ZERO", "ROR", ZeroPage, Cycles[0x66]},
	0x6E: {"ROR_ABS", "ROR", Absolute, Cycles[0x6E]},
	0x76: {"ROR_ZEROX", "ROR", ZeroPageX, Cycles[0x76]},
	0x7E: {"ROR_ABSX", "ROR", AbsoluteX, Cycles[0x7E]},
	0xE8: {"INX", "INX", Implied, Cycles[0xE8]},
	0xC8: {"INY", "INY", Implied, Cycles[0xC8]},
	0xE6: {"INC_ZERO", "INC", ZeroPage, Cycles[0xE6]},
	0xEE: {"INC_ABS", "INC", Absolute, Cycles[0xEE]},
	0xF6: {"INC_ZEROX", "INC", ZeroPageX, Cycles[0xF6]},
	0xFE: {"INC_ABSX", "INC", AbsoluteX, Cycles[0xFE]},
	0xCA: {"DEX", "DEX", Implied, Cycles[0xCA]},
	0x88: {"DEY", "DEY", Implied, Cycles[0x88]},
	0xC6: {"DEC_ZERO", "DEC", ZeroPage, Cycles[0xC6]},
	0xCE: {"DEC_ABS", "DEC", Absolute, Cycles[0xCE]},
	0xD6: {"DEC_ZEROX", "DEC", ZeroPageX, Cycles[0xD6]},
	0xDE: {"DEC_ABSX", "DEC", AbsoluteX, Cycles[0xDE]},
	0x18: {"CLC", "CLC", Implied, Cycles[0x18]},
	0x58: {"CLI", "CLI", Implied, Cycles[0x58]},
	0xB8: {"CLV", "CLV", Implied, Cycles[0xB8]},
	0x38: {"SEC", "SEC", Implied, Cycles[0x38]},
	0x78: {"SEI", "SEI", Implied, Cycles[0x78]},
	0xEA: {"NOP", "NOP", Implied, Cycles[0xEA]},
	0x0: {"BRK", "BRK", Implied, Cycles[0x00]},
	0x20: {"JSR_ABS", "JSR", Absolute, Cycles[0x20]},
	0x4C: {"JMP_ABS", "JMP", Absolute, Cycles[0x4C]},
	0x6C: {"JMP_INDABS", "JMP", IndirectAbsolute, Cycles[0x6C]},
	0x40: {"RTI", "RTI", Implied, Cycles[0x40]},
	0x60: {"RTS", "RTS", Implied, Cycles[0x60]},
	0x10: {"BPL", "BPL", Relative, Cycles[0x10]},
	0x30: {"BMI", "BMI", Relative, Cycles[0x30]},
	0x50: {"BVC", "BVC", Relative, Cycles[0x50]},
	0x70: {"BVS", "BVS", Relative, Cycles[0x70]},
	0x90: {"BCC", "BCC", Relative, Cycles[0x90]},
	0xB0: {"BCS", "BCS", Relative, Cycles[0xB0]},
	0xD0: {"BNE", "BNE", Relative, Cycles[0xD0]},
	0xF0: {"BEQ", "BEQ", Relative, Cycles[0xF0]},
	0xF8: {"SED", "SED", Implied, Cycles[0xF8]},
	0xD8: {"CLD", "CLD", Implied, Cycles[0xD8]},
	// unofficial opecode
	// Also see https://wiki.nesdev.com/w/index.php/CPU_unofficial_opcodes
	0x1A: {"NOP", "NOP", Implied, Cycles[0x1A]},
	0x3A: {"NOP", "NOP", Implied, Cycles[0x3A]},
	0x5A: {"NOP", "NOP", Implied, Cycles[0x5A]},
	0x7A: {"NOP", "NOP", Implied, Cycles[0x7A]},
	0xDA: {"NOP", "NOP", Implied, Cycles[0xDA]},
	0xFA: {"NOP", "NOP", Implied, Cycles[0xFA]},
	0x02: {"NOP", "NOP", Implied, Cycles[0x02]},
	0x12: {"NOP", "NOP", Implied, Cycles[0x12]},
	0x22: {"NOP", "NOP", Implied, Cycles[0x22]},
	0x32: {"NOP", "NOP", Implied, Cycles[0x32]},
	0x42: {"NOP", "NOP", Implied, Cycles[0x42]},
	0x52: {"NOP", "NOP", Implied, Cycles[0x52]},
	0x62: {"NOP", "NOP", Implied, Cycles[0x62]},
	0x72: {"NOP", "NOP", Implied, Cycles[0x72]},
	0x92: {"NOP", "NOP", Implied, Cycles[0x92]},
	0xB2: {"NOP", "NOP", Implied, Cycles[0xB2]},
	0xD2: {"NOP", "NOP", Implied, Cycles[0xD2]},
	0xF2: {"NOP", "NOP", Implied, Cycles[0xF2]},
	0x80: {"NOPD", "NOPD", Implied, Cycles[0x80]},
	0x82: {"NOPD", "NOPD", Implied, Cycles[0x82]},
	0x89: {"NOPD", "NOPD", Implied, Cycles[0x89]},
	0xC2: {"NOPD", "NOPD", Implied, Cycles[0xC2]},
	0xE2: {"NOPD", "NOPD", Implied, Cycles[0xE2]},
	0x04: {"NOPD", "NOPD", Implied, Cycles[0x04]},
	0x44: {"NOPD", "NOPD", Implied, Cycles[0x44]},
	0x64: {"NOPD", "NOPD", Implied, Cycles[0x64]},
	0x14: {"NOPD", "NOPD", Implied, Cycles[0x14]},
	0x34: {"NOPD", "NOPD", Implied, Cycles[0x34]},
	0x54: {"NOPD", "NOPD", Implied, Cycles[0x54]},
	0x74: {"NOPD", "NOPD", Implied, Cycles[0x74]},
	0xD4: {"NOPD", "NOPD", Implied, Cycles[0xD4]},
	0xF4: {"NOPD", "NOPD", Implied, Cycles[0xF4]},
	0x0C: {"NOPI", "NOPI", Implied, Cycles[0x0C]},
	0x1C: {"NOPI", "NOPI", Implied, Cycles[0x1C]},
	0x3C: {"NOPI", "NOPI", Implied, Cycles[0x3C]},
	0x5C: {"NOPI", "NOPI", Implied, Cycles[0x5C]},
	0x7C: {"NOPI", "NOPI", Implied, Cycles[0x7C]},
	0xDC: {"NOPI", "NOPI", Implied, Cycles[0xDC]},
	0xFC: {"NOPI", "NOPI", Implied, Cycles[0xFC]},
	// LAX
	0xA7: {"LAX_ZERO", "LAX", ZeroPage, Cycles[0xA7]},
	0xB7: {"LAX_ZEROY", "LAX", ZeroPageY, Cycles[0xB7]},
	0xAF: {"LAX_ABS", "LAX", Absolute, Cycles[0xAF]},
	0xBF: {"LAX_ABSY", "LAX", AbsoluteY, Cycles[0xBF]},
	0xA3: {"LAX_INDX", "LAX", PreIndexedIndirect, Cycles[0xA3]},
	0xB3: {"LAX_INDY", "LAX", PostIndexedIndirect, Cycles[0xB3]},
	// SAX
	0x87: {"SAX_ZERO", "SAX", ZeroPage, Cycles[0x87]},
	0x97: {"SAX_ZEROY", "SAX", ZeroPageY, Cycles[0x97]},
	0x8F: {"SAX_ABS", "SAX", Absolute, Cycles[0x8F]},
	0x83: {"SAX_INDX", "SAX", PreIndexedIndirect, Cycles[0x83]},
	// SBC
	0xEB: {"SBC_IMM", "SBC", Immediate, Cycles[0xEB]},
	// DCP
	0xC7: {"DCP_ZERO", "DCP", ZeroPage, Cycles[0xC7]},
	0xD7: {"DCP_ZEROX", "DCP", ZeroPageX, Cycles[0xD7]},
	0xCF: {"DCP_ABS", "DCP", Absolute, Cycles[0xCF]},
	0xDF: {"DCP_ABSX", "DCP", AbsoluteX, Cycles[0xDF]},
	0xDB: {"DCP_ABSY", "DCP", AbsoluteY, Cycles[0xD8]},
	0xC3: {"DCP_INDX", "DCP", PreIndexedIndirect, Cycles[0xC3]},
	0xD3: {"DCP_INDY", "DCP", PostIndexedIndirect, Cycles[0xD3]},
	// ISB
	0xE7: {"ISB_ZERO", "ISB", ZeroPage, Cycles[0xE7]},
	0xF7: {"ISB_ZEROX", "ISB", ZeroPageX, Cycles[0xF7]},
	0xEF: {"ISB_ABS", "ISB", Absolute, Cycles[0xEF]},
	0xFF: {"ISB_ABSX", "ISB", AbsoluteX, Cycles[0xFF]},
	0xFB: {"ISB_ABSY", "ISB", AbsoluteY, Cycles[0xF8]},
	0xE3: {"ISB_INDX", "ISB", PreIndexedIndirect, Cycles[0xE3]},
	0xF3: {"ISB_INDY", "ISB", PostIndexedIndirect, Cycles[0xF3]},
	// SLO
	0x07: {"SLO_ZERO", "SLO", ZeroPage, Cycles[0x07]},
	0x17: {"SLO_ZEROX", "SLO", ZeroPageX, Cycles[0x17]},
	0x0F: {"SLO_ABS", "SLO", Absolute, Cycles[0x0F]},
	0x1F: {"SLO_ABSX", "SLO", AbsoluteX, Cycles[0x1F]},
	0x1B: {"SLO_ABSY", "SLO", AbsoluteY, Cycles[0x1B]},
	0x03: {"SLO_INDX", "SLO", PreIndexedIndirect, Cycles[0x03]},
	0x13: {"SLO_INDY", "SLO", PostIndexedIndirect, Cycles[0x13]},
	// RLA
	0x27: {"RLA_ZERO", "RLA", ZeroPage, Cycles[0x27]},
	0x37: {"RLA_ZEROX", "RLA", ZeroPageX, Cycles[0x37]},
	0x2F: {"RLA_ABS", "RLA", Absolute, Cycles[0x2F]},
	0x3F: {"RLA_ABSX", "RLA", AbsoluteX, Cycles[0x3F]},
	0x3B: {"RLA_ABSY", "RLA", AbsoluteY, Cycles[0x3B]},
	0x23: {"RLA_INDX", "RLA", PreIndexedIndirect, Cycles[0x23]},
	0x33: {"RLA_INDY", "RLA", PostIndexedIndirect, Cycles[0x33]},
	// SRE
	0x47: {"SRE_ZERO", "SRE", ZeroPage, Cycles[0x47]},
	0x57: {"SRE_ZEROX", "SRE", ZeroPageX, Cycles[0x57]},
	0x4F: {"SRE_ABS", "SRE", Absolute, Cycles[0x4F]},
	0x5F: {"SRE_ABSX", "SRE", AbsoluteX, Cycles[0x5F]},
	0x5B: {"SRE_ABSY", "SRE", AbsoluteY, Cycles[0x5B]},
	0x43: {"SRE_INDX", "SRE", PreIndexedIndirect, Cycles[0x43]},
	0x53: {"SRE_INDY", "SRE", PostIndexedIndirect, Cycles[0x53]},
	// RRA
	0x67: {"RRA_ZERO", "RRA", ZeroPage, Cycles[0x67]},
	0x77: {"RRA_ZEROX", "RRA", ZeroPageX, Cycles[0x77]},
	0x6F: {"RRA_ABS", "RRA", Absolute, Cycles[0x6F]},
	0x7F: {"RRA_ABSX", "RRA", AbsoluteX, Cycles[0x7F]},
	0x7B: {"RRA_ABSY", "RRA", AbsoluteY, Cycles[0x7B]},
	0x63: {"RRA_INDX", "RRA", PreIndexedIndirect, Cycles[0x63]},
	0x73: {"RRA_INDY", "RRA", PostIndexedIndirect, Cycles[0x73]},
}