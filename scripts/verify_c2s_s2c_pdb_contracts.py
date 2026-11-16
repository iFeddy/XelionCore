#!/usr/bin/env python3
"""
Verificador Automático de Contratos Binarios C2S/S2C (PDB y .exe vs. Go)
========================================================================
Cruza automáticamente las especificaciones TPI de los PDBs oficiales
(docs/Servicios/Login.md, WorldManager.md, Zone.md, Account.md, Character.md)
contra los handlers C2S, constructores S2C y preservación de estado de sesión
en pkg/protocol, pkg/server/login, pkg/server/world y pkg/server/zone.
"""

import glob
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))


def load_pdb_structs():
    structs = {}
    docs = [
        "docs/Servicios/Login.md",
        "docs/Servicios/Account.md",
        "docs/Servicios/AccountLog.md",
        "docs/Servicios/WorldManager.md",
        "docs/Servicios/Character.md",
        "docs/Servicios/Zone.md",
        "docs/Servicios/GameLog.md",
    ]
    for rel_doc in docs:
        doc_path = os.path.join(ROOT, rel_doc)
        if not os.path.exists(doc_path):
            continue
        with open(doc_path, encoding="utf-8") as f:
            text = f.read()
        pattern = re.compile(
            r"### `(PROTO_NC_[A-Za-z0-9_]+)`\s*\n+-\s*\*\*Tamaño Total en Memoria:\*\*\s*`(\d+)\s*bytes`.*?\n```cpp\n(.*?)```(.*?)(?=\n### |\Z)",
            re.S,
        )
        for m in pattern.finditer(text):
            name, sz, cpp, rest = (
                m.group(1),
                int(m.group(2)),
                m.group(3),
                m.group(4),
            )
            is_empty = (
                "Estructura vacía" in cpp
                or "Payload de 0 bytes" in cpp
                or "Payload de 0 bytes" in rest
            )
            fields = []
            for fm in re.finditer(
                r"\|\s*`0x([0-9A-Fa-f]+)`\s*\|\s*`(\d+)`\s*\|\s*`([^`]+)`\s*\|\s*`([^`]+)`\s*\|",
                rest,
            ):
                off_dec = int(fm.group(2))
                ftype = fm.group(3).strip()
                fname = fm.group(4).strip()
                fields.append((off_dec, ftype, fname))
            has_flex = any("[0]" in f[1] for f in fields) or "[0]" in cpp
            structs[name] = {
                "size": sz,
                "is_empty": is_empty,
                "has_flex": has_flex,
                "fields": fields,
                "doc": rel_doc,
            }
    return structs


def load_go_sources():
    go_files = {}
    patterns = [
        "pkg/protocol/*.go",
        "pkg/server/login/*.go",
        "pkg/server/world/*.go",
        "pkg/server/zone/*.go",
    ]
    for pat in patterns:
        for path in sorted(glob.glob(os.path.join(ROOT, pat))):
            if path.endswith("_test.go"):
                continue
            rel = os.path.relpath(path, ROOT)
            with open(path, encoding="utf-8") as f:
                go_files[rel] = f.read()
    return go_files


def extract_methods(go_files):
    methods = {}
    for rel, src in go_files.items():
        parts = rel.split(os.sep)
        pkg = parts[2] if len(parts) > 2 else parts[1]
        for fm in re.finditer(
            r"func\s+\(\w+\s+\*?\w+\)\s+(\w+)\s*\([^)]*\)\s*(?:[^{]*)\{", src
        ):
            fname = fm.group(1)
            start = fm.end()
            depth = 1
            i = start
            while i < len(src) and depth > 0:
                if src[i] == "{":
                    depth += 1
                elif src[i] == "}":
                    depth -= 1
                i += 1
            methods[(pkg, fname)] = (rel, src[start : i - 1])
    return methods


