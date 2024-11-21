package parse

type symbolKind int

const (
	symNormal symbolKind = iota
	sym_alphabetic
	sym_binaryop
	sym_other
	sym_relation
	sym_normal
	sym_opening
	sym_closing
	sym_diacritic
	sym_large
)

type symbol struct {
	char   string
	entity string
	kind   symbolKind
}

var symbolTable = map[string]symbol{
	"$": {
		char:   "$",
		entity: "&dollar;",
		kind:   sym_normal,
	},
	"-": {
		char:   "\u00ad",
		entity: "&shy;",
		kind:   sym_other,
	},
	"Alpha": {
		char:   "Α",
		entity: "&Alpha;",
		kind:   sym_alphabetic,
	},
	"Angle": {
		char:   "⦜",
		entity: "&vangrt;",
		kind:   sym_other,
	},
	"BbbPi": {
		char:   "ℿ",
		entity: "&opfpi;",
		kind:   sym_alphabetic,
	},
	"Beta": {
		char:   "Β",
		entity: "&Bgr;",
		kind:   sym_alphabetic,
	},
	"Bumpeq": {
		char:   "≎",
		entity: "&bump;",
		kind:   sym_relation,
	},
	"Cap": {
		char:   "⋒",
		entity: "&Cap;",
		kind:   sym_binaryop,
	},
	"Chi": {
		char:   "Χ",
		entity: "&KHgr;",
		kind:   sym_alphabetic,
	},
	"Colon": {
		char:   "∷",
		entity: "&Colon;",
		kind:   sym_other,
	},
	"Cup": {
		char:   "⋓",
		entity: "&Cup;",
		kind:   sym_binaryop,
	},
	"Dashv": {
		char:   "⫤",
		entity: "&Dashv;",
		kind:   sym_relation,
	},
	"Ddownarrow": {
		char:   "⤋",
		entity: "&dAarr;",
		kind:   sym_relation,
	},
	"Delta": {
		char:   "Δ",
		entity: "&Delta;",
		kind:   sym_alphabetic,
	},
	"Digamma": {
		char:   "Ϝ",
		entity: "&Gammad;",
		kind:   sym_alphabetic,
	},
	"Doteq": {
		char:   "≑",
		entity: "&eDot;",
		kind:   sym_relation,
	},
	"DownArrowBar": {
		char:   "⤓",
		entity: "&darrb;",
		kind:   sym_relation,
	},
	"DownArrowUpArrow": {
		char:   "⇵",
		entity: "&duarr;",
		kind:   sym_relation,
	},
	"DownLeftRightVector": {
		char:   "⥐",
		entity: "&ldrdshar;",
		kind:   sym_relation,
	},
	"DownLeftTeeVector": {
		char:   "⥞",
		entity: "&bldhar;",
		kind:   sym_relation,
	},
	"DownLeftVectorBar": {
		char:   "⥖",
		entity: "&ldharb;",
		kind:   sym_relation,
	},
	"DownRightTeeVector": {
		char:   "⥟",
		entity: "&brdhar;",
		kind:   sym_relation,
	},
	"DownRightVectorBar": {
		char:   "⥗",
		entity: "&rdharb;",
		kind:   sym_relation,
	},
	"Downarrow": {
		char:   "⇓",
		entity: "&dArr;",
		kind:   sym_relation,
	},
	"ElOr": {
		char:   "⩖",
		entity: "&oror;",
		kind:   sym_binaryop,
	},
	"Elroang": {
		char:   "⦆",
		entity: "&ropar;",
		kind:   sym_closing,
	},
	"ElsevierGlyph{225A}": {
		char:   "⩣",
		entity: "&veeBar;",
		kind:   sym_binaryop,
	},
	"ElsevierGlyph{3018}": {
		char:   "〘",
		entity: "&loang;",
		kind:   sym_opening,
	},
	"ElsevierGlyph{E25E}": {
		char:   "⨵",
		entity: "&rotimes;",
		kind:   sym_binaryop,
	},
	"ElzAnd": {
		char:   "⩓",
		entity: "&And;",
		kind:   sym_binaryop,
	},
	"ElzLap": {
		char:   "⧊",
		entity: "&tridoto;",
		kind:   sym_other,
	},
	"ElzOr": {
		char:   "⩔",
		entity: "&Or;",
		kind:   sym_binaryop,
	},
	"ElzRlarr": {
		char:   "⥂",
		entity: "&arrlrsl;",
		kind:   sym_relation,
	},
	"ElzTimes": {
		char:   "⨯",
		entity: "&htimes;",
		kind:   sym_binaryop,
	},
	"Elzbar": {
		char:   "̶",
		entity: "",
		kind:   sym_diacritic,
	},
	"Elzbtdl": {
		char:   "ɬ",
		entity: "",
		kind:   sym_other,
	},
	"Elzcirfb": {
		char:   "◒",
		entity: "",
		kind:   sym_other,
	},
	"Elzcirfl": {
		char:   "◐",
		entity: "",
		kind:   sym_other,
	},
	"Elzcirfr": {
		char:   "◑",
		entity: "",
		kind:   sym_other,
	},
	"Elzclomeg": {
		char:   "ɷ",
		entity: "",
		kind:   sym_other,
	},
	"Elzddfnc": {
		char:   "⦙",
		entity: "&vellip4;",
		kind:   sym_other,
	},
	"Elzdefas": {
		char:   "⧋",
		entity: "&tribar;",
		kind:   sym_other,
	},
	"Elzdlcorn": {
		char:   "⎣",
		entity: "&vldash;",
		kind:   sym_relation,
	},
	"Elzdshfnc": {
		char:   "┆",
		entity: "&Bvert;",
		kind:   sym_other,
	},
	"Elzdyogh": {
		char:   "ʤ",
		entity: "",
		kind:   sym_other,
	},
	"Elzesh": {
		char:   "ʃ",
		entity: "",
		kind:   sym_other,
	},
	"Elzfhr": {
		char:   "ɾ",
		entity: "",
		kind:   sym_other,
	},
	"Elzglst": {
		char:   "ʔ",
		entity: "",
		kind:   sym_other,
	},
	"Elzhlmrk": {
		char:   "ˑ",
		entity: "",
		kind:   sym_other,
	},
	"Elzinglst": {
		char:   "ʖ",
		entity: "",
		kind:   sym_other,
	},
	"Elzinvv": {
		char:   "ʌ",
		entity: "",
		kind:   sym_other,
	},
	"Elzinvw": {
		char:   "ʍ",
		entity: "",
		kind:   sym_other,
	},
	"Elzlmrk": {
		char:   "ː",
		entity: "",
		kind:   sym_other,
	},
	"Elzlow": {
		char:   "˕",
		entity: "",
		kind:   sym_other,
	},
	"Elzlpargt": {
		char:   "⦠",
		entity: "&gtrpar;",
		kind:   sym_other,
	},
	"Elzltlmr": {
		char:   "ɱ",
		entity: "",
		kind:   sym_other,
	},
	"Elzltln": {
		char:   "ɲ",
		entity: "",
		kind:   sym_other,
	},
	"Elzminhat": {
		char:   "⩟",
		entity: "&wedbar;",
		kind:   sym_binaryop,
	},
	"Elzopeno": {
		char:   "ɔ",
		entity: "",
		kind:   sym_other,
	},
	"Elzpalh": {
		char:   "̡",
		entity: "",
		kind:   sym_other,
	},
	"Elzpbgam": {
		char:   "ɤ",
		entity: "",
		kind:   sym_other,
	},
	"Elzpes": {
		char:   "₧",
		entity: "",
		kind:   sym_other,
	},
	"Elzpgamma": {
		char:   "ɣ",
		entity: "",
		kind:   sym_other,
	},
	"Elzpscrv": {
		char:   "ʋ",
		entity: "",
		kind:   sym_other,
	},
	"Elzpupsil": {
		char:   "ʊ",
		entity: "",
		kind:   sym_other,
	},
	"ElzrLarr": {
		char:   "⥄",
		entity: "&arrsrll;",
		kind:   sym_relation,
	},
	"Elzrais": {
		char:   "˔",
		entity: "",
		kind:   sym_other,
	},
	"Elzrarrx": {
		char:   "⥇",
		entity: "&rarrx;",
		kind:   sym_relation,
	},
	"Elzreapos": {
		char:   "‛",
		entity: "",
		kind:   sym_other,
	},
	"Elzreglst": {
		char:   "ʕ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrh": {
		char:   "̢",
		entity: "",
		kind:   sym_diacritic,
	},
	"Elzrl": {
		char:   "ɼ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtld": {
		char:   "ɖ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtll": {
		char:   "ɭ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtln": {
		char:   "ɳ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtlr": {
		char:   "ɽ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtls": {
		char:   "ʂ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtlt": {
		char:   "ʈ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrtlz": {
		char:   "ʐ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrttrnr": {
		char:   "ɻ",
		entity: "",
		kind:   sym_other,
	},
	"Elzrvbull": {
		char:   "◘",
		entity: "",
		kind:   sym_other,
	},
	"Elzsbbrg": {
		char:   "̪",
		entity: "",
		kind:   sym_other,
	},
	"Elzsblhr": {
		char:   "˓",
		entity: "",
		kind:   sym_other,
	},
	"Elzsbrhr": {
		char:   "˒",
		entity: "",
		kind:   sym_other,
	},
	"Elzschwa": {
		char:   "ə",
		entity: "",
		kind:   sym_other,
	},
	"Elzsqfl": {
		char:   "◧",
		entity: "",
		kind:   sym_other,
	},
	"Elzsqfnw": {
		char:   "┙",
		entity: "",
		kind:   sym_other,
	},
	"Elzsqfr": {
		char:   "◨",
		entity: "",
		kind:   sym_other,
	},
	"Elzsqfse": {
		char:   "◪",
		entity: "",
		kind:   sym_other,
	},
	"Elzsqspne": {
		char:   "⋥",
		entity: "",
		kind:   sym_other,
	},
	"Elztdcol": {
		char:   "⫶",
		entity: "&vellipv;",
		kind:   sym_other,
	},
	"Elztesh": {
		char:   "ʧ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrna": {
		char:   "ɐ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnh": {
		char:   "ɥ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnm": {
		char:   "ɯ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnmlr": {
		char:   "ɰ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnr": {
		char:   "ɹ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnrl": {
		char:   "ɺ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnsa": {
		char:   "ɒ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrnt": {
		char:   "ʇ",
		entity: "",
		kind:   sym_other,
	},
	"Elztrny": {
		char:   "ʎ",
		entity: "",
		kind:   sym_other,
	},
	"Elzverti": {
		char:   "ˌ",
		entity: "",
		kind:   sym_other,
	},
	"Elzverts": {
		char:   "ˈ",
		entity: "",
		kind:   sym_other,
	},
	"Elzvrecto": {
		char:   "▯",
		entity: "",
		kind:   sym_other,
	},
	"Elzxh": {
		char:   "ħ",
		entity: "&hstrok;",
		kind:   sym_alphabetic,
	},
	"Elzxl": {
		char:   "̵",
		entity: "",
		kind:   sym_diacritic,
	},
	"Elzxrat": {
		char:   "℞",
		entity: "&rx;",
		kind:   sym_normal,
	},
	"Elzyogh": {
		char:   "ʒ",
		entity: "",
		kind:   sym_other,
	},
	"Epsilon": {
		char:   "Ε",
		entity: "&Egr;",
		kind:   sym_alphabetic,
	},
	"Equal": {
		char:   "⩵",
		entity: "&eqeq;",
		kind:   sym_relation,
	},
	"Eta": {
		char:   "Η",
		entity: "&EEgr;",
		kind:   sym_alphabetic,
	},
	"Game": {
		char:   "⅁",
		entity: "&Game;",
		kind:   sym_normal,
	},
	"Gamma": {
		char:   "Γ",
		entity: "&Gamma;",
		kind:   sym_alphabetic,
	},
	"Im": {
		char:   "ℑ",
		entity: "&Im;",
		kind:   sym_alphabetic,
	},
	"Iota": {
		char:   "Ι",
		entity: "&Igr;",
		kind:   sym_alphabetic,
	},
	"Kappa": {
		char:   "Κ",
		entity: "&Kgr;",
		kind:   sym_alphabetic,
	},
	"Koppa": {
		char:   "Ϟ",
		entity: "&koppa;",
		kind:   sym_alphabetic,
	},
	"Lambda": {
		char:   "Λ",
		entity: "&Lambda;",
		kind:   sym_alphabetic,
	},
	"LeftDownTeeVector": {
		char:   "⥡",
		entity: "&bdlhar;",
		kind:   sym_relation,
	},
	"LeftDownVectorBar": {
		char:   "⥙",
		entity: "&dlharb;",
		kind:   sym_relation,
	},
	"LeftRightVector": {
		char:   "⥎",
		entity: "&lurushar;",
		kind:   sym_relation,
	},
	"LeftTeeVector": {
		char:   "⥚",
		entity: "&bluhar;",
		kind:   sym_relation,
	},
	"LeftTriangleBar": {
		char:   "⧏",
		entity: "&ltrivb;",
		kind:   sym_other,
	},
	"LeftUpDownVector": {
		char:   "⥑",
		entity: "&uldlshar;",
		kind:   sym_relation,
	},
	"LeftUpTeeVector": {
		char:   "⥠",
		entity: "&bulhar;",
		kind:   sym_relation,
	},
	"LeftUpVectorBar": {
		char:   "⥘",
		entity: "&ulharb;",
		kind:   sym_relation,
	},
	"LeftVectorBar": {
		char:   "⥒",
		entity: "&luharb;",
		kind:   sym_relation,
	},
	"Leftarrow": {
		char:   "⇐",
		entity: "&lArr;",
		kind:   sym_relation,
	},
	"Leftrightarrow": {
		char:   "⇔",
		entity: "&hArr;",
		kind:   sym_relation,
	},
	"Lleftarrow": {
		char:   "⇚",
		entity: "&lAarr;",
		kind:   sym_relation,
	},
	"Longleftarrow": {
		char:   "⟸",
		entity: "&xlArr;",
		kind:   sym_relation,
	},
	"Longleftrightarrow": {
		char:   "⟺",
		entity: "&xhArr;",
		kind:   sym_relation,
	},
	"Longrightarrow": {
		char:   "⟹",
		entity: "&xrArr;",
		kind:   sym_relation,
	},
	"Lsh": {
		char:   "↰",
		entity: "&lsh;",
		kind:   sym_relation,
	},
	"Mapsfrom": {
		char:   "⤆",
		entity: "&Mapfrom;",
		kind:   sym_relation,
	},
	"Mapsto": {
		char:   "⤇",
		entity: "&Mapto;",
		kind:   sym_relation,
	},
	"NestedGreaterGreater": {
		char:   "⪢",
		entity: "&Gt;",
		kind:   sym_relation,
	},
	"NestedLessLess": {
		char:   "⪡",
		entity: "&Lt;",
		kind:   sym_relation,
	},
	"NotGreaterGreater": {
		char:   "≫̸",
		entity: "&nGtv;",
		kind:   sym_relation,
	},
	"NotLeftTriangleBar": {
		char:   "⧏̸",
		entity: "&nltrivb;",
		kind:   sym_other,
	},
	"NotLessLess": {
		char:   "≪̸",
		entity: "&nLtv;",
		kind:   sym_relation,
	},
	"NotNestedGreaterGreater": {
		char:   "⪢̸",
		entity: "&nsGt;",
		kind:   sym_relation,
	},
	"NotNestedLessLess": {
		char:   "⪡̸",
		entity: "&nsLt;",
		kind:   sym_relation,
	},
	"NotRightTriangleBar": {
		char:   "⧐̸",
		entity: "&nvbrtri;",
		kind:   sym_other,
	},
	"NotSquareSubset": {
		char:   "⊏̸",
		entity: "&nsqsub;",
		kind:   sym_relation,
	},
	"NotSquareSuperset": {
		char:   "⊐̸",
		entity: "&nsqsup;",
		kind:   sym_relation,
	},
	"Omega": {
		char:   "Ω",
		entity: "&OHgr;",
		kind:   sym_alphabetic,
	},
	"P": {
		char:   "¶",
		entity: "&para;",
		kind:   sym_normal,
	},
	"Phi": {
		char:   "Φ",
		entity: "&PHgr;",
		kind:   sym_alphabetic,
	},
	"Pi": {
		char:   "Π",
		entity: "&Pgr;",
		kind:   sym_alphabetic,
	},
	"Psi": {
		char:   "Ψ",
		entity: "&PSgr;",
		kind:   sym_alphabetic,
	},
	"Re": {
		char:   "ℜ",
		entity: "&Re;",
		kind:   sym_alphabetic,
	},
	"ReverseUpEquilibrium": {
		char:   "⥯",
		entity: "&duhar;",
		kind:   sym_relation,
	},
	"Rho": {
		char:   "Ρ",
		entity: "&Rgr;",
		kind:   sym_alphabetic,
	},
	"RightDownTeeVector": {
		char:   "⥝",
		entity: "&bdrhar;",
		kind:   sym_relation,
	},
	"RightDownVectorBar": {
		char:   "⥕",
		entity: "&drharb;",
		kind:   sym_relation,
	},
	"RightTeeVector": {
		char:   "⥛",
		entity: "&bruhar;",
		kind:   sym_relation,
	},
	"RightTriangleBar": {
		char:   "⧐",
		entity: "&vbrtri;",
		kind:   sym_other,
	},
	"RightUpDownVector": {
		char:   "⥏",
		entity: "&urdrshar;",
		kind:   sym_relation,
	},
	"RightUpTeeVector": {
		char:   "⥜",
		entity: "&burhar;",
		kind:   sym_relation,
	},
	"RightUpVectorBar": {
		char:   "⥔",
		entity: "&urharb;",
		kind:   sym_relation,
	},
	"RightVectorBar": {
		char:   "⥓",
		entity: "&ruharb;",
		kind:   sym_relation,
	},
	"Rightarrow": {
		char:   "⇒",
		entity: "&rArr;",
		kind:   sym_relation,
	},
	"RoundImplies": {
		char:   "⥰",
		entity: "&rimply;",
		kind:   sym_relation,
	},
	"Rrightarrow": {
		char:   "⇛",
		entity: "&rAarr;",
		kind:   sym_relation,
	},
	"Rsh": {
		char:   "↱",
		entity: "&rsh;",
		kind:   sym_relation,
	},
	"RuleDelayed": {
		char:   "⧴",
		entity: "&;",
		kind:   sym_other,
	},
	"S": {
		char:   "§",
		entity: "&sect;",
		kind:   sym_normal,
	},
	"Sampi": {
		char:   "Ϡ",
		entity: "&sampi;",
		kind:   sym_alphabetic,
	},
	"Sigma": {
		char:   "Σ",
		entity: "&Sgr;",
		kind:   sym_alphabetic,
	},
	"Stigma": {
		char:   "Ϛ",
		entity: "&stigma;",
		kind:   sym_alphabetic,
	},
	"Subset": {
		char:   "⋐",
		entity: "&Sub;",
		kind:   sym_relation,
	},
	"Supset": {
		char:   "⋑",
		entity: "&Sup;",
		kind:   sym_relation,
	},
	"Tau": {
		char:   "Τ",
		entity: "&Tgr;",
		kind:   sym_alphabetic,
	},
	"Theta": {
		char:   "Θ",
		entity: "&THgr;",
		kind:   sym_alphabetic,
	},
	"UpArrowBar": {
		char:   "⤒",
		entity: "&uarrb;",
		kind:   sym_relation,
	},
	"UpEquilibrium": {
		char:   "⥮",
		entity: "&udhar;",
		kind:   sym_relation,
	},
	"Uparrow": {
		char:   "⇑",
		entity: "&uArr;",
		kind:   sym_relation,
	},
	"Updownarrow": {
		char:   "⇕",
		entity: "&vArr;",
		kind:   sym_relation,
	},
	"Upsilon": {
		char:   "Υ",
		entity: "&Ugr;",
		kind:   sym_alphabetic,
	},
	"Uuparrow": {
		char:   "⤊",
		entity: "&uAarr;",
		kind:   sym_relation,
	},
	"VDash": {
		char:   "⊫",
		entity: "&VDash;",
		kind:   sym_relation,
	},
	"Vdash": {
		char:   "⊩",
		entity: "&Vdash;",
		kind:   sym_relation,
	},
	"Vert": {
		char:   "‖",
		entity: "&Vert;",
		kind:   sym_other,
	},
	"Vvdash": {
		char:   "⊪",
		entity: "&Vvdash;",
		kind:   sym_relation,
	},
	"Vvert": {
		char:   "⦀",
		entity: "&tverbar;",
		kind:   sym_other,
	},
	"Xi": {
		char:   "Ξ",
		entity: "&Xgr;",
		kind:   sym_alphabetic,
	},
	"Zeta": {
		char:   "Ζ",
		entity: "&Zgr;",
		kind:   sym_alphabetic,
	},
	"_": {
		char:   "_",
		entity: "&lowbar;",
		kind:   sym_other,
	},
	"acute": {
		char:   "́",
		entity: "",
		kind:   sym_diacritic,
	},
	"adots": {
		char:   "⋰",
		entity: "&utdot;",
		kind:   sym_other,
	},
	"aleph": {
		char:   "ℵ",
		entity: "&alefsym;",
		kind:   sym_alphabetic,
	},
	"allequal": {
		char:   "≌",
		entity: "&bcong;",
		kind:   sym_other,
	},
	"alpha": {
		char:   "α",
		entity: "&agr;",
		kind:   sym_alphabetic,
	},
	"amalg": {
		char:   "⨿",
		entity: "&amalg;",
		kind:   sym_binaryop,
	},
	"angle": {
		char:   "∠",
		entity: "&ang;",
		kind:   sym_normal,
	},
	"approx": {
		char:   "≈",
		entity: "&asymp;",
		kind:   sym_relation,
	},
	"approxeq": {
		char:   "≊",
		entity: "&ape;",
		kind:   sym_relation,
	},
	"approxnotequal": {
		char:   "≆",
		entity: "&simne;",
		kind:   sym_relation,
	},
	"ast": {
		char:   "*",
		entity: "&ast;",
		kind:   sym_other,
	},
	"asymp": {
		char:   "≍",
		entity: "&CupCap;",
		kind:   sym_relation,
	},
	"backepsilon": {
		char:   "϶",
		entity: "&bepsi;",
		kind:   sym_other,
	},
	"backprime": {
		char:   "‵",
		entity: "&bprime;",
		kind:   sym_other,
	},
	"backsim": {
		char:   "∽",
		entity: "&bsim;",
		kind:   sym_relation,
	},
	"backsimeq": {
		char:   "⋍",
		entity: "&bsime;",
		kind:   sym_relation,
	},
	"backslash": {
		char:   "\\",
		entity: "&bsol;",
		kind:   sym_normal,
	},
	"bar": {
		char:   "̄",
		entity: "",
		kind:   sym_diacritic,
	},
	"barwedge": {
		char:   "⌅",
		entity: "&barwed;",
		kind:   sym_other,
	},
	"bbsum": {
		char:   "⅀",
		entity: "&opfsum;",
		kind:   sym_large,
	},
	"because": {
		char:   "∵",
		entity: "&becaus;",
		kind:   sym_normal,
	},
	"beta": {
		char:   "β",
		entity: "&beta;",
		kind:   sym_alphabetic,
	},
	"beth": {
		char:   "ℶ",
		entity: "&beth;",
		kind:   sym_alphabetic,
	},
	"between": {
		char:   "≬",
		entity: "&twixt;",
		kind:   sym_relation,
	},
	"bigcap": {
		char:   "⋂",
		entity: "&xcap;",
		kind:   sym_large,
	},
	"bigcirc": {
		char:   "○",
		entity: "&cir;",
		kind:   sym_binaryop,
	},
	"bigcup": {
		char:   "⋃",
		entity: "&xcup;",
		kind:   sym_large,
	},
	"bigcupdot": {
		char:   "⨃",
		entity: "&xcupdot;",
		kind:   sym_large,
	},
	"bigodot": {
		char:   "⨀",
		entity: "&xodot;",
		kind:   sym_large,
	},
	"bigoplus": {
		char:   "⨁",
		entity: "&xoplus;",
		kind:   sym_large,
	},
	"bigotimes": {
		char:   "⨂",
		entity: "&xotime;",
		kind:   sym_large,
	},
	"bigsqcap": {
		char:   "⨅",
		entity: "&xsqcap;",
		kind:   sym_large,
	},
	"bigsqcup": {
		char:   "⨆",
		entity: "&xsqcup;",
		kind:   sym_large,
	},
	"bigtimes": {
		char:   "⨉",
		entity: "&xtimes;",
		kind:   sym_large,
	},
	"bigtriangledown": {
		char:   "▽",
		entity: "&xdtri;",
		kind:   sym_other,
	},
	"bigtriangleup": {
		char:   "△",
		entity: "&xutri;",
		kind:   sym_other,
	},
	"biguplus": {
		char:   "⨄",
		entity: "&xuplus;",
		kind:   sym_large,
	},
	"bigvee": {
		char:   "⋁",
		entity: "&Vee;",
		kind:   sym_large,
	},
	"bigwedge": {
		char:   "⋀",
		entity: "&Wedge;",
		kind:   sym_large,
	},
	"bkarow": {
		char:   "⤍",
		entity: "&rbarr;",
		kind:   sym_other,
	},
	"blacklozenge": {
		char:   "⧫",
		entity: "&lozf;",
		kind:   sym_other,
	},
	"blacksquare": {
		char:   "▪",
		entity: "&squf;",
		kind:   sym_other,
	},
	"blacktriangle": {
		char:   "▴",
		entity: "&utrif;",
		kind:   sym_other,
	},
	"blacktriangledown": {
		char:   "▾",
		entity: "&dtrif;",
		kind:   sym_other,
	},
	"blacktriangleleft": {
		char:   "◂",
		entity: "&ltrif;",
		kind:   sym_other,
	},
	"blacktriangleright": {
		char:   "▸",
		entity: "&rtrif;",
		kind:   sym_other,
	},
	"bowtie": {
		char:   "⋈",
		entity: "&bowtie;",
		kind:   sym_relation,
	},
	"boxast": {
		char:   "⧆",
		entity: "&astb;",
		kind:   sym_other,
	},
	"boxbslash": {
		char:   "⧅",
		entity: "&bsolb;",
		kind:   sym_other,
	},
	"boxcircle": {
		char:   "⧇",
		entity: "&cirb;",
		kind:   sym_other,
	},
	"boxdiag": {
		char:   "⧄",
		entity: "&solb;",
		kind:   sym_other,
	},
	"boxdot": {
		char:   "⊡",
		entity: "&sdotb;",
		kind:   sym_binaryop,
	},
	"boxminus": {
		char:   "⊟",
		entity: "&minusb;",
		kind:   sym_binaryop,
	},
	"boxplus": {
		char:   "⊞",
		entity: "&plusb;",
		kind:   sym_binaryop,
	},
	"boxtimes": {
		char:   "⊠",
		entity: "&timesb;",
		kind:   sym_binaryop,
	},
	"breve": {
		char:   "̆",
		entity: "",
		kind:   sym_diacritic,
	},
	"btimes": {
		char:   "⨲",
		entity: "&btimes;",
		kind:   sym_binaryop,
	},
	"bullet": {
		char:   "•",
		entity: "&bull;",
		kind:   sym_binaryop,
	},
	"bumpeq": {
		char:   "≏",
		entity: "&bumpe;",
		kind:   sym_relation,
	},
	"bumpeqq": {
		char:   "⪮",
		entity: "&bumpE;",
		kind:   sym_relation,
	},
	"cap": {
		char:   "∩",
		entity: "&cap;",
		kind:   sym_binaryop,
	},
	"cdot": {
		char:   "⋅",
		entity: "&sdot;",
		kind:   sym_binaryop,
	},
	"cdotp": {
		char:   "·",
		entity: "&middot;",
		kind:   sym_binaryop,
	},
	"cdots": {
		char:   "⋯",
		entity: "&ctdot;",
		kind:   sym_other,
	},
	"check": {
		char:   "̌",
		entity: "",
		kind:   sym_diacritic,
	},
	"chi": {
		char:   "χ",
		entity: "&chi;",
		kind:   sym_alphabetic,
	},
	"circ": {
		char:   "∘",
		entity: "&compfn;",
		kind:   sym_binaryop,
	},
	"circeq": {
		char:   "≗",
		entity: "&cire;",
		kind:   sym_relation,
	},
	"circlearrowleft": {
		char:   "↺",
		entity: "&olarr;",
		kind:   sym_other,
	},
	"circlearrowright": {
		char:   "↻",
		entity: "&orarr;",
		kind:   sym_other,
	},
	"circledR": {
		char:   "®",
		entity: "&reg;",
		kind:   sym_normal,
	},
	"circledS": {
		char:   "Ⓢ",
		entity: "&oS;",
		kind:   sym_other,
	},
	"circledast": {
		char:   "⊛",
		entity: "&oast;",
		kind:   sym_binaryop,
	},
	"circledcirc": {
		char:   "⊚",
		entity: "&ocir;",
		kind:   sym_binaryop,
	},
	"circleddash": {
		char:   "⊝",
		entity: "&odash;",
		kind:   sym_binaryop,
	},
	"clockoint": {
		char:   "⨏",
		entity: "&slint;",
		kind:   sym_large,
	},
	"clwintegral": {
		char:   "∱",
		entity: "&cwint;",
		kind:   sym_large,
	},
	"complement": {
		char:   "∁",
		entity: "&comp;",
		kind:   sym_normal,
	},
	"cong": {
		char:   "≅",
		entity: "&cong;",
		kind:   sym_relation,
	},
	"conjquant": {
		char:   "⨇",
		entity: "&xandand;",
		kind:   sym_large,
	},
	"coprod": {
		char:   "∐",
		entity: "&coprod;",
		kind:   sym_large,
	},
	"copyright": {
		char:   "©",
		entity: "&copy;",
		kind:   sym_normal,
	},
	"cup": {
		char:   "∪",
		entity: "&cup;",
		kind:   sym_binaryop,
	},
	"cupdot": {
		char:   "⊍",
		entity: "&cupdot;",
		kind:   sym_binaryop,
	},
	"curlyeqprec": {
		char:   "⋞",
		entity: "&cuepr;",
		kind:   sym_relation,
	},
	"curlyeqsucc": {
		char:   "⋟",
		entity: "&cuesc;",
		kind:   sym_relation,
	},
	"curlyvee": {
		char:   "⋎",
		entity: "&cuvee;",
		kind:   sym_binaryop,
	},
	"curlywedge": {
		char:   "⋏",
		entity: "&cuwed;",
		kind:   sym_binaryop,
	},
	"curvearrowleft": {
		char:   "↶",
		entity: "&cularr;",
		kind:   sym_relation,
	},
	"curvearrowright": {
		char:   "↷",
		entity: "&curarr;",
		kind:   sym_relation,
	},
	"dagger": {
		char:   "†",
		entity: "&dagger;",
		kind:   sym_other,
	},
	"daleth": {
		char:   "ℸ",
		entity: "&daleth;",
		kind:   sym_alphabetic,
	},
	"dashV": {
		char:   "⫣",
		entity: "&dashV;",
		kind:   sym_relation,
	},
	"dashv": {
		char:   "⊣",
		entity: "&dashv;",
		kind:   sym_relation,
	},
	"dbkarow": {
		char:   "⤏",
		entity: "&rBarr;",
		kind:   sym_relation,
	},
	"dblarrowupdown": {
		char:   "⇅",
		entity: "&udarr;",
		kind:   sym_relation,
	},
	"ddagger": {
		char:   "‡",
		entity: "&Dagger;",
		kind:   sym_other,
	},
	"ddddot": {
		char:   "⃜",
		entity: "&DotDot;",
		kind:   sym_diacritic,
	},
	"dddot": {
		char:   "⃛",
		entity: "&tdot;",
		kind:   sym_diacritic,
	},
	"ddot": {
		char:   "̈",
		entity: "",
		kind:   sym_diacritic,
	},
	"ddots": {
		char:   "⋱",
		entity: "&dtdot;",
		kind:   sym_other,
	},
	"ddotseq": {
		char:   "⩷",
		entity: "&eDDot;",
		kind:   sym_relation,
	},
	"degree": {
		char:   "°",
		entity: "&deg;",
		kind:   sym_other,
	},
	"delta": {
		char:   "δ",
		entity: "&delta;",
		kind:   sym_alphabetic,
	},
	"diagdown": {
		char:   "╲",
		entity: "&xsol;",
		kind:   sym_other,
	},
	"diagup": {
		char:   "╱",
		entity: "&xbsol;",
		kind:   sym_other,
	},
	"diamond": {
		char:   "⋄",
		entity: "&diam;",
		kind:   sym_binaryop,
	},
	"diamondsuit": {
		char:   "♢",
		entity: "",
		kind:   sym_normal,
	},
	"digamma": {
		char:   "ϝ",
		entity: "&gammad;",
		kind:   sym_alphabetic,
	},
	"disjquant": {
		char:   "⨈",
		entity: "&xoror;",
		kind:   sym_large,
	},
	"div": {
		char:   "÷",
		entity: "&div;",
		kind:   sym_binaryop,
	},
	"divideontimes": {
		char:   "⋇",
		entity: "&divonx;",
		kind:   sym_binaryop,
	},
	"dot": {
		char:   "̇",
		entity: "",
		kind:   sym_diacritic,
	},
	"doteq": {
		char:   "≐",
		entity: "&esdot;",
		kind:   sym_relation,
	},
	"dotminus": {
		char:   "∸",
		entity: "&minusd;",
		kind:   sym_binaryop,
	},
	"dotplus": {
		char:   "∔",
		entity: "&plusdo;",
		kind:   sym_binaryop,
	},
	"dots": {
		char:   "…",
		entity: "&#x2026;",
		kind:   sym_other,
	},
	"doublebarwedge ?": {
		char:   "⌆",
		entity: "&Barwed;",
		kind:   sym_binaryop,
	},
	"downarrow": {
		char:   "↓",
		entity: "&darr;",
		kind:   sym_relation,
	},
	"downdownarrows": {
		char:   "⇊",
		entity: "&darr2;",
		kind:   sym_relation,
	},
	"downharpoonleft": {
		char:   "⇃",
		entity: "&dharl;",
		kind:   sym_relation,
	},
	"downharpoonright": {
		char:   "⇂",
		entity: "&dharr;",
		kind:   sym_relation,
	},
	"drbkarrow": {
		char:   "⤐",
		entity: "&RBarr;",
		kind:   sym_relation,
	},
	"dualmap": {
		char:   "⧟",
		entity: "&dumap;",
		kind:   sym_relation,
	},
	"ell": {
		char:   "ℓ",
		entity: "&ell;",
		kind:   sym_alphabetic,
	},
	"emdash": {
		char:   "—",
		entity: "&mdash;",
		kind:   sym_normal,
	},
	"epsilon": {
		char:   "ϵ",
		entity: "&epsi;",
		kind:   sym_alphabetic,
	},
	"eqcirc": {
		char:   "≖",
		entity: "&ecir;",
		kind:   sym_relation,
	},
	"eqcolon": {
		char:   "≕",
		entity: "&ecolon;",
		kind:   sym_relation,
	},
	"eqsim": {
		char:   "≂",
		entity: "&esim;",
		kind:   sym_relation,
	},
	"eqslantgtr": {
		char:   "⪖",
		entity: "&egs;",
		kind:   sym_relation,
	},
	"eqslantless": {
		char:   "⪕",
		entity: "&els;",
		kind:   sym_relation,
	},
	"equiv": {
		char:   "≡",
		entity: "&equiv;",
		kind:   sym_relation,
	},
	"eta": {
		char:   "η",
		entity: "&eegr;",
		kind:   sym_alphabetic,
	},
	"eth": {
		char:   "ƪ",
		entity: "",
		kind:   sym_other,
	},
	"exists": {
		char:   "∃",
		entity: "&exist;",
		kind:   sym_normal,
	},
	"fallingdotseq": {
		char:   "≒",
		entity: "&efDot;",
		kind:   sym_relation,
	},
	"fdiagovnearrow": {
		char:   "⤯",
		entity: "&fdonearr;",
		kind:   sym_other,
	},
	"fdiagovrdiag": {
		char:   "⤬",
		entity: "&fdiordi;",
		kind:   sym_other,
	},
	"flat": {
		char:   "♭",
		entity: "&flat;",
		kind:   sym_normal,
	},
	"forall": {
		char:   "∀",
		entity: "&forall;",
		kind:   sym_normal,
	},
	"forks": {
		char:   "⫝̸",
		entity: "&;",
		kind:   sym_relation,
	},
	"forksnot": {
		char:   "⫝",
		entity: "&;",
		kind:   sym_relation,
	},
	"frown": {
		char:   "⌢",
		entity: "&frown;",
		kind:   sym_relation,
	},
	"gamma": {
		char:   "γ",
		entity: "&gamma;",
		kind:   sym_alphabetic,
	},
	"ge": {
		char:   "≥",
		entity: "&ge;",
		kind:   sym_relation,
	},
	"geqq": {
		char:   "≧",
		entity: "&gE;",
		kind:   sym_relation,
	},
	"geqslant": {
		char:   "⩾",
		entity: "&ges;",
		kind:   sym_relation,
	},
	"gg": {
		char:   "≫",
		entity: "&Gt;",
		kind:   sym_relation,
	},
	"ggg": {
		char:   "⋙",
		entity: "&Gg;",
		kind:   sym_relation,
	},
	"gimel": {
		char:   "ℷ",
		entity: "&gimel;",
		kind:   sym_alphabetic,
	},
	"gnapprox": {
		char:   "⪊",
		entity: "&gnap;",
		kind:   sym_relation,
	},
	"gneq": {
		char:   "⪈",
		entity: "&gne;",
		kind:   sym_relation,
	},
	"gneqq": {
		char:   "≩",
		entity: "&gnE;",
		kind:   sym_relation,
	},
	"gnsim": {
		char:   "⋧",
		entity: "&gnsim;",
		kind:   sym_relation,
	},
	"grave": {
		char:   "̀",
		entity: "",
		kind:   sym_diacritic,
	},
	"greater": {
		char:   ">",
		entity: "&gt;",
		kind:   sym_relation,
	},
	"gtrapprox": {
		char:   "⪆",
		entity: "&gap;",
		kind:   sym_relation,
	},
	"gtrdot": {
		char:   "⋗",
		entity: "&gsdot;",
		kind:   sym_relation,
	},
	"gtreqless": {
		char:   "⋛",
		entity: "&gel;",
		kind:   sym_relation,
	},
	"gtreqqless": {
		char:   "⪌",
		entity: "&gEl;",
		kind:   sym_relation,
	},
	"gtrless": {
		char:   "≷",
		entity: "&gl;",
		kind:   sym_relation,
	},
	"gtrsim": {
		char:   "≳",
		entity: "&gsim;",
		kind:   sym_relation,
	},
	"guilsinglleft": {
		char:   "‹",
		entity: "&lsaquo;",
		kind:   sym_opening,
	},
	"guilsinglright": {
		char:   "›",
		entity: "&rsaquo;",
		kind:   sym_closing,
	},
	"gvertneqq": {
		char:   "≩︀",
		entity: "&gvnE;",
		kind:   sym_relation,
	},
	"hat": {
		char:   "̂",
		entity: "",
		kind:   sym_diacritic,
	},
	"heartsuit": {
		char:   "♡",
		entity: "",
		kind:   sym_normal,
	},
	"hermitconjmatrix": {
		char:   "⊹",
		entity: "&hercon;",
		kind:   sym_other,
	},
	"hksearow": {
		char:   "⤥",
		entity: "&searhk;",
		kind:   sym_relation,
	},
	"hkswarow": {
		char:   "⤦",
		entity: "&swarhk;",
		kind:   sym_relation,
	},
	"hookleftarrow": {
		char:   "↩",
		entity: "&larrhk;",
		kind:   sym_relation,
	},
	"hookrightarrow": {
		char:   "↪",
		entity: "&rarrhk;",
		kind:   sym_relation,
	},
	"hslash": {
		char:   "ℏ",
		entity: "&hbar;",
		kind:   sym_alphabetic,
	},
	"hspace": {
		char:   " ",
		entity: "&hairsp;",
		kind:   sym_other,
	},
	"iff": {
		char:   "⟺",
		entity: "&DoubleLongLeftRightArrow;",
		kind:   sym_relation,
	},
	"iiiint": {
		char:   "⨌",
		entity: "&qint;",
		kind:   sym_large,
	},
	"iiint": {
		char:   "∭",
		entity: "&tint;",
		kind:   sym_large,
	},
	"iint": {
		char:   "∬",
		entity: "&Int;",
		kind:   sym_large,
	},
	"image": {
		char:   "⊷",
		entity: "&imof;",
		kind:   sym_relation,
	},
	"imath": {
		char:   "ı",
		entity: "&imath;",
		kind:   sym_alphabetic,
	},
	"in": {
		char:   "∈",
		entity: "&in;",
		kind:   sym_relation,
	},
	"infty": {
		char:   "∞",
		entity: "&infin;",
		kind:   sym_normal,
	},
	"int": {
		char:   "∫",
		entity: "&int;",
		kind:   sym_large,
	},
	"intBar": {
		char:   "⨎",
		entity: "&Barint;",
		kind:   sym_large,
	},
	"intbar": {
		char:   "⨍",
		entity: "&fpartint;",
		kind:   sym_large,
	},
	"intcap": {
		char:   "⨙",
		entity: "&capint;",
		kind:   sym_large,
	},
	"intcup": {
		char:   "⨚",
		entity: "&cupint;",
		kind:   sym_large,
	},
	"intercal": {
		char:   "⊺",
		entity: "&intcal;",
		kind:   sym_binaryop,
	},
	"interleave": {
		char:   "⫴",
		entity: "&vert3;",
		kind:   sym_binaryop,
	},
	"intprod": {
		char:   "⨼",
		entity: "&iprod;",
		kind:   sym_binaryop,
	},
	"intprodr": {
		char:   "⨽",
		entity: "&iprodr;",
		kind:   sym_binaryop,
	},
	"intx": {
		char:   "⨘",
		entity: "&timeint;",
		kind:   sym_large,
	},
	"iota": {
		char:   "ι",
		entity: "&igr;",
		kind:   sym_alphabetic,
	},
	"jupiter": {
		char:   "♃",
		entity: "",
		kind:   sym_other,
	},
	"k": {
		char:   "̨",
		entity: "",
		kind:   sym_diacritic,
	},
	"kappa": {
		char:   "κ",
		entity: "&kappa;",
		kind:   sym_alphabetic,
	},
	"kernelcontraction": {
		char:   "∻",
		entity: "&homtht;",
		kind:   sym_other,
	},
	"lVert": {
		char:   "‖",
		entity: "&Vert;",
		kind:   sym_opening,
	},
	"lambda": {
		char:   "λ",
		entity: "&lambda;",
		kind:   sym_alphabetic,
	},
	"langle": {
		char:   "⟨",
		entity: "&lang;",
		kind:   sym_opening,
	},
	"lazysinv": {
		char:   "∾",
		entity: "&ac;",
		kind:   sym_other,
	},
	"lbrace": {
		char:   "{",
		entity: "&lcub;",
		kind:   sym_opening,
	},
	"lceil": {
		char:   "⌈",
		entity: "&lceil;",
		kind:   sym_opening,
	},
	"le": {
		char:   "≤",
		entity: "&le;",
		kind:   sym_relation,
	},
	"leftarrow": {
		char:   "←",
		entity: "&larr;",
		kind:   sym_relation,
	},
	"leftarrowtail": {
		char:   "↢",
		entity: "&larrtl;",
		kind:   sym_relation,
	},
	"leftarrowtriangle": {
		char:   "⇽",
		entity: "&loarr;",
		kind:   sym_relation,
	},
	"leftharpoondown": {
		char:   "↽",
		entity: "&lhard;",
		kind:   sym_relation,
	},
	"leftharpoonup": {
		char:   "↼",
		entity: "&lharu;",
		kind:   sym_relation,
	},
	"leftleftarrows": {
		char:   "⇇",
		entity: "&larr2;",
		kind:   sym_relation,
	},
	"leftrightarrow": {
		char:   "↔",
		entity: "&harr;",
		kind:   sym_relation,
	},
	"leftrightarrows": {
		char:   "⇆",
		entity: "&lrarr;",
		kind:   sym_relation,
	},
	"leftrightarrowtria*": {
		char:   "⇿",
		entity: "&hoarr;",
		kind:   sym_relation,
	},
	"leftrightharpoons": {
		char:   "⇋",
		entity: "&lrhar;",
		kind:   sym_relation,
	},
	"leftrightsquigarrow": {
		char:   "↭",
		entity: "&harrw;",
		kind:   sym_relation,
	},
	"leftsquigarrow": {
		char:   "↜",
		entity: "",
		kind:   sym_relation,
	},
	"leftthreetimes": {
		char:   "⋋",
		entity: "&lthree;",
		kind:   sym_binaryop,
	},
	"leqq": {
		char:   "≦",
		entity: "&lE;",
		kind:   sym_relation,
	},
	"leqslant": {
		char:   "⩽",
		entity: "&les;",
		kind:   sym_relation,
	},
	"less": {
		char:   "&lt;",
		entity: "&lt;",
		kind:   sym_relation,
	},
	"lessapprox": {
		char:   "⪅",
		entity: "&lap;",
		kind:   sym_relation,
	},
	"lessdot": {
		char:   "⋖",
		entity: "&ldot;",
		kind:   sym_relation,
	},
	"lesseqgtr": {
		char:   "⋚",
		entity: "&leg;",
		kind:   sym_relation,
	},
	"lesseqqgtr": {
		char:   "⪋",
		entity: "&lEg;",
		kind:   sym_relation,
	},
	"lessgtr": {
		char:   "≶",
		entity: "&lg;",
		kind:   sym_relation,
	},
	"lesssim": {
		char:   "≲",
		entity: "&lsim;",
		kind:   sym_relation,
	},
	"lfloor": {
		char:   "⌊",
		entity: "&lfloor;",
		kind:   sym_opening,
	},
	"ll": {
		char:   "≪",
		entity: "&Lt;",
		kind:   sym_relation,
	},
	"llcorner": {
		char:   "⌞",
		entity: "&dlcorn;",
		kind:   sym_opening,
	},
	"lmoustache": {
		char:   "⎰",
		entity: "&lmoust;",
		kind:   sym_other,
	},
	"lnapprox": {
		char:   "⪉",
		entity: "&lnap;",
		kind:   sym_relation,
	},
	"lneq": {
		char:   "⪇",
		entity: "&lne;",
		kind:   sym_relation,
	},
	"lneqq": {
		char:   "≨",
		entity: "&lnE;",
		kind:   sym_relation,
	},
	"lnsim": {
		char:   "⋦",
		entity: "&lnsim;",
		kind:   sym_relation,
	},
	"longleftarrow": {
		char:   "⟵",
		entity: "&xlarr;",
		kind:   sym_relation,
	},
	"longleftrightarrow": {
		char:   "⟷",
		entity: "&xharr;",
		kind:   sym_relation,
	},
	"longmapsto": {
		char:   "⟼",
		entity: "&xmap;",
		kind:   sym_relation,
	},
	"longrightarrow": {
		char:   "⟶",
		entity: "&xrarr;",
		kind:   sym_relation,
	},
	"looparrowleft": {
		char:   "↫",
		entity: "&larrlp;",
		kind:   sym_relation,
	},
	"looparrowright": {
		char:   "↬",
		entity: "&rarrlp;",
		kind:   sym_relation,
	},
	"lowint": {
		char:   "⨜",
		entity: "&lowint;",
		kind:   sym_large,
	},
	"lozenge": {
		char:   "◊",
		entity: "&loz;",
		kind:   sym_other,
	},
	"lrcorner": {
		char:   "⌟",
		entity: "&drcorn;",
		kind:   sym_closing,
	},
	"lt": {
		char:   "&lt;",
		entity: "&lt;",
		kind:   sym_binaryop,
	},
	"ltimes": {
		char:   "⋉",
		entity: "&ltimes;",
		kind:   sym_binaryop,
	},
	"lvert": {
		char:   "|",
		entity: "|",
		kind:   sym_opening,
	},
	"lvertneqq": {
		char:   "≨︀",
		entity: "&lvnE;",
		kind:   sym_relation,
	},
	"mapsto": {
		char:   "↦",
		entity: "&map;",
		kind:   sym_relation,
	},
	"measuredangle": {
		char:   "∡",
		entity: "&angmsd;",
		kind:   sym_normal,
	},
	"mercury": {
		char:   "☿",
		entity: "",
		kind:   sym_other,
	},
	"mho": {
		char:   "℧",
		entity: "&mho;",
		kind:   sym_normal,
	},
	"mid": {
		char:   "∣",
		entity: "&mid;",
		kind:   sym_relation,
	},
	"minusdot": {
		char:   "⨪",
		entity: "&minusdu;",
		kind:   sym_binaryop,
	},
	"mkern4mu": {
		char:   " ",
		entity: "&MediumSpace;",
		kind:   sym_other,
	},
	"mlcp": {
		char:   "⫛",
		entity: "&mlcp;",
		kind:   sym_relation,
	},
	"models": {
		char:   "⊧",
		entity: "&models;",
		kind:   sym_relation,
	},
	"mp": {
		char:   "∓",
		entity: "&mp;",
		kind:   sym_binaryop,
	},
	"mu": {
		char:   "μ",
		entity: "&mgr;",
		kind:   sym_alphabetic,
	},
	"multimap": {
		char:   "⊸",
		entity: "&mumap;",
		kind:   sym_relation,
	},
	"nBumpeq": {
		char:   "≎̸",
		entity: "&nbump;",
		kind:   sym_relation,
	},
	"nLeftarrow": {
		char:   "⇍",
		entity: "&nlArr;",
		kind:   sym_relation,
	},
	"nRightarrow": {
		char:   "⇏",
		entity: "&nrArr;",
		kind:   sym_relation,
	},
	"nVDash": {
		char:   "⊯",
		entity: "&nVDash;",
		kind:   sym_relation,
	},
	"nVdash": {
		char:   "⊮",
		entity: "&nVdash;",
		kind:   sym_relation,
	},
	"nabla": {
		char:   "∇",
		entity: "&Del;",
		kind:   sym_normal,
	},
	"napprox": {
		char:   "≉",
		entity: "&nap;",
		kind:   sym_relation,
	},
	"natural": {
		char:   "♮",
		entity: "&natur;",
		kind:   sym_normal,
	},
	"nbumpeq": {
		char:   "≏̸",
		entity: "&nbumpe;",
		kind:   sym_relation,
	},
	"ncong": {
		char:   "≇",
		entity: "&ncong;",
		kind:   sym_relation,
	},
	"ne": {
		char:   "≠",
		entity: "&ne;",
		kind:   sym_relation,
	},
	"nearrow": {
		char:   "↗",
		entity: "&nearr;",
		kind:   sym_relation,
	},
	"neg": {
		char:   "¬",
		entity: "&not;",
		kind:   sym_normal,
	},
	"neovnwarrow": {
		char:   "⤱",
		entity: "&neonwarr;",
		kind:   sym_other,
	},
	"neovsearrow": {
		char:   "⤮",
		entity: "&neosearr;",
		kind:   sym_other,
	},
	"neptune": {
		char:   "♆",
		entity: "",
		kind:   sym_other,
	},
	"neqsim": {
		char:   "≂̸",
		entity: "&nesim;",
		kind:   sym_relation,
	},
	"nequiv": {
		char:   "≢",
		entity: "&nequiv;",
		kind:   sym_relation,
	},
	"nexists": {
		char:   "∄",
		entity: "&nexist;",
		kind:   sym_normal,
	},
	"ngeq": {
		char:   "≱",
		entity: "&nge;",
		kind:   sym_relation,
	},
	"ngeqslant": {
		char:   "⩾̸",
		entity: "&nges;",
		kind:   sym_relation,
	},
	"ngtr": {
		char:   "≯",
		entity: "&ngt;",
		kind:   sym_relation,
	},
	"ni": {
		char:   "∋",
		entity: "&niv;",
		kind:   sym_relation,
	},
	"nleftarrow": {
		char:   "↚",
		entity: "&nlarr;",
		kind:   sym_relation,
	},
	"nleftrightarrow": {
		char:   "↮",
		entity: "&nharr;",
		kind:   sym_relation,
	},
	"nleq": {
		char:   "≰",
		entity: "&nle;",
		kind:   sym_relation,
	},
	"nleqslant": {
		char:   "⩽̸",
		entity: "&nles;",
		kind:   sym_relation,
	},
	"nless": {
		char:   "≮",
		entity: "&nlt;",
		kind:   sym_relation,
	},
	"nmid": {
		char:   "∤",
		entity: "&nmid;",
		kind:   sym_relation,
	},
	"nolinebreak": {
		char:   "\u2060",
		entity: "&NoBreak;",
		kind:   sym_normal,
	},
	"notgreaterless": {
		char:   "≹",
		entity: "&ntgl;",
		kind:   sym_relation,
	},
	"notin": {
		char:   "∉",
		entity: "&notin;",
		kind:   sym_relation,
	},
	"notlessgreater": {
		char:   "≸",
		entity: "&ntlg;",
		kind:   sym_relation,
	},
	"nparallel": {
		char:   "∦",
		entity: "&npar;",
		kind:   sym_relation,
	},
	"nprec": {
		char:   "⊀",
		entity: "&npr;",
		kind:   sym_relation,
	},
	"npreceq": {
		char:   "⪯̸",
		entity: "&npre;",
		kind:   sym_relation,
	},
	"nprecsim": {
		char:   "≾̸",
		entity: "&nprsim;",
		kind:   sym_relation,
	},
	"nrightarrow": {
		char:   "↛",
		entity: "&nrarr;",
		kind:   sym_relation,
	},
	"nsim": {
		char:   "≁",
		entity: "&nsim;",
		kind:   sym_relation,
	},
	"nsime": {
		char:   "≄",
		entity: "&nsime;",
		kind:   sym_relation,
	},
	"nsubset": {
		char:   "⊄",
		entity: "&nsub;",
		kind:   sym_relation,
	},
	"nsubseteq": {
		char:   "⊈",
		entity: "&nsube;",
		kind:   sym_relation,
	},
	"nsubseteqq": {
		char:   "⫅̸",
		entity: "&nsubE;",
		kind:   sym_relation,
	},
	"nsucc": {
		char:   "⊁",
		entity: "&nsc;",
		kind:   sym_relation,
	},
	"nsucceq": {
		char:   "⪰̸",
		entity: "&nsce;",
		kind:   sym_relation,
	},
	"nsuccsim": {
		char:   "≿̸",
		entity: "&nscsim;",
		kind:   sym_relation,
	},
	"nsupset": {
		char:   "⊅",
		entity: "&nsup;",
		kind:   sym_relation,
	},
	"nsupseteq": {
		char:   "⊉",
		entity: "&nsupe;",
		kind:   sym_relation,
	},
	"nsupseteqq": {
		char:   "⫆̸",
		entity: "&nsupE;",
		kind:   sym_relation,
	},
	"ntriangleleft": {
		char:   "⋪",
		entity: "&nltri;",
		kind:   sym_relation,
	},
	"ntrianglelefteq": {
		char:   "⋬",
		entity: "&nltrie;",
		kind:   sym_relation,
	},
	"ntriangleright": {
		char:   "⋫",
		entity: "&nrtri;",
		kind:   sym_relation,
	},
	"ntrianglerighteq": {
		char:   "⋭",
		entity: "&nrtrie;",
		kind:   sym_relation,
	},
	"nu": {
		char:   "ν",
		entity: "&ngr;",
		kind:   sym_alphabetic,
	},
	"nvDash": {
		char:   "⊭",
		entity: "&nvDash;",
		kind:   sym_relation,
	},
	"nvdash": {
		char:   "⊬",
		entity: "&nvdash;",
		kind:   sym_relation,
	},
	"nwarrow": {
		char:   "↖",
		entity: "&nwarr;",
		kind:   sym_relation,
	},
	"nwovnearrow": {
		char:   "⤲",
		entity: "&nwonearr;",
		kind:   sym_other,
	},
	"obar": {
		char:   "⌽",
		entity: "&ovbar;",
		kind:   sym_binaryop,
	},
	"obslash": {
		char:   "⦸",
		entity: "&obsol;",
		kind:   sym_binaryop,
	},
	"odot": {
		char:   "⊙",
		entity: "&odot;",
		kind:   sym_binaryop,
	},
	"oiiint": {
		char:   "∰",
		entity: "&Cconint;",
		kind:   sym_large,
	},
	"oiint": {
		char:   "∯",
		entity: "&Conint;",
		kind:   sym_large,
	},
	"oint": {
		char:   "∮",
		entity: "&oint;",
		kind:   sym_large,
	},
	"omega": {
		char:   "ω",
		entity: "&ohgr;",
		kind:   sym_alphabetic,
	},
	"ominus": {
		char:   "⊖",
		entity: "&ominus;",
		kind:   sym_binaryop,
	},
	"openbracketleft": {
		char:   "〚",
		entity: "&lobrk;",
		kind:   sym_opening,
	},
	"openbracketright": {
		char:   "〛",
		entity: "&robrk;",
		kind:   sym_closing,
	},
	"oplus": {
		char:   "⊕",
		entity: "&oplus;",
		kind:   sym_binaryop,
	},
	"original": {
		char:   "⊶",
		entity: "&origof;",
		kind:   sym_relation,
	},
	"oslash": {
		char:   "⊘",
		entity: "&osol;",
		kind:   sym_binaryop,
	},
	"otimes": {
		char:   "⊗",
		entity: "&otimes;",
		kind:   sym_binaryop,
	},
	"parallel": {
		char:   "∥",
		entity: "&par;",
		kind:   sym_relation,
	},
	"partial": {
		char:   "∂",
		entity: "&part;",
		kind:   sym_normal,
	},
	"partialmeetcontraction": {
		char:   "⪣",
		entity: "&Ltbar;",
		kind:   sym_relation,
	},
	"perp": {
		char:   "⊥",
		entity: "&bot;",
		kind:   sym_relation,
	},
	"perspcorrespond": {
		char:   "⩞",
		entity: "&Barwedl;",
		kind:   sym_binaryop,
	},
	"phi": {
		char:   "ϕ",
		entity: "&phi;",
		kind:   sym_alphabetic,
	},
	"pi": {
		char:   "π",
		entity: "&pgr;",
		kind:   sym_alphabetic,
	},
	"pitchfork": {
		char:   "⋔",
		entity: "&fork;",
		kind:   sym_other,
	},
	"plusdot": {
		char:   "⨥",
		entity: "&plusdu;",
		kind:   sym_binaryop,
	},
	"pm": {
		char:   "±",
		entity: "&pm;",
		kind:   sym_binaryop,
	},
	"prec": {
		char:   "≺",
		entity: "&pr;",
		kind:   sym_relation,
	},
	"precapprox": {
		char:   "⪷",
		entity: "&prap;",
		kind:   sym_relation,
	},
	"preccurlyeq": {
		char:   "≼",
		entity: "&cupre;",
		kind:   sym_relation,
	},
	"preceq": {
		char:   "⪯",
		entity: "&pre;",
		kind:   sym_relation,
	},
	"precnapprox": {
		char:   "⪹",
		entity: "&prnap;",
		kind:   sym_relation,
	},
	"precneqq": {
		char:   "⪵",
		entity: "&prnE;",
		kind:   sym_relation,
	},
	"precnsim": {
		char:   "⋨",
		entity: "&prnsim;",
		kind:   sym_relation,
	},
	"precsim": {
		char:   "≾",
		entity: "&prsim;",
		kind:   sym_relation,
	},
	"prime": {
		char:   "′",
		entity: "&prime;",
		kind:   sym_other,
	},
	"prod": {
		char:   "∏",
		entity: "&prod;",
		kind:   sym_large,
	},
	"propto": {
		char:   "∝",
		entity: "&prop;",
		kind:   sym_relation,
	},
	"psi": {
		char:   "ψ",
		entity: "&psgr;",
		kind:   sym_alphabetic,
	},
	"questeq": {
		char:   "≟",
		entity: "&equest;",
		kind:   sym_relation,
	},
	"rVert": {
		char:   "‖",
		entity: "&Vert;",
		kind:   sym_opening,
	},
	"rangle": {
		char:   "⟩",
		entity: "&rang;",
		kind:   sym_closing,
	},
	"rbrace": {
		char:   "}",
		entity: "&rcub;",
		kind:   sym_closing,
	},
	"rceil": {
		char:   "⌉",
		entity: "&rceil;",
		kind:   sym_closing,
	},
	"rdiagovfdiag": {
		char:   "⤫",
		entity: "&rdiofdi;",
		kind:   sym_other,
	},
	"rdiagovsearrow": {
		char:   "⤰",
		entity: "&rdosearr;",
		kind:   sym_other,
	},
	"recorder": {
		char:   "⌕",
		entity: "&telrec;",
		kind:   sym_other,
	},
	"rfloor": {
		char:   "⌋",
		entity: "&rfloor;",
		kind:   sym_closing,
	},
	"rho": {
		char:   "ρ",
		entity: "&rgr;",
		kind:   sym_alphabetic,
	},
	"rightangle": {
		char:   "∟",
		entity: "&ang90;",
		kind:   sym_normal,
	},
	"rightanglearc": {
		char:   "⊾",
		entity: "&angrtvb;",
		kind:   sym_other,
	},
	"rightarrow": {
		char:   "→",
		entity: "&rarr;",
		kind:   sym_relation,
	},
	"rightarrowtail": {
		char:   "↣",
		entity: "&rarrtl;",
		kind:   sym_relation,
	},
	"rightarrowtriangle": {
		char:   "⇾",
		entity: "&roarr;",
		kind:   sym_relation,
	},
	"rightharpoondown": {
		char:   "⇁",
		entity: "&rhard;",
		kind:   sym_relation,
	},
	"rightharpoonup": {
		char:   "⇀",
		entity: "&rharu;",
		kind:   sym_relation,
	},
	"rightleftarrows": {
		char:   "⇄",
		entity: "&rlarr;",
		kind:   sym_relation,
	},
	"rightleftharpoons": {
		char:   "⇌",
		entity: "&rlhar;",
		kind:   sym_relation,
	},
	"rightmoon": {
		char:   "☾",
		entity: "",
		kind:   sym_other,
	},
	"rightrightarrows": {
		char:   "⇉",
		entity: "&rarr2;",
		kind:   sym_relation,
	},
	"rightsquigarrow": {
		char:   "↝",
		entity: "&rarrw;",
		kind:   sym_relation,
	},
	"rightthreetimes": {
		char:   "⋌",
		entity: "&rthree;",
		kind:   sym_binaryop,
	},
	"risingdotseq": {
		char:   "≓",
		entity: "&erDot;",
		kind:   sym_relation,
	},
	"rmoustache": {
		char:   "⎱",
		entity: "&rmoust;",
		kind:   sym_other,
	},
	"rtimes": {
		char:   "⋊",
		entity: "&rtimes;",
		kind:   sym_binaryop,
	},
	"rvert": {
		char:   "|",
		entity: "|",
		kind:   sym_closing,
	},
	"saturn": {
		char:   "♄",
		entity: "",
		kind:   sym_other,
	},
	"searrow": {
		char:   "↘",
		entity: "&drarr;",
		kind:   sym_relation,
	},
	"sector": {
		char:   "⌔",
		entity: "&#x2314",
		kind:   sym_other,
	},
	"seovnearrow": {
		char:   "⤭",
		entity: "&seonearr;",
		kind:   sym_other,
	},
	"setminus": {
		char:   "∖",
		entity: "&setmn;",
		kind:   sym_binaryop,
	},
	"sharp": {
		char:   "♯",
		entity: "&sharp;",
		kind:   sym_normal,
	},
	"shuffle": {
		char:   "⧢",
		entity: "&shuffle;",
		kind:   sym_other,
	},
	"sigma": {
		char:   "σ",
		entity: "&sgr;",
		kind:   sym_alphabetic,
	},
	"sim": {
		char:   "∼",
		entity: "&sim;",
		kind:   sym_relation,
	},
	// SPURIOUS
	//"sim\\joinrel\\leadsto": {
	//	char:   "⟿",
	//	entity: "&dzigrarr;",
	//	kind:   sym_relation,
	//},
	"simeq": {
		char:   "≃",
		entity: "&sime;",
		kind:   sym_relation,
	},
	"smile": {
		char:   "⌣",
		entity: "&smile;",
		kind:   sym_relation,
	},
	"sphericalangle": {
		char:   "∢",
		entity: "&angsph;",
		kind:   sym_normal,
	},
	"sqcap": {
		char:   "⊓",
		entity: "&sqcap;",
		kind:   sym_binaryop,
	},
	"sqcup": {
		char:   "⊔",
		entity: "&sqcup;",
		kind:   sym_binaryop,
	},
	"sqrint": {
		char:   "⨖",
		entity: "&quatint;",
		kind:   sym_large,
	},
	"sqsubset": {
		char:   "⊏",
		entity: "&sqsub;",
		kind:   sym_relation,
	},
	"sqsubseteq": {
		char:   "⊑",
		entity: "&sqsube;",
		kind:   sym_relation,
	},
	"sqsupset": {
		char:   "⊐",
		entity: "&sqsup;",
		kind:   sym_relation,
	},
	"sqsupseteq": {
		char:   "⊒",
		entity: "&sqsupe;",
		kind:   sym_relation,
	},
	"square": {
		char:   "□",
		entity: "&squ;",
		kind:   sym_other,
	},
	"star": {
		char:   "⋆",
		entity: "&Star;",
		kind:   sym_binaryop,
	},
	"starequal": {
		char:   "≛",
		entity: "",
		kind:   sym_relation,
	},
	"sterling": {
		char:   "£",
		entity: "&pound;",
		kind:   sym_normal,
	},
	"subset": {
		char:   "⊂",
		entity: "&sub;",
		kind:   sym_relation,
	},
	"subseteq": {
		char:   "⊆",
		entity: "&sube;",
		kind:   sym_relation,
	},
	"subseteqq": {
		char:   "⫅",
		entity: "&subE;",
		kind:   sym_relation,
	},
	"subsetneq": {
		char:   "⊊",
		entity: "&subne;",
		kind:   sym_relation,
	},
	"subsetneqq": {
		char:   "⫋",
		entity: "&subnE;",
		kind:   sym_relation,
	},
	"succ": {
		char:   "≻",
		entity: "&sc;",
		kind:   sym_relation,
	},
	"succapprox": {
		char:   "⪸",
		entity: "&scap;",
		kind:   sym_relation,
	},
	"succcurlyeq": {
		char:   "≽",
		entity: "&sccue;",
		kind:   sym_relation,
	},
	"succeq": {
		char:   "⪰",
		entity: "&sce;",
		kind:   sym_relation,
	},
	"succnapprox": {
		char:   "⪺",
		entity: "&scnap;",
		kind:   sym_relation,
	},
	"succneqq": {
		char:   "⪶",
		entity: "&scnE;",
		kind:   sym_relation,
	},
	"succnsim": {
		char:   "⋩",
		entity: "&scnsim;",
		kind:   sym_relation,
	},
	"succsim": {
		char:   "≿",
		entity: "&scsim;",
		kind:   sym_relation,
	},
	"sum": {
		char:   "∑",
		entity: "&sum;",
		kind:   sym_large,
	},
	"supset": {
		char:   "⊃",
		entity: "&sup;",
		kind:   sym_relation,
	},
	"supseteq": {
		char:   "⊇",
		entity: "&supe;",
		kind:   sym_relation,
	},
	"supseteqq": {
		char:   "⫆",
		entity: "&supE;",
		kind:   sym_relation,
	},
	"supsetneq": {
		char:   "⊋",
		entity: "&supne;",
		kind:   sym_relation,
	},
	"supsetneqq": {
		char:   "⫌",
		entity: "&supnE;",
		kind:   sym_relation,
	},
	"surd": {
		char:   "√",
		entity: "&Sqrt;",
		kind:   sym_other,
	},
	"swarrow": {
		char:   "↙",
		entity: "&dlarr;",
		kind:   sym_relation,
	},
	"tau": {
		char:   "τ",
		entity: "&tau;",
		kind:   sym_alphabetic,
	},
	"textTheta": {
		char:   "ϴ",
		entity: "&Thetav;",
		kind:   sym_alphabetic,
	},
	"therefore": {
		char:   "∴",
		entity: "&there4;",
		kind:   sym_normal,
	},
	"theta": {
		char:   "θ",
		entity: "&theta;",
		kind:   sym_alphabetic,
	},
	"tilde": {
		char:   "̃",
		entity: "",
		kind:   sym_diacritic,
	},
	"tildetrpl": {
		char:   "≋",
		entity: "&apid;",
		kind:   sym_relation,
	},
	"times": {
		char:   "×",
		entity: "&times;",
		kind:   sym_binaryop,
	},
	"to": {
		char:   "→",
		entity: "&rarr;",
		kind:   sym_relation,
	},
	"toea": {
		char:   "⤨",
		entity: "&toea;",
		kind:   sym_relation,
	},
	"tona": {
		char:   "⤧",
		entity: "&nwnear;",
		kind:   sym_relation,
	},
	"top": {
		char:   "⊤",
		entity: "&top;",
		kind:   sym_normal,
	},
	"tosa": {
		char:   "⤩",
		entity: "&tosa;",
		kind:   sym_relation,
	},
	"towa": {
		char:   "⤪",
		entity: "&swnwar;",
		kind:   sym_relation,
	},
	"triangle": {
		char:   "△",
		entity: "&#x25B3;",
		kind:   sym_other,
	},
	"triangledown": {
		char:   "▿",
		entity: "&dtri;",
		kind:   sym_other,
	},
	"triangleleft": {
		char:   "◃",
		entity: "&ltri;",
		kind:   sym_other,
	},
	"trianglelefteq": {
		char:   "⊴",
		entity: "&ltrie;",
		kind:   sym_relation,
	},
	"triangleq": {
		char:   "≜",
		entity: "&trie;",
		kind:   sym_relation,
	},
	"triangleright": {
		char:   "▹",
		entity: "&rtri;",
		kind:   sym_other,
	},
	"trianglerighteq": {
		char:   "⊵",
		entity: "&rtrie;",
		kind:   sym_relation,
	},
	"twoheadleftarrow": {
		char:   "↞",
		entity: "&Larr;",
		kind:   sym_relation,
	},
	"twoheadrightarrow": {
		char:   "↠",
		entity: "&Rarr;",
		kind:   sym_relation,
	},
	"twoheadrightarrowtail": {
		char:   "⤖",
		entity: "&Rarrtl;",
		kind:   sym_relation,
	},
	"u": {
		char:   "˘",
		entity: "&breve;",
		kind:   sym_other,
	},
	"ulcorner": {
		char:   "⌜",
		entity: "&ulcorn;",
		kind:   sym_opening,
	},
	"uparrow": {
		char:   "↑",
		entity: "&uarr;",
		kind:   sym_relation,
	},
	"updownarrow": {
		char:   "↕",
		entity: "&varr;",
		kind:   sym_relation,
	},
	"upharpoonleft": {
		char:   "↾",
		entity: "&uharr;",
		kind:   sym_relation,
	},
	"upharpoonright": {
		char:   "↿",
		entity: "&uharl;",
		kind:   sym_relation,
	},
	"upint": {
		char:   "⨛",
		entity: "&upint;",
		kind:   sym_large,
	},
	"uplus": {
		char:   "⊎",
		entity: "&uplus;",
		kind:   sym_binaryop,
	},
	"upsilon": {
		char:   "υ",
		entity: "&ugr;",
		kind:   sym_alphabetic,
	},
	"upuparrows": {
		char:   "⇈",
		entity: "&uarr2;",
		kind:   sym_relation,
	},
	"uranus": {
		char:   "♅",
		entity: "",
		kind:   sym_other,
	},
	"urcorner": {
		char:   "⌝",
		entity: "&urcorn;",
		kind:   sym_closing,
	},
	"vDash": {
		char:   "⊨",
		entity: "&vDash;",
		kind:   sym_relation,
	},
	"varepsilon": {
		char:   "ɛ",
		entity: "",
		kind:   sym_other,
	},
	"varkappa": {
		char:   "ϰ",
		entity: "&kappav;",
		kind:   sym_alphabetic,
	},
	"varnothing": {
		char:   "∅",
		entity: "&empty;",
		kind:   sym_normal,
	},
	"varphi": {
		char:   "φ",
		entity: "&phgr;",
		kind:   sym_alphabetic,
	},
	"varpi": {
		char:   "ϖ",
		entity: "&piv;",
		kind:   sym_alphabetic,
	},
	"varrho": {
		char:   "ϱ",
		entity: "&rhov;",
		kind:   sym_alphabetic,
	},
	"varsigma": {
		char:   "ς",
		entity: "&sfgr;",
		kind:   sym_alphabetic,
	},
	"varsubsetneqq": {
		char:   "⊊︀",
		entity: "&vsubne;",
		kind:   sym_relation,
	},
	"varsupsetneq": {
		char:   "⊋︀",
		entity: "&vsupne;",
		kind:   sym_relation,
	},
	"vartheta": {
		char:   "ϑ",
		entity: "&thetav;",
		kind:   sym_alphabetic,
	},
	"vartriangle": {
		char:   "▵",
		entity: "&utri;",
		kind:   sym_other,
	},
	"vartriangleleft": {
		char:   "⊲",
		entity: "&vltri;",
		kind:   sym_relation,
	},
	"vartriangleright": {
		char:   "⊳",
		entity: "&vrtri;",
		kind:   sym_relation,
	},
	"vdash": {
		char:   "⊢",
		entity: "&vdash;",
		kind:   sym_relation,
	},
	"vdots": {
		char:   "⋮",
		entity: "&vellip;",
		kind:   sym_other,
	},
	"vee": {
		char:   "∨",
		entity: "&or;",
		kind:   sym_binaryop,
	},
	"veebar": {
		char:   "⊻",
		entity: "&veebar;",
		kind:   sym_binaryop,
	},
	"vert": {
		char:   "|",
		entity: "&vert;",
		kind:   sym_other,
	},
	"verymuchless": {
		char:   "⋘",
		entity: "&Ll;",
		kind:   sym_relation,
	},
	"wedge": {
		char:   "∧",
		entity: "&and;",
		kind:   sym_binaryop,
	},
	"wedgeq": {
		char:   "≙",
		entity: "&wedgeq;",
		kind:   sym_relation,
	},
	"wp": {
		char:   "℘",
		entity: "&wp;",
		kind:   sym_alphabetic,
	},
	"wr": {
		char:   "≀",
		entity: "&wr;",
		kind:   sym_binaryop,
	},
	"xi": {
		char:   "ξ",
		entity: "&xgr;",
		kind:   sym_alphabetic,
	},
	"yen": {
		char:   "¥",
		entity: "&yen;",
		kind:   sym_normal,
	},
	"zeta": {
		char:   "ζ",
		entity: "&zeta;",
		kind:   sym_alphabetic,
	},
}
