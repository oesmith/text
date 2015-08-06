package text

import (
	"strings"

	"golang.org/x/exp/utf8string"
)

// Ruby implementation of the Double Metaphone algorithm by Lawrence Philips,
// originally published in the June 2000 issue of C/C++ Users Journal.
//
// Based on double_metaphone.rb from https://github.com/threedaymonk/text
//
// Returns the primary and secondary double metaphone tokens (the secondary
// will be nil if equal to the primary).
func DoubleMetaphone(str string) (string, string) {
	primary := make([]string, 0, 5)
	secondary := make([]string, 0, 5)
	original := utf8string.NewString(strings.ToUpper(str))
	current := 0

	if mm(ss(original, 0, 2), []string{"GN", "KN", "PN", "WR", "PS"}) {
		current += 1
	}
	if ss(original, 0, 1) == "X" {
		primary = append(primary, "S")
		secondary = append(secondary, "S")
		current += 1
	}
	for len(primary) < 4 || len(secondary) < 4 {
		if current > original.RuneCount() {
			break
		}
		p, s, n := doubleMetaphoneLookup(original, current)
		if p != "" {
			primary = append(primary, p)
		}
		if s != "" {
			secondary = append(secondary, s)
		}
		current += n
	}
	p := strings.Join(primary, "")
	if len(p) > 4 {
		p = p[0:4]
	}
	s := strings.Join(secondary, "")
	if len(s) > 4 {
		s = s[0:4]
	}
	if p == s {
		return p, ""
	}
	return p, s
}

func slavoGermanic(s *utf8string.String) bool {
	str := s.String()
	return strings.Contains(str, "W") ||
		strings.Contains(str, "K") ||
		strings.Contains(str, "CZ") ||
		strings.Contains(str, "WITZ")
}

func vowel(str string) bool {
	return str == "A" || str == "E" || str == "I" || str == "O" || str == "U" || str == "Y"
}

func mm(str string, matches []string) bool {
	for _, s := range(matches) {
		if s == str {
			return true
		}
	}
	return false
}

func ss(str *utf8string.String, pos, length int) string {
	p0, p1 := pos, pos + length
	if p0 < 0 {
		p0 = 0
	}
	l := p1 - p0
	if p1 > str.RuneCount() {
		p1 = str.RuneCount()
	}
	s := str.Slice(p0, p1)
	if l > (p1 - p0) {
		s = s + strings.Repeat(" ", l - (p1 - p0))
	}
	return s
}