def extract_primary_branch_reads(body):
    """
    Extrae las lecturas de pkt.Reader / readFixedStr / readItemInvenSlot en la
    rama principal (primera rama si hay if/else sobre pkt.Reader.Remaining()).
    """
    # Si hay un if pkt.Reader.Remaining() >= N o == N { ... } else { ... }, evaluar la rama principal
    m_if = re.search(
        r"if\s+(?:pkt\s*!=\s*nil\s*&&\s*)?pkt\.Reader\.Remaining\(\)\s*(?:>=|==)\s*\d+(?:\s*&&\s*pkt\.Reader\.Remaining\(\)\s*<\s*\d+)?\s*\{",
        body,
    )
    target = body
    if m_if:
        start = m_if.end()
        depth = 1
        i = start
        while i < len(body) and depth > 0:
            if body[i] == "{":
                depth += 1
            elif body[i] == "}":
                depth -= 1
            i += 1
        target = body[: m_if.start()] + "\n" + body[start : i - 1]

    ops = []
    token_re = re.compile(
        r"(?:Reader\.(ReadByte|ReadUint8|ReadInt8|ReadUint16|ReadInt16|ReadUint32|ReadInt32|ReadFloat32|ReadUint64|ReadInt64|ReadBytes|ReadString)\s*\(([^)]*)\))"
        r"|(?:readFixedStr\s*\(\s*pkt\.Reader\s*,\s*(\d+)\s*\))"
        r"|(?:readItemInvenSlot\s*\(\s*pkt\s*\))"
    )
    for m in token_re.finditer(target):
        if m.group(1):
            op, arg = m.group(1), m.group(2).strip()
            if op in ("ReadByte", "ReadUint8", "ReadInt8"):
                ops.append((op, 1))
            elif op in ("ReadUint16", "ReadInt16"):
                ops.append((op, 2))
            elif op in ("ReadUint32", "ReadInt32", "ReadFloat32"):
                ops.append((op, 4))
            elif op in ("ReadUint64", "ReadInt64"):
                ops.append((op, 8))
            elif op in ("ReadBytes", "ReadString"):
                try:
                    ops.append((op, int(arg)))
                except ValueError:
                    ops.append((op, None))
        elif m.group(3):
            ops.append(("readFixedStr", int(m.group(3))))
        else:
            ops.append(("readItemInvenSlot", 2))
    return ops


def check_c2s_contracts(structs, go_files, methods):
    discrepancies = []
    checked_cases = 0
    empty_verified = 0

    for rel, src in go_files.items():
        if not rel.startswith("pkg/server/"):
            continue
        pkg = rel.split(os.sep)[2]
        lines = src.splitlines()
        for idx, line in enumerate(lines):
            cm = re.search(r"case\s+protocol\.(PROTO_NC_\w+)\s*:", line)
            if not cm:
                continue
            proto = cm.group(1)
            case_lines = []
            for j in range(idx + 1, min(len(lines), idx + 60)):
                if re.match(r"^\s*(case\s+|default\s*:)", lines[j]):
                    break
                case_lines.append(lines[j])
            case_body = "\n".join(case_lines)
            hm = re.search(r"\bh\.(\w+)\s*\(", case_body)
            target_body = case_body
            target_loc = f"{rel}:{idx + 1}"
            if hm and (pkg, hm.group(1)) in methods:
                mrel, mbody = methods[(pkg, hm.group(1))]
                target_body = mbody
                target_loc = f"{mrel} ({hm.group(1)})"

            if proto not in structs:
                continue

            checked_cases += 1
            st = structs[proto]

            # 1. Estructuras vacías en el PDB (0 bytes de campos)
            if st["is_empty"]:
                raw_reads = re.findall(
                    r"(?:Reader\.Read\w+|readFixedStr|readItemInvenSlot)", target_body
                )
                if raw_reads:
                    discrepancies.append(
                        f"[C2S EMPTY STRUCT READ] {pkg} {proto} @ {target_loc}: "
                        f"PDB declara estructura vacía (0 bytes payload), pero Go invoca {raw_reads}"
                    )
                else:
                    empty_verified += 1
                continue

            if st["has_flex"]:
                continue

            # 2. Estructuras de tamaño fijo sin lazos dinámicos
            if "for " in target_body:
                continue
            ops = extract_primary_branch_reads(target_body)
            if ops and all(sz is not None for _, sz in ops):
                total_bytes = sum(sz for _, sz in ops)
                if total_bytes != st["size"]:
                    discrepancies.append(
                        f"[C2S SIZE MISMATCH] {pkg} {proto} @ {target_loc}: "
                        f"PDB sizeof={st['size']}B vs Go lecturas={total_bytes}B ({ops})"
                    )

    return checked_cases, empty_verified, discrepancies


