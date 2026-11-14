#!/usr/bin/env python3
"""
disasm_pdb_handler.py
Desensamblador oficial asistido por PDB para XelionCore (Etapa 4).
Cruza los símbolos de pdo/*.pdb (o ../AstraServer/Server/*/*.pdb) con el binario
PE original (.exe) usando /usr/bin/objdump y anota automáticamente cada CALL
con el nombre de símbolo C++ desmanglado del PDB.

Uso:
  python3 scripts/disasm_pdb_handler.py zone sp_NC_ACT_CHAT_REQ [bytes=512]
  python3 scripts/disasm_pdb_handler.py world fc_NC_CHAR_LOGIN_REQ [bytes=512]
  python3 scripts/disasm_pdb_handler.py zone --search MoveRun
"""

import re
import struct
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

TARGETS = {
    "zone": {
        "pdb": [ROOT / "pdo" / "Zone.pdb", ROOT.parent / "AstraServer" / "Server" / "Zone00" / "Zone.pdb"],
        "exe": [ROOT.parent / "AstraServer" / "Server" / "Zone00" / "Zone.exe"],
        "image_base": 0x00400000,
        "text_rva": 0x00001000,
    },
    "world": {
        "pdb": [ROOT / "pdo" / "WorldManager.pdb", ROOT.parent / "AstraServer" / "Server" / "WorldManager" / "WorldManager.pdb"],
        "exe": [ROOT.parent / "AstraServer" / "Server" / "WorldManager" / "WorldManager.exe"],
        "image_base": 0x00400000,
        "text_rva": 0x00001000,
    },
    "character": {
        "pdb": [ROOT.parent / "AstraServer" / "Server" / "Character" / "Character.pdb"],
        "exe": [ROOT.parent / "AstraServer" / "Server" / "Character" / "Character.exe"],
        "image_base": 0x00400000,
        "text_rva": 0x00001000,
    },
    "login": {
        "pdb": [ROOT / "pdo" / "Login.pdb", ROOT.parent / "AstraServer" / "Server" / "Login" / "Login.pdb"],
        "exe": [ROOT.parent / "AstraServer" / "Server" / "Login" / "Login.exe"],
        "image_base": 0x00400000,
        "text_rva": 0x00001000,
    },
}


def pick_existing(paths):
    for p in paths:
        if p.exists():
            return p
    raise FileNotFoundError(f"Ninguno de los archivos existe: {[str(p) for p in paths]}")


def load_pdb_symbols(pdb_path: Path, image_base: int, text_rva: int):
    data = pdb_path.read_bytes()
    va_to_sym = {}
    name_to_vas = {}
    pos = 0
    while True:
        idx = data.find(b"\x0e\x11", pos)
        if idx == -1:
            break
        if idx >= 2:
            reclen = struct.unpack_from("<H", data, idx - 2)[0]
            if reclen >= 10:
                _flags, off, seg = struct.unpack_from("<IIH", data, idx + 2)
                name_end = data.find(b"\x00", idx + 12)
                if name_end != -1 and name_end <= idx + reclen:
                    try:
                        raw_name = data[idx + 12 : name_end].decode("latin1")
                        if seg == 1:
                            va = image_base + text_rva + off
                            clean = raw_name.split("@")[0].lstrip("?_")
                            va_to_sym[va] = f"{clean} ({raw_name})"
                            name_to_vas.setdefault(clean, []).append((va, raw_name))
                            name_to_vas.setdefault(raw_name, []).append((va, raw_name))
                    except Exception:
                        pass
        pos = idx + 2
    return va_to_sym, name_to_vas


def pe_va_to_offset(exe_bytes: bytes, va: int) -> int:
    pe_off = struct.unpack_from("<I", exe_bytes, 0x3C)[0]
    num_sec = struct.unpack_from("<H", exe_bytes, pe_off + 6)[0]
    opt_size = struct.unpack_from("<H", exe_bytes, pe_off + 20)[0]
    image_base = struct.unpack_from("<I", exe_bytes, pe_off + 24 + 28)[0]
    rva = va - image_base
    sec_table = pe_off + 24 + opt_size
    for i in range(num_sec):
        s_off = sec_table + i * 40
        vsize, vaddr, rawsize, rawptr = struct.unpack_from("<IIII", exe_bytes, s_off + 8)
        span = max(vsize, rawsize)
        if vaddr <= rva < vaddr + span:
            return rawptr + (rva - vaddr)
    return 0x400 + (rva - 0x1000)


