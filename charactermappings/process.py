#!/usr/bin/env python3

import xml.etree.ElementTree as ET
import json


def find_matching_brace(string, open_brace_index, braces):
    """Finds the index of the matching closing bracethesis.

    Args:
      string: The string to search.
      open_brace_index: The index of the opening bracethesis.

    Returns:
      The index of the matching closing bracethesis, or -1 if not found.
    """

    if string[open_brace_index] != "(":
        return -1

    count = 1
    for i in range(open_brace_index + 1, len(string)):
        if string[i] == braces[0]:
            count += 1
        elif string[i] == braces[1]:
            count -= 1
            if count == 0:
                return i

    return -1


def getEntityName(char: ET.Element):
    name = None
    for ent in char.findall("entity"):
        if name is None:
            name = ent.get("id")
            continue
        setname = ent.get("set")
        temp = ent.get("id")
        if setname == "mmlalias" and len(temp) < len(name):
            name = temp
        elif temp == name.lower():
            name = temp
    if name is not None:
        name = f"&{name};"

    return name


# One LaTeX 'character' may be multiple unicode characters!!!
def getCharacter(codepointStr: str):
    out = ""
    # remove the first 'U' character
    points = codepointStr[1:].split("-")
    for point in points:
        out += chr(int(point, 16))
    return out


def main():
    tree = ET.parse("unicode.xml")

    root = tree.getroot()
    with open("./additional_commands.json", 'r') as fp:
        symbols = json.load(fp)
    commands = {}
    command_arg_count = {}
    multi = {}

    for char in root.findall("character"):
        if char.get("mode") == "text":
            continue
        tex = char.find("AMS")
        if tex is None:
            tex = char.find("mathlatex")
        if tex is None:
            tex = char.find("latex")
        if tex is None:
            continue
        ctype = char.get("type")
        desc = char.find("description")
        codepoint = char.get("id")
        ent = getEntityName(char)
        sym = tex.text.strip()
        sym = sym.strip()
        if len(sym) > 0 and sym[0] != "\\":
            continue
        sym = tex.text.strip("\\").strip()
        if sym not in multi:
            multi[sym] = [
                {
                    "description": desc.text,
                    "codepoint": codepoint,
                    "entity": ent,
                    "char": getCharacter(codepoint),
                    "type": ctype,
                }
            ]
        else:
            multi[sym].append(
                {
                    "description": desc.text,
                    "codepoint": codepoint,
                    "entity": ent,
                    "char": getCharacter(codepoint),
                    "type": ctype,
                }
            )
        if "{" in sym and sym not in commands:
            commands[sym] = {
                "description": desc.text,
                "codepoint": codepoint,
                "entity": ent,
                "char": getCharacter(codepoint),
                "type": ctype,
            }
            argcount = sym.count("{")
            args = sym.split("{")
            if args[0] == "":
                continue
            if argcount not in command_arg_count:
                command_arg_count[argcount] = [args[0]]
            elif args[0] not in command_arg_count[argcount]:
                command_arg_count[argcount].append(args[0])

        elif sym not in symbols:
            symbols[sym] = {
                "description": desc.text,
                "codepoint": codepoint,
                "entity": ent,
                "char": getCharacter(codepoint),
                "type": ctype,
            }

    singular = []
    for sym in multi:
        if len(multi[sym]) <= 1:
            singular.append(sym)
    for sym in singular:
        del multi[sym]
    with open("commands.json", "w", encoding="utf-8") as fp:
        json.dump(commands, fp, sort_keys=True)
    with open("symbols.json", "w", encoding="utf-8") as fp:
        json.dump(symbols, fp, sort_keys=True)
    with open("counts.json", "w", encoding="utf-8") as fp:
        json.dump(command_arg_count, fp, sort_keys=True)
    with open("multi.json", "w", encoding="utf-8") as fp:
        json.dump(multi, fp, sort_keys=True)


def font_modifiers():
    tree = ET.parse("unicode.xml")

    root = tree.getroot()

    fonts = {}

    for char in root.findall("character"):
        if char.get("mode") == "text":
            continue
        tex = char.find("AMS")
        if tex is None:
            tex = char.find("mathlatex")
        if tex is None:
            tex = char.find("latex")
        if tex is None:
            continue
        ctype = char.get("type")
        desc = char.find("description")
        codepoint = char.get("id")
        ent = getEntityName(char)
        sym = tex.text.strip()
        if len(sym) > 0 and sym[0] != "\\":
            continue
        sym = tex.text.strip("\\")
        if len(sym) < len("math") or sym.find("math") != 0:
            continue
        arg_start = sym.find("{")
        if arg_start < 0:
            continue
        arg_end = find_matching_brace(sym, arg_start, "{}")
        arg = sym[arg_start + 1 : arg_end]
        arg = arg.strip().strip("\\")
        sym = sym[:arg_start]
        if sym not in fonts:
            fonts[sym] = {
                arg: getCharacter(codepoint),
            }
        else:
            fonts[sym][arg] = getCharacter(codepoint)
    with open("fonts.json", "w", encoding="utf-8") as fp:
        json.dump(fonts, fp, sort_keys=True)


def negations():
    tree = ET.parse("unicode.xml")

    root = tree.getroot()

    negs = {}

    for char in root.findall("character"):
        if char.get("mode") == "text":
            continue
        tex = char.find("latex")
        if tex is None:
            continue
        codepoint = char.get("id")
        sym = tex.text.strip()
        if len(sym) > 0 and sym[0] != "\\":
            continue
        sym = tex.text.strip().strip("\\")
        if len(sym) < len("not") or sym.find("not") != 0:
            continue
        print(sym)
        arg = sym[3:]
        sym = sym[:3]
        arg = arg.strip().strip("\\")
        if sym not in negs:
            negs[sym] = {
                arg: getCharacter(codepoint),
            }
        else:
            negs[sym][arg] = getCharacter(codepoint)
    with open("negations.json", "w", encoding="utf-8") as fp:
        json.dump(negs, fp, sort_keys=True)


if __name__ == "__main__":
    main()
    # font_modifiers()
    negations()
