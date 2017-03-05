// +build !testing

package stemmer

import (
	"regexp"
	"strings"
)

type Stemmer struct{}

var rootWords map[string]int

func init() {
	InitRootWords()
}

func InitRootWords() {
	rootWords = make(map[string]int)

	words := strings.Split(string(data), " ")
	for k := range words {
		rootWords[words[k]] = k
	}
}

func New() *Stemmer {
	return &Stemmer{}
}

func (s *Stemmer) Stemm(ws ...string) []string {
	result := []string{}
	for _, w := range ws {
		word := []byte(strings.ToLower(w))
		if s.IsRootWord(word) {
			result = append(result, w)
		} else {
			result = append(result, s.removingProcess(word))
		}
	}
	return result
}

func (s *Stemmer) IsRootWord(word []byte) bool {
	if _, ok := rootWords[string(word)]; ok {
		return true
	}
	return false
}

func (s *Stemmer) removingProcess(word []byte) string {
	// try to only remove derivation prefixes
	p0 := s.removeDerivationPrefixes(word)
	if s.IsRootWord(p0) {
		return string(p0)
	}

	p1 := s.removeInflectionSuffixes(word)
	p2 := s.removeDerivationSuffixes(p1)
	if s.IsRootWord(p2) {
		return string(p2)
	}
	p3 := s.removeDerivationPrefixes(p2)
	p4 := s.removeDerivationPeople(p3)
	return string(p4)
}

// isRulePrecedence checks the Rule Precedence
// combination of Prefix and Suffix
// "be-lah" "be-an" "me-i" "di-i" "pe-i" or "te-i"
func (s *Stemmer) isRulePrecedence(word []byte) bool {
	if match, _ := regexp.Match(`^(be)([a-z\-]+)(lah|an)$`, word); match {
		return match
	} else if match, _ := regexp.Match(`^(di|[mpt]e)([a-z\-]+)(i)$`, word); match {
		return match
	}

	return false
}

// isdisallowedprefixsuffixes checks Disallowed Prefix-Suffix Combinations
// "be-i" . "di-an" . "ke-i|kan" . "me-an" . "se-i|kan" or "te-an"
func (s *Stemmer) isDisallowedPrefixSuffixes(word []byte) bool {
	if match, _ := regexp.Match(`^(be)([a-z\-]+)(i)$`, word); match {
		return true
	} else if match, _ := regexp.Match(`^(di)([a-z\-]+)(an)$`, word); match {
		return true
	} else if match, _ := regexp.Match(`^(ke)([a-z\-]+)(i|kan)$`, word); match {
		return true
	} else if match, _ := regexp.Match(`^(me)([a-z\-]+)(an)$`, word); match {
		return true
	} else if match, _ := regexp.Match(`^(se)([a-z\-]+)(i|kan)$`, word); match {
		return true
	} else if match, _ := regexp.Match(`^(te)([a-z\-]+)(an)$`, word); match {
		return true
	}
	return false
}

// removeInflectionSuffixes
// 1. Particle "-lah" "-kah" "-tah" and "-pun"
// 2. Possesive Pronoun "-ku" "-mu" "-nya"
func (s *Stemmer) removeInflectionSuffixes(word []byte) []byte {
	if match, _ := regexp.Match(`([klt]ah|pun|[km]u|nya)$`, word); match {
		re, _ := regexp.Compile(`([klt]ah|pun|[km]u|nya)$`)
		infSuf := re.ReplaceAll(word, []byte([]byte("")))

		if match, _ := regexp.Match(`([km]u|nya)$`, infSuf); match {
			re, _ := regexp.Compile(`([km]u|nya)$`)
			posPron := re.ReplaceAll(infSuf, []byte(""))
			return posPron
		}

		return infSuf
	}

	return word
}