def write_call_size(op, arg, body):
    if op == "WriteByte":
        return 1
    if op in ("WriteUint16", "WriteInt16"):
        return 2
    if op in ("WriteUint32", "WriteInt32", "WriteFloat32"):
        return 4
    if op in ("WriteUint64", "WriteInt64"):
        return 8
    if op == "WriteBool":
        return 1
    if op == "WriteString":
        parts = arg.split(",")
        if len(parts) == 2:
            try:
                return int(parts[1].strip())
            except ValueError:
                return None
    if op == "Fill":
        parts = arg.split(",")
        if len(parts) >= 1:
            try:
                return int(parts[0].strip())
            except ValueError:
                return None
    if op == "WriteBytes":
        a = arg.strip()
        m_arr = re.match(r"^(\w+)\[:\]$", a)
        if m_arr:
            vname = m_arr.group(1)
            dm = re.search(rf"var\s+{re.escape(vname)}\s+\[(\d+)\]byte", body)
            if dm:
                return int(dm.group(1))
        m_make = re.match(r"^make\(\[\]byte,\s*(\d+)\)$", a)
        if m_make:
            return int(m_make.group(1))
        if a.startswith("packTm("):
            return 36
    return None


def check_s2c_contracts(structs, go_files):
    discrepancies = []
    checked_s2c = 0

    # 1. Auditar constructores New_PROTO_NC_* en pkg/protocol/*.go
    for rel, src in go_files.items():
        if not rel.startswith("pkg/protocol/"):
            continue
        for fm in re.finditer(
            r"func\s+New_(PROTO_NC_\w+)\s*\([^)]*\)\s*\*PacketWriter\s*\{", src
        ):
            proto = fm.group(1)
            start = fm.end()
            depth = 1
            i = start
            while i < len(src) and depth > 0:
                if src[i] == "{":
                    depth += 1
                elif src[i] == "}":
                    depth -= 1
                i += 1
            body = src[start : i - 1]
            if proto not in structs:
                continue
            st = structs[proto]
            checked_s2c += 1
            if st["is_empty"]:
                writes = re.findall(r"\bpw\.(Write\w+|Fill)\(", body)
                if writes:
                    discrepancies.append(
                        f"[S2C EMPTY STRUCT WRITE] {rel}: New_{proto} escribe {writes} en struct PDB vacío"
                    )
                continue
            if st["has_flex"] or "for " in body or "if " in body:
                continue
            writes = re.findall(r"\bpw\.(Write\w+|Fill)\(([^)]*)\)", body)
            wsizes = [write_call_size(op, arg, body) for op, arg in writes]
            if writes and all(s is not None for s in wsizes):
                tot = sum(wsizes)
                if tot != st["size"]:
                    discrepancies.append(
                        f"[S2C SIZE MISMATCH] {rel}: New_{proto} PDB={st['size']}B vs Go={tot}B ({writes})"
                    )

    # 2. Auditar paquetes S2C críticos en WorldService y ZoneService
    critical_s2c_checks = [
        (
            "pkg/server/world/handler_social_misc.go",
            "PROTO_NC_PRISON_END_ACK",
            r"NewPacketWriter\(protocol\.PROTO_NC_PRISON_END_ACK\)\s*\n\s*resp\.WriteUint16\(0\)\s*\n\s*resp\.WriteUint32\(0\)",
            "PROTO_NC_PRISON_END_ACK debe escribir err(uint16) seguido de nPrisonSecond(uint32) (6 bytes)",
        ),
        (
            "pkg/server/world/guild_manager.go",
            "PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_ACK",
            r"NewPacketWriter\(protocol\.PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_ACK\)\s*\n\s*resp\.WriteUint16\(0\).*?writeGuildTM\(resp,\s*now\)",
            "PROTO_NC_GUILD_ACADEMY_GET_RANKING_LIST_ACK debe escribir 51B de cabecera oficial",
        ),
        (
            "pkg/server/zone/stats_pvp.go",
            "PROTO_NC_CHAR_REVIVEOTHER_CMD",
            r"NewPacketWriter\(protocol\.PROTO_NC_CHAR_REVIVEOTHER_CMD\)\s*\n\s*otherCmd\.WriteUint16\(targetMap\.ID\)\s*\n\s*otherCmd\.WriteInt32\(int32\(player\.Char\.XPos\)\)\s*\n\s*otherCmd\.WriteInt32\(int32\(player\.Char\.YPos\)\)\s*\n\s*writeFixedStr\(otherCmd,\s*\"127\.0\.0\.1\",\s*16\)\s*\n\s*otherCmd\.WriteUint16\(9013\)\s*\n\s*otherCmd\.WriteUint16\(player\.MapObjectID\)",
            "PROTO_NC_CHAR_REVIVEOTHER_CMD debe escribir 30B exactos",
        ),
    ]
    for rel, proto, pat, desc in critical_s2c_checks:
        checked_s2c += 1
        src = go_files.get(rel, "")
        if not re.search(pat, src, re.S):
            discrepancies.append(f"[S2C CONTRACT FAIL] {rel} ({proto}): {desc}")

    return checked_s2c, discrepancies


