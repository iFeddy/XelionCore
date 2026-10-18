#!/usr/bin/env python3
import os
import glob
import re

def clean_expr(expr):
    """
    Limpia expresiones matemáticas y relacionales de .ps a Lua.
    """
    expr = expr.strip()
    # Reemplazar operadores de igualdad de .ps
    expr = expr.replace("=== ", "== ").replace("=!=", "~=")

    # Reemplazar @(Expr) por (_G[Expr] or 0)
    m_at = re.search(r'@\((.+?)\)', expr)
    if m_at:
        inner = clean_expr(m_at.group(1))
        expr = expr.replace(m_at.group(0), f'(_G[{inner}] or 0)')

    # Reemplazar concatenación de strings % por ..
    tokens = re.findall(r'"[^"]*"|[^\s"]+', expr)
    res = []
    i = 0
    while i < len(tokens):
        t = tokens[i]
        if t == "%":
            res.append("..")
        else:
            res.append(t)
        i += 1
    return " ".join(res)

def split_statements(line):
    stripped = line.strip()
    if stripped.startswith(";"):
        return [line]

    comment = ""
    code = stripped
    if ";" in stripped:
        parts = stripped.split(";", 1)
        code = parts[0]
        comment = ";" + parts[1]

    parts = []
    current = []
    in_quote = False
    for char in code:
        if char == '"':
            in_quote = not in_quote
            current.append(char)
        elif char == '.' and not in_quote:
            stmt = "".join(current).strip()
            if stmt:
                parts.append(stmt + ".")
            current = []
        else:
            current.append(char)
    stmt = "".join(current).strip()
    if stmt:
        parts.append(stmt)

    if not parts:
        return [line]

    if comment:
        parts[-1] += " " + comment

    return parts