// removeDerivationSuffixes
// "-i" . "-kan" . "-an"
func (s *Stemmer) removeDerivationSuffixes(word []byte) []byte {
	var base []byte

	if match, _ := regexp.Match(`(kan)$`, word); match {
		re, _ := regexp.Compile(`(kan)$`)
		base = re.ReplaceAll(word, []byte([]byte("")))
		if s.IsRootWord(base) {
			return base
		}
	}

	if match, _ := regexp.Match(`(an|i)$`, word); match {
		re, _ := regexp.Compile(`(an|i)$`)
		base = re.ReplaceAll(word, []byte([]byte("")))
		if s.IsRootWord(base) {
			return base
		}
	}

	if s.isDisallowedPrefixSuffixes(base) {
		return word
	}

	return word
}

// removeDerivationPeople
// "-man" . "-wan" . "-wati"
func (s *Stemmer) removeDerivationPeople(word []byte) []byte {
	var base []byte

	if match, _ := regexp.Match(`([mw]an)$`, word); match {
		re, _ := regexp.Compile(`([mw]an)$`)
		base = re.ReplaceAll(word, []byte([]byte("")))
		if s.IsRootWord(base) {
			return base
		}
	}

	if match, _ := regexp.Match(`(wati)$`, word); match {
		re, _ := regexp.Compile(`(wati)$`)
		base = re.ReplaceAll(word, []byte([]byte("")))
		if s.IsRootWord(base) {
			return base
		}
	}

	if s.isDisallowedPrefixSuffixes(base) {
		return word
	}

	return word
}

