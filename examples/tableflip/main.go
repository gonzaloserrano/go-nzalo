package main

import (
	"fmt"
	"os"
	"strings"
)

var conversionTable = map[string]string{
	string("a"):      string("\u0250"),
	string("b"):      string("q"),
	string("c"):      string("\u0254"),
	string("d"):      string("p"),
	string("e"):      string("\u01DD"),
	string("f"):      string("\u025F"),
	string("g"):      string("b"),
	string("h"):      string("\u0265"),
	string("i"):      string("\u0131"),
	string("j"):      string("\u0638"),
	string("k"):      string("\u029E"),
	string("l"):      string("\u05DF"),
	string("m"):      string("\u026F"),
	string("n"):      string("u"),
	string("o"):      string("o"),
	string("p"):      string("d"),
	string("q"):      string("b"),
	string("r"):      string("\u0279"),
	string("s"):      string("s"),
	string("t"):      string("\u0287"),
	string("u"):      string("n"),
	string("v"):      string("\u028C"),
	string("w"):      string("\u028D"),
	string("x"):      string("x"),
	string("y"):      string("\u028E"),
	string("z"):      string("z"),
	string("["):      string("]"),
	string("]"):      string("["),
	string(""):       string(""),
	string("{"):      string("}"),
	string("}"):      string("{"),
	string("?"):      string("\u00BF"),
	string("\u00BF"): string("?"),
	string("!"):      string("\u00A1"),
	string("\""):     string(","),
	string(","):      string("\""),
	string("."):      string("\u02D9"),
	string("_"):      string("\u203E"),
	string(";"):      string("\u061B"),
	string("9"):      string("6"),
	string("6"):      string("9"),
}

func main() {
	if len(os.Args) != 2 {
		println("need a word as a first argument")
		os.Exit(1)
	}
	println(Flip(os.Args[1]))
}

func Flip(word string) string {
	word = strings.ToLower(word)
	flipped := ""
	for i := len(word) - 1; i >= 0; i-- {
		flipped += flip(string(word[i]))
	}
	return fmt.Sprintf("(╯°. °）╯︵ ┻%s┻", flipped)
}

func flip(c string) string {
	converted, ok := conversionTable[c]
	if !ok {
		return ""
	}
	return converted
}
