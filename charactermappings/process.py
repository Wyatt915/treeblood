#!/usr/bin/env python3

import xml.etree.ElementTree as ET
import json


def getEntityName(char: ET.Element):
    name = None
    for ent in char.findall('entity'):
        if name is None:
            name = ent.get('id')
            continue
        setname = ent.get('set')
        temp = ent.get('id')
        if setname == 'mmlalias' and len(temp) < len(name):
            name = temp
        elif temp == name.lower():
            name = temp
    if name is not None:
        name = f'&{name};'


    return name


# One LaTeX 'character' may be multiple unicode characters!!!
def getCharacter(codepointStr: str):
    out = ""
    #remove the first 'U' character
    points =  codepointStr[1:].split('-')
    for point in points:
        out += chr(int(point, 16))
    return out

def main():
    tree = ET.parse('unicode.xml')

    root = tree.getroot()

    symbols = {}
    commands = {}
    command_arg_count = {}
    multi = {}

    for char in root.findall('character'):
        if char.get('mode') == 'text':
            continue
        tex = char.find('AMS')
        if tex is None:
            tex = char.find('mathlatex')
        if tex is None:
            tex = char.find('latex')
        if tex is None:
            continue
        ctype = char.get('type')
        desc = char.find('description')
        codepoint = char.get('id')
        ent = getEntityName(char)
        sym = tex.text.strip(' \\')
        if sym not in multi:
            multi[sym] = [{
                'description': desc.text,
                'codepoint': codepoint,
                'entity': ent,
                'char': getCharacter(codepoint),
                'type': ctype
                }]
        else:
            multi[sym].append({
                'description': desc.text,
                'codepoint': codepoint,
                'entity': ent,
                'char': getCharacter(codepoint),
                'type': ctype
                })
        if '{' in sym and sym not in commands:
            commands[sym] = {
                    'description': desc.text,
                    'codepoint': codepoint,
                    'entity': ent,
                    'char': getCharacter(codepoint),
                    'type': ctype
                    }
            argcount = sym.count('{')
            args = sym.split('{')
            if args[0] == "":
                continue
            if argcount not in command_arg_count:
                command_arg_count[argcount] = [args[0]]
            elif args[0] not in command_arg_count[argcount]:
                command_arg_count[argcount].append(args[0])

        elif sym not in symbols:
            symbols[sym] = {
                    'description': desc.text,
                    'codepoint': codepoint,
                    'entity': ent,
                    'char': getCharacter(codepoint),
                    'type': ctype
                    }

    singular = []
    for sym in multi:
        if len(multi[sym]) <= 1:
            singular.append(sym)
    for sym in singular:
        del multi[sym]
    with open('commands.json', 'w', encoding='utf-8') as fp:
        json.dump(commands, fp, sort_keys=True, indent=4)
    with open('symbols.json', 'w', encoding='utf-8') as fp:
        json.dump(symbols, fp, sort_keys=True, indent=4)
    with open('counts.json', 'w', encoding='utf-8') as fp:
        json.dump(command_arg_count, fp, sort_keys=True, indent=4)
    with open('multi.json', 'w', encoding='utf-8') as fp:
        json.dump(multi, fp, sort_keys=True, indent=4)


if __name__ == "__main__":
    main()