// removeDerivationPrefixes
// "di-" . "ke-" . "se-" . "me-" . "be-" . "pe-" or "te-"
func (s *Stemmer) removeDerivationPrefixes(word []byte) []byte {
	base := word
	var strRemovedDerivSuff []byte
	var strRemovedStdPref []byte
	var strRemovedCmplxPref []byte
	var prefixes []byte

	var re *regexp.Regexp

	if match, _ := regexp.Match(`^(di|[ks]e)\S{1,}`, word); match {
		re, _ = regexp.Compile("^(di|[ks]e)")

		strRemovedStdPref = re.ReplaceAll(word, []byte(""))
		if s.IsRootWord(strRemovedStdPref) {
			return strRemovedStdPref
		}

		strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedStdPref)
		if s.IsRootWord(strRemovedDerivSuff) {
			return strRemovedDerivSuff
		}

		if match, _ := regexp.Match(`^(keber)\S{0,}`, word); match {
			re, _ = regexp.Compile("^([ks]eber)")
			strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
			if s.IsRootWord(strRemovedCmplxPref) {
				return strRemovedCmplxPref
			}

			strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
			if s.IsRootWord(strRemovedDerivSuff) {
				return strRemovedDerivSuff
			}
		}

	} else if match, _ := regexp.Match(`^([^aiueo])e\1[aiueo]\S{1,}`, word); match {
		re, _ = regexp.Compile("^([^aiueo])e")

		prefixes = re.ReplaceAll(word, []byte(""))
		if s.IsRootWord(prefixes) {
			return prefixes
		}

		strRemovedDerivSuff = s.removeDerivationSuffixes(prefixes)
		if s.IsRootWord(strRemovedDerivSuff) {
			return strRemovedDerivSuff
		}

	} else if match, _ := regexp.Match(`^([tmbp]e)\S{1,}`, word); match {
		if match, _ := regexp.Match(`^(be)\S{1,}`, word); match {
			if match, _ := regexp.Match(`^(ber)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ber)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("r"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(ber)[^aiueor]([a-z\-]+)\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ber)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(ber)[^aiueor]([a-z\-]+)er[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ber)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^belajar\S{0,}`, word); match {
				re, _ = regexp.Compile("^(bel)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(be)[^aiueolr]er[^aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(be)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
			}

		}

		if match, _ := regexp.Match(`^(te)\S{1,}`, word); match {
			if match, _ := regexp.Match(`^(terr)\S{1,}`, word); match {
				return word

			} else if match, _ := regexp.Match(`^(ter)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ter)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}
				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
				strRemovedCmplxPref = re.ReplaceAll(word, []byte("r"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}
				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(ter)[^aiueor]er[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ter)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(ter)[^aiueor]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ter)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(te)[^aiueor]er\S{1,}`, word); match {
				re, _ = regexp.Compile("^(te)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(ter)[^aiueor]er[^aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(ter)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
			}
		}

		if match, _ := regexp.Match(`^(me)\S{1,}`, word); match {
			if match, _ := regexp.Match(`^(me)[lrwyv][aiueo]`, word); match {
				re, _ = regexp.Compile("^(me)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(mem)[bfvp]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(mem)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				if match, _ := regexp.Match(`^(member)\S{0,}`, word); match {
					re, _ = regexp.Compile("^(member)")

					strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
					if s.IsRootWord(strRemovedCmplxPref) {
						return strRemovedCmplxPref
					}

					strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
					if s.IsRootWord(strRemovedDerivSuff) {
						return strRemovedDerivSuff
					}
				}

			} else if match, _ := regexp.Match(`^(mem)((r[aiueo])|[aiueo])\S{1,}`, word); match {
				re, _ = regexp.Compile("^(mem)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("m"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("p"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(men)[cdjszt]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(men)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(men)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(men)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("n"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("t"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(meng)[ghqk]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(meng)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(meng)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(meng)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("k"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("ng"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				re, _ = regexp.Compile("^(menge)")
				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(meny)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(meny)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("s"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				re, _ = regexp.Compile("^(me)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
			}
		}

		if len(strRemovedDerivSuff) != 0 {
			word = strRemovedDerivSuff
		}

		if match, _ := regexp.Match(`^(pe)\S{1,}`, word); match {
			if match, _ := regexp.Match(`^(pe)[wy]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pe)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(per)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(per)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("r"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(per)[^aiueor]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(per)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				if match, _ := regexp.Match(`^(perse)\S{0,}`, word); match {
					re, _ = regexp.Compile("^(perse)")

					strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
					if s.IsRootWord(strRemovedCmplxPref) {
						return strRemovedCmplxPref
					}

					strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
					if s.IsRootWord(strRemovedDerivSuff) {
						return strRemovedDerivSuff
					}
				}

			} else if match, _ := regexp.Match(`^(per)[^aiueor]([a-z\-]+)(er)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(per)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pem)[bfv]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pem)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pem)(r[aiueo]|[aiueo])\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pem)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("m"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("p"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pen)[cdjzts]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pen)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pen)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pen)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("n"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("t"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
			} else if match, _ := regexp.Match(`^(peng)[ghq]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(peng)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(peng)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(peng)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("k"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				re, _ = regexp.Compile("^(penge)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(peng)[^ghq]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(peng)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(peny)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(peny)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("s"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				re, _ = regexp.Compile("^(pe)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pel)[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pel)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte("l"))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

				if match, _ := regexp.Match(`^(pelajar)i*\S{0,}`, word); match {
					re, _ = regexp.Compile("^(pel)")

					strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
					if s.IsRootWord(strRemovedCmplxPref) {
						return strRemovedCmplxPref
					}

					strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
					if s.IsRootWord(strRemovedDerivSuff) {
						return strRemovedDerivSuff
					}
				}

			} else if match, _ := regexp.Match(`^(pe)[^rwylmn]er[aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pe)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pe)[^rwylmn]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pe)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}

			} else if match, _ := regexp.Match(`^(pe)[^aiueor]er[^aiueo]\S{1,}`, word); match {
				re, _ = regexp.Compile("^(pe)")

				strRemovedCmplxPref = re.ReplaceAll(word, []byte(""))
				if s.IsRootWord(strRemovedCmplxPref) {
					return strRemovedCmplxPref
				}

				strRemovedDerivSuff = s.removeDerivationSuffixes(strRemovedCmplxPref)
				if s.IsRootWord(strRemovedDerivSuff) {
					return strRemovedDerivSuff
				}
			}
		}
	}

	return base
}
