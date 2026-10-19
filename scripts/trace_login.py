import struct
from capstone import Cs, CS_ARCH_X86, CS_MODE_32

# Load PDB symbols
with open('pdo/Zone.pdb', 'rb') as f:
    pdb = f.read()

symbols = {}
pos = 0
while True:
    idx = pdb.find(b'\x0e\x11', pos)
    if idx == -1:
        break
    if idx >= 2:
        reclen = struct.unpack_from('<H', pdb, idx-2)[0]
        if reclen >= 10:
            flags, off, seg = struct.unpack_from('<IIH', pdb, idx+2)
            name_end = pdb.find(b'\x00', idx+12)
            if name_end != -1 and name_end <= idx + reclen:
                try:
                    name = pdb[idx+12:name_end].decode('latin1')
                    if seg == 1:
                        symbols[0x401000 + off] = name
                except:
                    pass
    pos = idx + 2

sorted_syms = sorted(symbols.items())

def find_enclosing_symbol(va):
    prev = None
    for s_va, s_name in sorted_syms:
        if s_va > va:
            break
        prev = (s_va, s_name)
    return prev

# Load Zone.exe
with open('../NA2016/Server/Zone00/Zone.exe', 'rb') as f:
    exe = f.read()

md = Cs(CS_ARCH_X86, CS_MODE_32)

def disassemble_func(va, max_len=1500):
    rva = va - 0x400000
    foff = 0x51990 + (va - 0x452590) # based on .text section
    code = exe[foff:foff+max_len]
    print(f"\nDisassembly of {symbols.get(va, hex(va))} (0x{va:x}):")
    for insn in md.disasm(code, va):
        sym_str = ""
        if insn.mnemonic in ('call', 'jmp'):
            try:
                target = int(insn.op_str, 16)
                if target in symbols:
                    sym_str = f" ; {symbols[target]}"
            except:
                pass
        # Check for immediate push / constants
        elif 'push' in insn.mnemonic:
            try:
                val = int(insn.op_str, 16)
                if val in symbols:
                    sym_str = f" ; {symbols[val]}"
                elif val > 0x1000 and val < 0x9000:
                    dept = val >> 10
                    cmd = val & 0x3ff
                    sym_str = f" ; OpCode 0x{val:04x} (dept={dept}, cmd={cmd})"
            except:
                pass
        print(f"  0x{insn.address:x}: {insn.mnemonic:8s} {insn.op_str:24s} {sym_str}")
        if insn.mnemonic == 'ret':
            break

if __name__ == '__main__':
    disassemble_func(0x452590)