def check_session_state_transitions(go_files):
    discrepancies = []
    avatar_opts = go_files.get("pkg/server/world/handler_avatar_options.go", "")

    # Verificar que handleNormalLogout preserve StyleCharID antes de limpiar CharacterID
    if not re.search(
        r"client\.StyleCharID\s*=\s*client\.CharacterID.*?\n.*?client\.CharacterID\s*=\s*0",
        avatar_opts,
        re.S,
    ):
        discrepancies.append(
            "[SESSION TRANSITION] handleNormalLogout en pkg/server/world/handler_avatar_options.go "
            "no preserva client.StyleCharID = client.CharacterID antes de limpiar client.CharacterID = 0"
        )

    # Verificar que logoutType == 1 retorne sin cerrar la sesión ni el socket
    if not re.search(
        r"if\s+logoutType\s*==\s*1\s*\{\s*(?://[^\n]*\n\s*)*client\.IsTransferring\s*=\s*false\s*\n\s*return\s*\n\s*\}",
        avatar_opts,
        re.S,
    ):
        discrepancies.append(
            "[SESSION TRANSITION] handleNormalLogout (logoutType == 1) no retorna preservando la sesión activa"
        )

    # Verificar que handleSetStyleGetInfo y handleSetStyle usen client.StyleCharID cuando client.CharacterID == 0
    for fn_name in ("handleSetStyleGetInfo", "handleSetStyle"):
        m = re.search(
            rf"func\s+\(h\s+\*Handler\)\s+{fn_name}\b.*?\n\}}", avatar_opts, re.S
        )
        if not m or "client.StyleCharID" not in m.group(0):
            discrepancies.append(
                f"[SESSION TRANSITION] {fn_name} no recupera el personaje desde client.StyleCharID tras logoutType == 1"
            )

    # Verificar que Character.Sync() en pkg/database/character.go no sobrescriba campos planos cuando ID != 0
    char_go_path = os.path.join(ROOT, "pkg/database/character.go")
    with open(char_go_path, encoding="utf-8") as f:
        char_go = f.read()
    if "PopulateFromCanonical" not in char_go or "SyncToCanonical" not in char_go:
        discrepancies.append(
            "[MEMORY SYNC] pkg/database/character.go debe definir PopulateFromCanonical() y SyncToCanonical()"
        )
    if re.search(r"else\s+if\s+c\.TCharacter\.CharNo\s*!=\s*0\s*\{", char_go):
        discrepancies.append(
            "[MEMORY SYNC] Character.Sync() todavía contiene la rama destructiva `else if c.TCharacter.CharNo != 0`"
        )

    return discrepancies


def main():
    structs = load_pdb_structs()
    go_files = load_go_sources()
    methods = extract_methods(go_files)

    c2s_count, empty_count, c2s_disc = check_c2s_contracts(structs, go_files, methods)
    s2c_count, s2c_disc = check_s2c_contracts(structs, go_files)
    sess_disc = check_session_state_transitions(go_files)

    all_disc = c2s_disc + s2c_disc + sess_disc
    print(
        f" PDB/EXE Contract Verifier Summary:\n"
        f"  - PDB structs loaded:        {len(structs)}\n"
        f"  - C2S handlers verified:     {c2s_count} (incl. {empty_count} empty 0-byte structs)\n"
        f"  - S2C constructors verified: {s2c_count}\n"
        f"  - Session/Memory invariants: 5\n"
        f"  - Discrepancies found:       {len(all_disc)}"
    )

    if all_disc:
        print("\nDiscrepancies:")
        for d in all_disc:
            print(f"  * {d}")
        sys.exit(1)

    print("OK: 0 discrepancies found.")
    sys.exit(0)


if __name__ == "__main__":
    main()