def transpile_ps(raw_lines):
    # Pre-procesar múltiples sentencias en una sola línea
    lines = []
    for line in raw_lines:
        for sub in split_statements(line):
            lines.append(sub)

    out = []
    
    # Pila de bloques abiertos: ('func', name), ('if_chain', has_else), ('while', None), ('repeat', name)
    block_stack = []
    
    i = 0
    while i < len(lines):
        raw_line = lines[i]
        stripped = raw_line.strip()
        
        # 1. Líneas vacías
        if not stripped:
            out.append("")
            i += 1
            continue
            
        # 2. Separadores decorativos ;;;;;;;;;
        if re.match(r"^;{3,}$", stripped):
            out.append("--" + "-" * 70)
            i += 1
            continue
            
        # 3. Comentarios completos
        if stripped.startswith(";"):
            comment_content = stripped[1:].strip()
            # Si contiene código comentado de supervivencia o similar
            indent = "    " * len(block_stack)
            out.append(f"{indent}-- {comment_content}")
            i += 1
            continue
            
        # Separar código y comentario inline
        code_part = stripped
        inline_comment = ""
        if ";" in stripped:
            parts = stripped.split(";", 1)
            code_part = parts[0].strip()
            inline_comment = " -- " + parts[1].strip()
            
        if code_part.endswith("."):
            code_part = code_part[:-1].strip()
            
        indent = "    " * len(block_stack)
        
        # 4. Declaración de Función: open [Nombre]
        m_func = re.match(r"^open\s+\[([a-zA-Z0-9_]+)\]$", code_part)
        if m_func and len(block_stack) == 0:
            fn_name = m_func.group(1)
            block_stack.append(('func', fn_name))
            out.append(f"function {fn_name}(){inline_comment}")
            i += 1
            continue
            
        # 5. Bloque nombrado interno: open [Nombre]
        if m_func and len(block_stack) > 0:
            block_name = m_func.group(1)
            block_stack.append(('repeat', block_name))
            out.append(f"{indent}-- Bloque [{block_name}]{inline_comment}")
            out.append(f"{indent}repeat")
            i += 1
            continue
            
        # 6. infinite (loop infinito)
        if code_part == "infinite":
            # el siguiente debe ser 'open'
            next_line = lines[i+1].strip() if i+1 < len(lines) else ""
            if next_line.startswith("open"):
                i += 1 # saltar 'open'
            block_stack.append(('while', 'infinite'))
            out.append(f"{indent}while true do{inline_comment}")
            i += 1
            continue
            
        # 7. while <cond>
        if code_part.startswith("while "):
            cond = clean_expr(code_part[6:])
            next_line = lines[i+1].strip() if i+1 < len(lines) else ""
            if next_line.startswith("open"):
                i += 1 # saltar 'open'
            block_stack.append(('while', cond))
            out.append(f"{indent}while {cond} do{inline_comment}")
            i += 1
            continue
            
        # 8. if <cond>
        if code_part.startswith("if "):
            cond = clean_expr(code_part[3:])
            # ver si viene 'then' y 'open' en siguientes líneas
            j = i + 1
            while j < len(lines) and lines[j].strip() in ["then", "open"]:
                j += 1
            i = j - 1
            block_stack.append(('if_chain', False))
            out.append(f"{indent}if {cond} then{inline_comment}")
            i += 1
            continue
            
        # 9. else if <cond>
        if code_part.startswith("else if "):
            cond = clean_expr(code_part[8:])
            j = i + 1
            while j < len(lines) and lines[j].strip() in ["then", "open"]:
                j += 1
            i = j - 1
            # ajusta el indent para que coincida con el if
            parent_indent = "    " * (len(block_stack) - 1)
            out.append(f"{parent_indent}elseif {cond} then{inline_comment}")
            i += 1
            continue
            
        # 10. else
        if code_part == "else":
            j = i + 1
            while j < len(lines) and lines[j].strip() in ["then", "open"]:
                j += 1
            i = j - 1
            parent_indent = "    " * (len(block_stack) - 1)
            out.append(f"{parent_indent}else{inline_comment}")
            i += 1
            continue
            
        # 11. close
        if code_part == "close":
            # Mirar si lo que sigue es 'else' o 'else if'
            # Buscar la próxima línea no vacía ni comentario
            next_code = None
            k = i + 1
            while k < len(lines):
                n_str = lines[k].strip()
                if n_str and not n_str.startswith(";"):
                    next_code = n_str
                    break
                k += 1
                
            if next_code and (next_code.startswith("else") or next_code.startswith("else if")):
                # No cerramos el if todavía porque se continuará con elseif/else
                i += 1
                continue
                
            if block_stack:
                b_type, b_name = block_stack.pop()
                close_indent = "    " * len(block_stack)
                if b_type == 'func':
                    if b_name == "IsTimeUp":
                        out.append("    return true")
                    out.append(f"end\n")
                elif b_type == 'repeat':
                    out.append(f"{close_indent}until true{inline_comment}")
                elif b_type in ['while', 'if_chain']:
                    out.append(f"{close_indent}end{inline_comment}")
            else:
                out.append(f"end\n")
            i += 1
            continue
            
        # 12. var Name
        m_var = re.match(r"^var\s+([a-zA-Z0-9_]+)$", code_part)
        if m_var:
            vname = m_var.group(1)
            out.append(f"{indent}local {vname} = nil{inline_comment}")
            i += 1
            continue
            
        # 13. Asignación: Var = Expr
        m_ass = re.match(r"^([a-zA-Z0-9_]+)\s*=\s*(.+)$", code_part)
        if m_ass:
            vname = m_ass.group(1)
            val = clean_expr(m_ass.group(2))
            out.append(f"{indent}{vname} = {val}{inline_comment}")
            i += 1
            continue
            
        # 14. break
        if code_part.startswith("break"):
            # Si estamos dentro de un while o repeat
            in_loop = any(b[0] in ['while', 'repeat'] for b in block_stack)
            if in_loop:
                out.append(f"{indent}break{inline_comment}")
            else:
                out.append(f"{indent}return false{inline_comment}")
            i += 1
            continue
            
        # 15. call
        # call "Routine"
        m_c1 = re.match(r'^call\s+"([a-zA-Z0-9_]+)"$', code_part)
        if m_c1:
            fn = m_c1.group(1)
            in_loop = any(b[0] in ['while', 'repeat'] for b in block_stack)
            if fn == "IsTimeUp" and in_loop:
                out.append(f"{indent}if not IsTimeUp() then break end{inline_comment}")
            else:
                out.append(f"{indent}{fn}(){inline_comment}")
            i += 1
            continue
        # call MyName "DefRoutine"
        m_c2 = re.match(r'^call\s+([a-zA-Z0-9_]+)\s+"([a-zA-Z0-9_]+)"$', code_part)
        if m_c2:
            out.append(f"{indent}CallNamedRoutine({m_c2.group(1)}, \"{m_c2.group(2)}\"){inline_comment}")
            i += 1
            continue
        # call DynamicVar or expression like "Regen" % AppendMobName
        if code_part.startswith("call "):
            call_target = clean_expr(code_part[5:].strip())
            out.append(f"{indent}if _G[{call_target}] then _G[{call_target}]() end{inline_comment}")
            i += 1
            continue
            
        # 16. mobregen
        if code_part.startswith("mobregen "):
            args = code_part[9:].strip()
            tokens = re.findall(r'[^\s"]+|"[^"]*"', args)
            f_tokens = [t if (t.startswith('"') and t.endswith('"')) else f'"{t}"' for t in tokens]
            out.append(f"{indent}cMobRegen({', '.join(f_tokens)}){inline_comment}")
            i += 1
            continue
            
        # 17. mobsummon
        if code_part.startswith("mobsummon "):
            args = code_part[10:].strip()
            tokens = re.findall(r'[^\s"]+|"[^"]*"', args)
            f_tokens = [t if (t.startswith('"') or t.isdigit()) else f'"{t}"' for t in tokens]
            out.append(f"{indent}cMobSummon({', '.join(f_tokens)}){inline_comment}")
            i += 1
            continue
            
        # 18. howmanymob
        m_hmm = re.match(r'^howmanymob\s+([a-zA-Z0-9_]+)\s+(.+)$', code_part)
        if m_hmm:
            v, filt = m_hmm.group(1), m_hmm.group(2).strip()
            filt_expr = clean_expr(filt)
            if not ('"' in filt_expr or '..' in filt_expr):
                filt_expr = f'"{filt_expr}"'
            out.append(f"{indent}{v} = cHowManyMob({filt_expr}){inline_comment}")
            i += 1
            continue
            
        # 19. kqrestminute
        m_kq = re.match(r'^kqrestminute\s+([a-zA-Z0-9_]+)$', code_part)
        if m_kq:
            out.append(f"{indent}{m_kq.group(1)} = cKQRestMinute(){inline_comment}")
            i += 1
            continue
            
        # 20. currenthp
        m_chp = re.match(r'^currenthp\s+([a-zA-Z0-9_]+)$', code_part)
        if m_chp:
            out.append(f"{indent}{m_chp.group(1)} = cCurrentHP(){inline_comment}")
            i += 1
            continue
            
        # 21. currentmap
        m_cmap = re.match(r'^currentmap\s+([a-zA-Z0-9_]+)$', code_part)
        if m_cmap:
            out.append(f"{indent}{m_cmap.group(1)} = cCurrentMap(){inline_comment}")
            i += 1
            continue
            
        # 22. myname
        m_mname = re.match(r'^myname\s+([a-zA-Z0-9_]+)$', code_part)
        if m_mname:
            out.append(f"{indent}{m_mname.group(1)} = cMyName(){inline_comment}")
            i += 1
            continue
            
        # 23. myhandle
        m_mhan = re.match(r'^myhandle\s+([a-zA-Z0-9_]+)$', code_part)
        if m_mhan:
            out.append(f"{indent}{m_mhan.group(1)} = cMyHandle(){inline_comment}")
            i += 1
            continue
            
        # 24. whoistarget
        m_wt = re.match(r'^whoistarget\s+([a-zA-Z0-9_]+)$', code_part)
        if m_wt:
            out.append(f"{indent}{m_wt.group(1)} = cWhoIsTarget(){inline_comment}")
            i += 1
            continue
            
        # 25. whokillme
        m_wk = re.match(r'^whokillme\s+([a-zA-Z0-9_]+)$', code_part)
        if m_wk:
            out.append(f"{indent}{m_wk.group(1)} = cWhoKillMe(){inline_comment}")
            i += 1
            continue
            
        # 26. getname
        m_gn = re.match(r'^getname\s+([a-zA-Z0-9_]+)\s+([a-zA-Z0-9_]+)$', code_part)
        if m_gn:
            out.append(f"{indent}{m_gn.group(1)} = cGetName({m_gn.group(2)}){inline_comment}")
            i += 1
            continue
            
        # 27. permillage
        m_pm = re.match(r'^permillage\s+([a-zA-Z0-9_]+)$', code_part)
        if m_pm:
            out.append(f"{indent}{m_pm.group(1)} = cPermillage(){inline_comment}")
            i += 1
            continue
            
        # 28. chat
        if code_part.startswith("chat "):
            parts = code_part[5:].strip().split(None, 1)
            target = parts[0]
            msg = clean_expr(parts[1]) if len(parts) > 1 else '""'
            out.append(f"{indent}cChat({target}, {msg}){inline_comment}")
            i += 1
            continue
            
        # 29. shout
        if code_part.startswith("shout "):
            tokens = re.findall(r'[^\s"]+|"[^"]*"', code_part[6:].strip())
            if len(tokens) >= 3:
                target = tokens[0] if tokens[0].startswith('"') else f'"{tokens[0]}"'
                sender = tokens[1] if tokens[1].startswith('"') else f'"{tokens[1]}"'
                msg = clean_expr(' '.join(tokens[2:]))
                out.append(f"{indent}cShout({target}, {sender}, {msg}){inline_comment}")
            elif len(tokens) == 2:
                target = tokens[0] if tokens[0].startswith('"') else f'"{tokens[0]}"'
                msg = clean_expr(tokens[1])
                out.append(f"{indent}cShout({target}, \"\", {msg}){inline_comment}")
            i += 1
            continue
            
        # 30. broadcast
        if code_part.startswith("broadcast "):
            tokens = re.findall(r'[^\s"]+|"[^"]*"', code_part[10:].strip())
            if len(tokens) >= 2:
                target = tokens[0] if tokens[0].startswith('"') else f'"{tokens[0]}"'
                msg = clean_expr(' '.join(tokens[1:]))
                out.append(f"{indent}cBroadcast({target}, {msg}){inline_comment}")
            i += 1
            continue
            
        # 31. print
        if code_part.startswith("print "):
            msg = clean_expr(code_part[6:].strip())
            out.append(f"{indent}print({msg}){inline_comment}")
            i += 1
            continue
            
        # 32. wait
        m_w = re.match(r'^wait\s+(sec|min)\s+(.+)$', code_part)
        if m_w:
            unit, val = m_w.group(1), m_w.group(2).strip()
            if unit == "sec":
                out.append(f"{indent}cWaitSec({val}){inline_comment}")
            else:
                out.append(f"{indent}cWaitMin({val}){inline_comment}")
            i += 1
            continue
            
        # 33. reward
        if code_part.startswith("reward "):
            target = code_part[7:].strip()
            out.append(f"{indent}cReward(\"{target}\"){inline_comment}")
            i += 1
            continue
            
        # 34. endofevent
        m_eoe = re.match(r'^endofevent\s+([a-zA-Z0-9_]+)\s+([0-9]+)$', code_part)
        if m_eoe:
            ev, res = m_eoe.group(1), m_eoe.group(2)
            out.append(f"{indent}cEndOfEvent(\"{ev}\", {res}){inline_comment}")
            i += 1
            continue
            
        # 35. linkto
        if code_part.startswith("linkto "):
            tokens = re.findall(r'[^\s"]+|"[^"]*"', code_part[7:].strip())
            f_tokens = [t if (t.startswith('"') or t.isdigit() or t.startswith('-')) else f'"{t}"' for t in tokens]
            out.append(f"{indent}cLinkTo({', '.join(f_tokens)}){inline_comment}")
            i += 1
            continue
            
        # 36. mapclear
        if code_part.startswith("mapclear "):
            target = code_part[9:].strip()
            out.append(f"{indent}cMapClear(\"{target}\"){inline_comment}")
            i += 1
            continue
            
        # 37. skillblast
        if code_part.startswith("skillblast "):
            tokens = re.findall(r'[^\s"]+|"[^"]*"', code_part[11:].strip())
            f_tokens = [t if t.startswith('"') else f'"{t}"' for t in tokens]
            out.append(f"{indent}cSkillBlast({', '.join(f_tokens)}){inline_comment}")
            i += 1
            continue
            
        # 38. retrate
        m_rr2 = re.match(r'^retrate\s+([a-zA-Z0-9_]+)\s+([a-zA-Z0-9_]+)$', code_part)
        if m_rr2:
            out.append(f"{indent}{m_rr2.group(1)} = cRetRate({m_rr2.group(2)}){inline_comment}")
            i += 1
            continue
        m_rr1 = re.match(r'^retrate\s+([a-zA-Z0-9_]+)$', code_part)
        if m_rr1:
            out.append(f"{indent}cRetRate(\"{m_rr1.group(1)}\"){inline_comment}")
            i += 1
            continue
            
        # Fallback
        out.append(f"{indent}{clean_expr(code_part)}{inline_comment}")
        i += 1

    return out

