// Extracts the data from transformmappings.html to machine-readable json. The json object consists of mappings from the
// base unicode character to all its different variants. For example, here is the set of mappings for the ascii "x" character:
// "x": {
//   "bold": {
//     "char": "𝐱",
//     "delta": "1D3B9"
//   },
//   "bold-fraktur": {
//     "char": "𝖝",
//     "delta": "1D525"
//   },
//   "bold-italic": {
//     "char": "𝒙",
//     "delta": "1D421"
//   },
//   "bold-sans-serif": {
//     "char": "𝘅",
//     "delta": "1D58D"
//   },
//   "bold-script": {
//     "char": "𝔁",
//     "delta": "1D489"
//   },
//   "double-struck": {
//     "char": "𝕩",
//     "delta": "1D4F1"
//   },
//   "fraktur": {
//     "char": "𝔵",
//     "delta": "1D4BD"
//   },
//   "italic": {
//     "char": "𝑥",
//     "delta": "1D3ED"
//   },
//   "monospace": {
//     "char": "𝚡",
//     "delta": "1D629"
//   },
//   "sans-serif": {
//     "char": "𝗑",
//     "delta": "1D559"
//   },
//   "sans-serif-bold-italic": {
//     "char": "𝙭",
//     "delta": "1D5F5"
//   },
//   "sans-serif-italic": {
//     "char": "𝘹",
//     "delta": "1D5C1"
//   },
//   "script": {
//     "char": "𝓍",
//     "delta": "1D455"
//   }
// }
//
// where "char" holds the actual unicode character and "delta" is the offset from the base codepoint.

function mapToJson(map){
  const obj = {};
  for (const [key, value] of map) {
    obj[key] = value instanceof Map ? mapToJson(value) : value;
  }
  return obj;
}

function extract() {
    let result = new Map();
    const variants = document.querySelectorAll("section section");
    variants.forEach(function(val){
        let variant = val.id.replaceAll("-mappings","");
        let table = val.querySelectorAll("tr");
        table.forEach(function(row){
            let cells = row.querySelectorAll("td");
            if (cells.length == 3) {
                let base = cells[0].innerHTML.split(" ")[0];
                let xform = cells[1].innerHTML.split(" ")[0];
                let delta = cells[2].innerHTML.split(" ")[0];
                if (!result.has(base)){
                    result.set(base, new Map());
                }
                result.get(base).set(variant, new Map([["char", xform], ["delta", delta]]));
            }
        });
    });
    return JSON.stringify(mapToJson(result));
};