func doubleMetaphoneLookup(str *utf8string.String, pos int) (string, string, int) {
	switch ss(str, pos, 1) {
	case "A", "E", "I", "O", "U", "Y":
		if 0 == pos {
			return "A", "A", 1
		} else {
			return "", "", 1
		}
	case "B":
		n := 1
		if str.RuneCount() > pos + 1 && "B" == ss(str, pos + 1, 1) {
			n = 2
		}
		return "P", "P", n
	case "Ç":
		return "S", "S", 1
	case "C":
		if pos > 1 &&
				!vowel(ss(str, pos - 2, 1)) &&
				"ACH" == ss(str, pos - 1, 3) &&
				ss(str, pos + 2, 1) != "I" && (
					ss(str, pos + 2, 1) != "E" ||
					mm(ss(str, pos - 2, 6), []string{"BACHER", "MACHER"})) {
			return "K", "K", 2
		} else if 0 == pos && "CAESAR" == ss(str, pos, 6) {
			return "S", "S", 2
		} else if "CHIA" == ss(str, pos, 4) {
			return "K", "K", 2
		} else if "CH" == ss(str, pos, 2) {
			if pos > 0 && "CHAE" == ss(str, pos, 4) {
				return "K", "X", 2
			} else if 0 == pos &&
					(mm(ss(str, pos + 1, 5), []string{"HARAC", "HARIS"}) ||
					mm(ss(str, pos + 1, 3), []string{"HOR", "HYM", "HIA", "HEM"})) &&
					ss(str, 0, 5) != "CHORE" {
				return "K", "K", 2
			} else if mm(ss(str, 0, 4), []string{"VAN ", "VON "}) ||
					"SCH" == ss(str, 0, 3) ||
					mm(ss(str, pos - 2, 6), []string{"ORCHES", "ARCHIT", "ORCHID"}) ||
					mm(ss(str, pos + 2, 1), []string{"T", "S"}) || (
						(pos == 0 || mm(ss(str, pos - 1, 1), []string{"A", "O", "U", "E"})) &&
						mm(ss(str, pos + 2, 1), []string{"L", "R", "N", "M", "B", "H", "F", "V", "W", " "})) {
				return "K", "K", 2
			} else if pos > 0 {
				if "MC" == ss(str, 0, 2) {
					return "K", "K", 2
				} else {
					return "X", "K", 2
				}
			} else {
				return "X", "X", 2
			}
		} else if "CZ" == ss(str, pos, 2) && "WICZ" != ss(str, pos - 2, 4) {
			return "S", "X", 2
		} else if "CIA" == ss(str, pos + 1, 3) {
			return "X", "X", 3
		} else if "CC" == ss(str, pos, 2) && !(1 == pos && "M" == ss(str, 0, 1)) {
			if mm(ss(str, pos + 2, 1), []string{"I", "E", "H"}) && "HU" != ss(str, pos + 2, 2) {
				if (1 == pos && "A" == ss(str, pos - 1, 1)) ||
						mm(ss(str, pos - 1, 5), []string{"UCCEE", "UCCES"}) {
					return "KS", "KS", 3
				} else {
					return "X", "X", 3
				}
			} else {
				return "K", "K", 2
			}
		} else if mm(ss(str, pos, 2), []string{"CK", "CG", "CQ"}) {
			return "K", "K", 2
		} else if mm(ss(str, pos, 2), []string{"CI", "CE", "CY"}) {
			if mm(ss(str, pos, 3), []string{"CIO", "CIE", "CIA"}) {
				return "S", "X", 2
			} else {
				return "S", "S", 2
			}
		} else if mm(ss(str, pos + 1, 2), []string{" C", " Q", "G"}) {
			return "K", "K", 3
		} else if mm(ss(str, pos + 1, 1), []string{"C", "K", "Q"}) && !mm(ss(str, pos + 1, 2), []string{"CE", "CI"}) {
			return "K", "K", 2
		} else {
			return "K", "K", 1
		}
	case "D":
		if "DG" == ss(str, pos, 2) {
			if mm(ss(str, pos + 2, 1), []string{"I", "E", "Y"}) {
				return "J", "J", 3
			} else {
				return "TK", "TK", 2
			}
		} else if mm(ss(str, pos, 2), []string{"DT", "DD"}) {
			return "T", "T", 2
		} else {
			return "T", "T", 1
		}
	case "F":
		if "F" == ss(str, pos + 1, 1) {
			return "F", "F", 2
		} else {
			return "F", "F", 1
		}
	case "G":
		if "H" == ss(str, pos + 1, 1) {
			if pos > 0 && !vowel(ss(str, pos - 1, 1)) {
				return "K", "K", 2
			} else if 0 == pos {
				if "I" == ss(str, pos + 2, 1) {
					return "J", "J", 2
				} else {
					return "K", "K", 2
				}
			} else if (pos > 1 && mm(ss(str, pos - 2, 1), []string{"B", "H", "D"})) ||
					(pos > 2 && mm(ss(str, pos - 3, 1), []string{"B", "H", "D"})) ||
					(pos > 3 && mm(ss(str, pos - 4, 1), []string{"B", "H"})) {
				return "", "", 2
			} else if pos > 2 && "U" == ss(str, pos - 1, 1) && mm(ss(str, pos - 3, 1), []string{"C", "G", "L", "R", "T"}) {
				return "F", "F", 2
			} else if pos > 0 && "I" != ss(str, pos - 1, 1) {
				return "K", "K", 2
			} else {
				return "", "", 2
			}
		} else if "N" == ss(str, pos + 1, 1) {
			if 1 == pos && vowel(ss(str, 0, 1)) && !slavoGermanic(str) {
				return "KN", "N", 2
			} else if "EY" != ss(str, pos + 2, 2) && "Y" != ss(str, pos + 1, 1) && !slavoGermanic(str) {
				return "N", "KN", 2
			} else {
				return "KN", "KN", 2
			}
		} else if "LI" == ss(str, pos + 1, 2) && !slavoGermanic(str) {
			return "KL", "L", 2
		} else if 0 == pos && ("Y" == ss(str, pos + 1, 1) ||
				mm(ss(str, pos + 1, 2), []string{"ES", "EP", "EB", "EL", "EY", "EI", "ER", "IB", "IL", "IN", "IE"})) {
			return "K", "J", 2
		} else if ("ER" == ss(str, pos + 1, 2) || "Y" == ss(str, pos + 1, 1)) &&
				!mm(ss(str, 0, 6), []string{"DANGER", "RANGER", "MANGER"}) &&
				!mm(ss(str, pos - 1, 1), []string{"E", "I"}) &&
				!mm(ss(str, pos - 1, 3), []string{"RGY", "OGY"}) {
			return "K", "J", 2
		} else if mm(ss(str, pos + 1, 1), []string{"E", "I", "Y"}) ||
				mm(ss(str, pos - 1, 4), []string{"AGGI", "OGGI"}) {
			if mm(ss(str, 0, 4), []string{"VAN ", "VON "}) || "SCH" == ss(str, 0, 3) || "ET" == ss(str, pos + 1, 2) {
				return "K", "K", 2
			} else if "IER " == ss(str, pos + 1, 4) {
				return "J", "J", 2
			} else {
				return "J", "K", 2
			}
		} else if "G" == ss(str, pos + 1, 1) {
			return "K", "K", 2
		} else {
			return "K", "K", 1
		}
	case "H":
		if (0 == pos || vowel(ss(str, pos - 1, 1))) && vowel(ss(str, pos + 1, 1)) {
			return "H", "H", 2
		} else {
			return "", "", 1
		}
	case "J":
		if "JOSE" == ss(str, pos, 4) || "SAN " == ss(str, 0, 4) {
			if (0 == pos && " " == ss(str, pos + 4, 1)) || "SAN " == ss(str, 0, 4) {
				return "H", "H", 1
			} else {
				return "J", "H", 1
			}
		} else {
			c := 1
			if "J" == ss(str, pos + 1, 1) {
				c = 2
			}
			if 0 == pos && "JOSE" != ss(str, pos, 4) {
				return "J", "A", c
			} else if vowel(ss(str, pos - 1, 1)) && !slavoGermanic(str) && mm(ss(str, pos + 1, 1), []string{"A", "O"}) {
				return "J", "H", c
			} else if pos == str.RuneCount() - 1 {
				return "J", "", c
			} else if !mm(ss(str, pos + 1, 1), []string{"L", "T", "K", "S", "N", "M", "B", "Z"}) &&
					!mm(ss(str, pos - 1, 1), []string{"S", "K", "L"}) {
				return "J", "J", c
			} else {
				return "", "", c
			}
		}
	case "K":
		c := 1
		if "K" == ss(str, pos + 1, 1) {
			c = 2
		}
		return "K", "K", c
	case "L":
		if "L" == ss(str, pos + 1, 1) {
			if (str.RuneCount() - 3 == pos && mm(ss(str, pos - 1, 4), []string{"ILLO", "ILLA", "ALLE"})) ||
					(mm(ss(str, str.RuneCount() - 2, 2), []string{"AS", "OS"}) || mm(ss(str, str.RuneCount() - 1, 1), []string{"A", "O"})) &&
					"ALLE" == ss(str, pos - 1, 4) {
				return "L", "", 2
			} else {
				return "L", "L", 2
			}
		} else {
			return "L", "L", 1
		}
	case "M":
		if ("UMB" == ss(str, pos - 1, 3) && (str.RuneCount() - 2 == pos || "ER" == ss(str, pos + 2, 2))) ||
				"M" == ss(str, pos + 1, 1) {
			return "M", "M", 2
		} else {
			return "M", "M", 1
		}
	case "N":
		c := 1
		if "N" == ss(str, pos + 1, 1) {
			c = 2
		}
		return "N", "N", c
	case "Ñ":
		return "N", "N", 1
	case "P":
		if "H" == ss(str, pos + 1, 1) {
			return "F", "F", 2
		} else {
			c := 1
			if mm(ss(str, pos + 1, 1), []string{"P", "B"}) {
				c = 2
			}
			return "P", "P", c
		}
	case "Q":
		c := 1
		if "Q" == ss(str, pos + 1, 1) {
			c = 2
		}
		return "K", "K", c
	case "R":
		c := 1
		if "R" == ss(str, pos + 1, 1) {
			c = 2
		}
		if pos == str.RuneCount() - 1 && !slavoGermanic(str) && "IE" == ss(str, pos - 2, 2) &&
				!mm(ss(str, pos - 4, 2), []string{"ME", "MA"}) {
			return "", "R", c
		} else {
			return "R", "R", c
		}
	case "S":
		if mm(ss(str, pos - 1, 3), []string{"ISL", "YSL"}) {
			return "", "", 1
		} else if 0 == pos && "SUGAR" == ss(str, 0, 5) {
			return "X", "S", 1
		} else if "H" == ss(str, pos + 1, 1) {
			if mm(ss(str, pos + 1, 4), []string{"HEIM", "HOEK", "HOLM", "HOLZ"}) {
				return "S", "S", 2
			} else {
				return "X", "X", 2
			}
		} else if mm(ss(str, pos, 3), []string{"SIO", "SIA"}) || "SIAN" == ss(str, pos, 4) {
			if slavoGermanic(str) {
				return "S", "S", 3
			} else {
				return "S", "X", 3
			}
		} else if 0 == pos && mm(ss(str, pos + 1, 1), []string{"M", "N", "L", "W"}) {
			return "S", "X", 1
		} else if "Z" == ss(str, pos + 1, 1) {
			return "S", "X", 2
		} else if "C" == ss(str, pos + 1, 1) {
			if "H" == ss(str, pos + 2, 1) {
				if mm(ss(str, pos + 3, 2), []string{"OO", "UY", "ED", "EM"}) {
					return "SK", "SK", 3
				} else if mm(ss(str, pos + 3, 2), []string{"ER", "EN"}) {
					return "X", "SK", 3
				}
			} else if mm(ss(str, pos + 2, 1), []string{"I", "E", "Y"}) {
				return "S", "S", 3
			} else {
				return "SK", "SK", 3
			}
		} else {
			p := "S"
			if pos == str.RuneCount() - 1 && mm(ss(str, pos - 2, 2), []string{"AI", "OI"}) {
				p = ""
			}
			c := 1
			if mm(ss(str, pos + 1, 1), []string{"S", "Z"}) {
				c = 2
			}
			return p, "S", c
		}
	case "T":
		if "TION" == ss(str, pos, 4) {
			return "X", "X", 3
		} else if mm(ss(str, pos, 3), []string{"TIA", "TCH"}) {
			return "X", "X", 3
		} else if "TH" == ss(str, pos, 2) || "TTH" == ss(str, pos, 3) {
			if mm(ss(str, pos + 2, 2), []string{"OM", "AM"}) ||
					mm(ss(str, 0, 4), []string{"VAN ", "VON "}) ||
					"SCH" == ss(str, 0, 3) {
				return "T", "T", 2
			} else {
				return "0", "T", 2
			}
		} else {
			c := 1
			if mm(ss(str, pos + 1, 1), []string{"T", "D"}) {
				c = 2
			}
			return "T", "T", c
		}
	case "V":
		c := 1
		if "V" == ss(str, pos + 1, 1) {
			c = 2
		}
		return "F", "F", c
	case "W":
		if "R" == ss(str, pos + 1, 1) {
			return "R", "R", 2
		}
		p, s := "", ""
		if 0 == pos && (vowel(ss(str, pos + 1, 1)) || "H" == ss(str, pos + 1, 1)) {
			p = "A"
			if vowel(ss(str, pos + 1, 1)) {
				s = "F"
			} else {
				s = "A"
			}
		}
		if (str.RuneCount() - 1 == pos && vowel(ss(str, pos - 1, 1))) ||
				"SCH" == ss(str, 0, 3) ||
				mm(ss(str, pos - 1, 5), []string{"EWSKI", "EWSKY", "OWSKI", "OWSKY"}) {
			return p, s + "F", 1
		} else if mm(ss(str, pos, 4), []string{"WICZ", "WITZ"}) {
			return p + "TS", s + "FX", 4
		} else {
			return p, s, 1
		}
	case "X":
		c := 1
		if mm(ss(str, pos + 1, 1), []string{"C", "X"}) {
			c = 2
		}
		if str.RuneCount() - 1 == pos && (mm(ss(str, pos - 3, 3), []string{"IAU", "EAU"}) ||
				mm(ss(str, pos - 2, 2), []string{"AU", "OU"})) {
			return "KS", "KS", c
		} else {
			return "", "", c
		}
	case "Z":
		if "H" == ss(str, pos + 1, 1) {
			return "J", "J", 2
		} else {
			c := 1
			if "Z" == ss(str, pos + 1, 1) {
				c = 2
			}
			if mm(ss(str, pos + 1, 2), []string{"ZO", "ZI", "ZA"}) || slavoGermanic(str) && (pos > 0 &&
					"T" != ss(str, pos - 1, 1)) {
				return "S", "TS", c
			} else {
				return "S", "S", c
			}
		}
	}
	return "", "", 1
}