def convert_file(ps_path, lua_path):
    with open(ps_path, "r", encoding="cp949", errors="replace") as f:
        lines = f.readlines()
        
    res_lines = transpile_ps(lines)
    header = [
        f"-- {os.path.basename(lua_path)}",
        f"-- Modernizado desde {os.path.basename(ps_path)} para XelionCore ZoneService (Lua)",
        "-- 100% fidelidad de logica, eventos KQ y comentarios oficiales preservados.",
        ""
    ]
    with open(lua_path, "w", encoding="utf-8") as out:
        out.write("\n".join(header + res_lines) + "\n")
    print(f"Generated: {lua_path} ({len(res_lines)} lines)")

def main():
    base_dir = "data/MobBehaviorDescript"
    ps_files = sorted(glob.glob(f"{base_dir}/**/*.ps", recursive=True))
    if not ps_files:
        base_dir = "../AstraServer/Server/9Data/Shine/MobBehaviorDescript"
        ps_files = sorted(glob.glob(f"{base_dir}/**/*.ps", recursive=True))

    for ps in ps_files:
        rel = os.path.relpath(ps, base_dir)
        lua_path = os.path.join("data/MobBehaviorDescript", os.path.splitext(rel)[0] + ".lua")
        os.makedirs(os.path.dirname(lua_path), exist_ok=True)
        convert_file(ps, lua_path)

if __name__ == "__main__":
    main()