def build_elf32_for_code(code: bytes, start_va: int) -> bytes:
    # Minimal 32-bit x86 ELF with 1 .text section at sh_addr = start_va
    shstrtab = b"\x00.text\x00.shstrtab\x00"
    ehdr_size = 52
    code_off = 64
    shstr_off = code_off + len(code)
    # align to 4
    while (shstr_off + len(shstrtab)) % 4 != 0:
        shstrtab += b"\x00"
    shdr_off = shstr_off + len(shstrtab)

    e_ident = b"\x7fELF\x01\x01\x01" + b"\x00" * 9
    ehdr = struct.pack(
        "<16sHHIIIIIHHHHHH",
        e_ident,
        1,          # ET_REL / ET_EXEC
        3,          # EM_386
        1,          # EV_CURRENT
        start_va,   # e_entry
        0,          # e_phoff
        shdr_off,   # e_shoff
        0,          # e_flags
        ehdr_size,
        0,
        0,
        40,         # e_shentsize
        3,          # e_shnum (NULL, .text, .shstrtab)
        2,          # e_shstrndx
    )
    pad = b"\x00" * (code_off - len(ehdr))
    sh_null = b"\x00" * 40
    sh_text = struct.pack(
        "<IIIIIIIIII",
        1,          # sh_name (.text)
        1,          # SHT_PROGBITS
        6,          # SHF_ALLOC | SHF_EXECINSTR
        start_va,   # sh_addr
        code_off,   # sh_offset
        len(code),  # sh_size
        0, 0, 16, 0
    )
    sh_str = struct.pack(
        "<IIIIIIIIII",
        7,          # sh_name (.shstrtab)
        3,          # SHT_STRTAB
        0,
        0,
        shstr_off,
        len(shstrtab),
        0, 0, 1, 0
    )
    return ehdr + pad + code + shstrtab + sh_null + sh_text + sh_str


def disassemble_range(exe_path: Path, start_va: int, length: int, va_to_sym: dict):
    exe_bytes = exe_path.read_bytes()
    file_off = pe_va_to_offset(exe_bytes, start_va)
    code = exe_bytes[file_off : file_off + length]
    elf_bytes = build_elf32_for_code(code, start_va)

    with tempfile.NamedTemporaryFile(suffix=".elf", delete=True) as tmp:
        tmp.write(elf_bytes)
        tmp.flush()
        cmd = [
            "/usr/bin/objdump",
            "-d",
            "--x86-asm-syntax=intel",
            tmp.name,
        ]
        out = subprocess.check_output(cmd, text=True, errors="replace")

    call_pat = re.compile(r"\b(call|jmp|je|jne|ja|jb|jg|jl|jge|jle)\s+0x([0-9a-fA-F]+)")
    lines = []
    for line in out.splitlines():
        if "file format" in line or "Disassembly of section" in line or "<.text>:" in line:
            continue
        m = call_pat.search(line)
        if m:
            target_va = int(m.group(2), 16)
            sym = va_to_sym.get(target_va)
            if sym:
                line = f"{line}    ; -> {sym}"
        lines.append(line)
    return "\n".join(lines).strip()


def main():
    if len(sys.argv) < 3:
        print(__doc__)
        sys.exit(1)

    target_key = sys.argv[1].lower()
    if target_key not in TARGETS:
        print(f"Target inválido '{target_key}'. Opciones: {list(TARGETS.keys())}")
        sys.exit(1)

    cfg = TARGETS[target_key]
    pdb_path = pick_existing(cfg["pdb"])
    exe_path = pick_existing(cfg["exe"])
    va_to_sym, name_to_vas = load_pdb_symbols(pdb_path, cfg["image_base"], cfg["text_rva"])

    if sys.argv[2] == "--search":
        query = sys.argv[3].lower() if len(sys.argv) > 3 else ""
        matches = []
        for va, desc in sorted(va_to_sym.items()):
            if query in desc.lower():
                matches.append((va, desc))
        print(f" Encontrados {len(matches)} símbolos para '{query}' en {pdb_path.name}:")
        for va, desc in matches[:100]:
            print(f"  0x{va:08X}  {desc}")
        return

    sym_arg = sys.argv[2]
    size = int(sys.argv[3]) if len(sys.argv) > 3 else 512

    if sym_arg.startswith("0x") or sym_arg.startswith("0X"):
        va = int(sym_arg, 16)
        print(f"=== {exe_path.name} @ 0x{va:08X} ({va_to_sym.get(va, 'unknown')}) ===")
        print(disassemble_range(exe_path, va, size, va_to_sym))
        return

    candidates = name_to_vas.get(sym_arg, [])
    if not candidates:
        q = sym_arg.lower()
        for va, desc in sorted(va_to_sym.items()):
            if q in desc.lower():
                candidates.append((va, desc))

    if not candidates:
        print(f"No se encontró el símbolo '{sym_arg}' en {pdb_path.name}")
        sys.exit(1)

    for va, raw in candidates[:3]:
        print(f"=== {exe_path.name} :: {raw} @ 0x{va:08X} (len={size}) ===")
        print(disassemble_range(exe_path, va, size, va_to_sym))
        print()


if __name__ == "__main__":
    main()
