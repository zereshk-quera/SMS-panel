package utils

import (
	"strconv"
	"strings"
)

func GenerateRegex(word string) string {
	word = strings.ToLower(word)
	letterReplacement := make(map[string][]string)
	for i := 97; i <= 122; i++ {
		letterReplacement[strconv.Itoa(i)] = make([]string, 0)
	}
	letterReplacement["a"] = append(letterReplacement["a"], "A", "a", "@", "^", "4")
	letterReplacement["b"] = append(letterReplacement["b"], "B", "b", "8", "ß")
	letterReplacement["c"] = append(letterReplacement["c"], "C", "c", "(", "<", "©", "¢")
	letterReplacement["d"] = append(letterReplacement["d"], "D", "d", "ð")
	letterReplacement["e"] = append(letterReplacement["e"], "E", "e", "∑", "3", "€")
	letterReplacement["f"] = append(letterReplacement["f"], "F", "f", "ƒ", "ʄ")
	letterReplacement["g"] = append(letterReplacement["g"], "G", "g", "ʛ", "6", "9", "ǥ")
	letterReplacement["h"] = append(letterReplacement["h"], "H", "h", "#", "ʜ", "ɦ", "λ")
	letterReplacement["i"] = append(letterReplacement["i"], "I", "i", "1", "!", "|", "ɨ")
	letterReplacement["j"] = append(letterReplacement["j"], "J", "j", "]", "ʝ")
	letterReplacement["k"] = append(letterReplacement["k"], "K", "k", "κ")
	letterReplacement["l"] = append(letterReplacement["l"], "L", "l", "£", "ʟ", "|", "1")
	letterReplacement["m"] = append(letterReplacement["m"], "M", "m")
	letterReplacement["n"] = append(letterReplacement["n"], "N", "n")
	letterReplacement["o"] = append(letterReplacement["o"], "O", "o", "°", "0")
	letterReplacement["p"] = append(letterReplacement["p"], "P", "p")
	letterReplacement["q"] = append(letterReplacement["q"], "Q", "q", "9", "¶")
	letterReplacement["r"] = append(letterReplacement["r"], "R", "r", "®")
	letterReplacement["s"] = append(letterReplacement["s"], "S", "s", "$", "§", "5", "Ƨ")
	letterReplacement["t"] = append(letterReplacement["t"], "T", "t", "+", "7")
	letterReplacement["u"] = append(letterReplacement["u"], "U", "u", "μ", "v")
	letterReplacement["v"] = append(letterReplacement["v"], "V", "v", "u", "√")
	letterReplacement["w"] = append(letterReplacement["w"], "W", "w", "₩")
	letterReplacement["x"] = append(letterReplacement["x"], "X", "x", "%")
	letterReplacement["y"] = append(letterReplacement["y"], "Y", "y", "¥", "γ")
	letterReplacement["z"] = append(letterReplacement["z"], "Z", "z", "2")

	reg := ""
	for i, l := range word {
		reg += "["
		for _, v := range letterReplacement[string(l)] {
			reg += v
		}
		reg += "]"
		if i != len(word)-1 {
			reg += "[^a-zA-Z]*"
		}
	}
	return reg
}
