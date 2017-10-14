package main

import (
	"fmt"
	"os"
	"strings"
)

var inversions = map[string]string{
	"a":      "\u0250",
	"b":      "q",
	"c":      "\u0254",
	"d":      "p",
	"e":      "\u01DD",
	"f":      "\u025F",
	"g":      "b",
	"h":      "\u0265",
	"i":      "\u0131",
	"j":      "\u0638",
	"k":      "\u029E",
	"l":      "\u05DF",
	"m":      "\u026F",
	"n":      "u",
	"o":      "o",
	"p":      "d",
	"q":      "b",
	"r":      "\u0279",
	"s":      "s",
	"t":      "\u0287",
	"u":      "n",
	"v":      "\u028C",
	"w":      "\u028D",
	"x":      "x",
	"y":      "\u028E",
	"z":      "z",
	"[":      "]",
	"]":      "[",
	"":       "",
	"{":      "}",
	"}":      "{",
	"?":      "\u00BF",
	"\u00BF": "?",
	"!":      "\u00A1",
	"\"":     ",",
	",":      "\"",
	".":      "\u02D9",
	"_":      "\u203E",
	";":      "\u061B",
	"9":      "6",
	"6":      "9",
}

func main() {
	if len(os.Args) != 2 {
		println("error: need a word as a first argument")
		os.Exit(1)
	}

	fmt.Printf("(╯°. °）╯︵ ┻%s┻\n", invert(os.Args[1]))
}

func invert(word string) string {
	word = strings.ToLower(word)
	output := ""
	for i := len(word) - 1; i >= 0; i-- {
		inverted, _ := inversions[string(word[i])]
		output += inverted
	}
	return output
}
