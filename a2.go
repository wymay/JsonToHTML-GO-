package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// Token in json file
type Token struct {
	t int
	v string
}

// define token types
const (
	LeftCurlyBrace = iota
	RightCurlyBrace
	LeftBracket
	RightBracket
	Colon
	Comma
	BoolOrNull
	String
	EscapeCharacters
	Number
)

var colorMap = map[int]string{
	LeftCurlyBrace:   "red",
	RightCurlyBrace:  "red",
	LeftBracket:      "green",
	RightBracket:     "green",
	Colon:            "DimGray",
	Comma:            "Plum",
	BoolOrNull:       "blue",
	String:           "purple",
	EscapeCharacters: "orange",
	Number:           "LightCoral",
}

func jsonScanner(fileName string) []Token {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)
	openQuote := false
	var str string
	var result []Token
	for {
		if c, err := r.ReadByte(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			switch {
			case c == '"' && !openQuote:
				openQuote = true
				str += string(c)
				if c, err = r.ReadByte(); err != nil {
					if err == io.EOF {
						return result
					}
					log.Fatal(err)
				}
				for c != '"' || str[len(str)-1] == '\\' {
					str += string(c)
					if c, err = r.ReadByte(); err != nil {
						if err == io.EOF {
							return result
						}
						log.Fatal(err)
					}
				}
				if c == '"' && openQuote && str[len(str)-1] != '\\' {
					openQuote = false
					str += string(c)
					result = append(result, Token{String, str})
					str = ""
				}
			case c == '{':
				result = append(result, Token{LeftCurlyBrace, "{"})
			case c == '}':
				result = append(result, Token{RightCurlyBrace, "}"})

			case c == '[':
				result = append(result, Token{LeftBracket, "["})

			case c == ']':
				result = append(result, Token{RightBracket, "]"})

			case c == ':':
				result = append(result, Token{Colon, ":"})

			case c == ',':
				result = append(result, Token{Comma, ","})
			case c == 't':
				result = append(result, Token{BoolOrNull, "true"})
				for i := 0; i < 3; i++ {
					if c, err = r.ReadByte(); err != nil {
						if err == io.EOF {
							return result
						}
						log.Fatal(err)
					}
				}
			case c == 'f':
				result = append(result, Token{BoolOrNull, "false"})
				for i := 0; i < 4; i++ {
					if c, err = r.ReadByte(); err != nil {
						if err == io.EOF {
							return result
						}
						log.Fatal(err)
					}
				}
			case c == 'n':
				result = append(result, Token{BoolOrNull, "null"})
				for i := 0; i < 3; i++ {
					if c, err = r.ReadByte(); err != nil {
						if err == io.EOF {
							return result
						}
						log.Fatal(err)
					}
				}
			case (c > '0' && c < '9') || c == 'e' || c == 'E' || c == '.' || c == '-' || c == '+':
				result = append(result, Token{Number, string(c)})

			}
		}
	}
	return result
}

func makeIndent(num int) {
	for i := 0; i < num; i++ {
		fmt.Printf("&nbsp;&nbsp;&nbsp;&nbsp;")
	}
}

func writeToHTML(tokens []Token) {
	fmt.Printf("<span style=\"font-family:monospace; white-space:pre\">")
	indention := 0
	for _, tok := range tokens {
		switch tok.t {
		case LeftCurlyBrace:
			fmt.Printf("<span style=\"color:" + colorMap[LeftCurlyBrace] + "\">")
			fmt.Printf("{\n" + "</span>")
			indention++
			makeIndent(indention)
		case RightCurlyBrace:
			fmt.Printf("<span style=\"color:" + colorMap[RightCurlyBrace] + "\">")
			fmt.Printf("\n")
			indention--
			makeIndent(indention)
			fmt.Printf("}\n" + "</span>")
		case LeftBracket:
			fmt.Printf("<span style=\"color:" + colorMap[LeftBracket] + "\">")
			fmt.Printf("[\n" + "</span>")
			indention++
			makeIndent(indention)
		case RightBracket:
			fmt.Printf("<span style=\"color:" + colorMap[RightBracket] + "\">")
			fmt.Printf("\n")
			indention--
			makeIndent(indention)
			fmt.Printf("]" + "</span>")
		case Colon:
			fmt.Printf("<span style=\"color:" + colorMap[Colon] + "\">")
			fmt.Printf(" : " + "</span>")
		case Comma:
			fmt.Printf("<span style=\"color:" + colorMap[Comma] + "\">")
			fmt.Printf(",\n" + "</span>")
			makeIndent(indention)
		case BoolOrNull:
			fmt.Printf("<span style=\"color:" + colorMap[BoolOrNull] + "\">")
			fmt.Printf(tok.v + "</span>")
		case String:
			tmpstr := tok.v
			var i int
			for i = 0; i < len(tmpstr); i++ {
				if tmpstr[i] == '\\' {
					fmt.Printf("<span style=\"color:" + colorMap[EscapeCharacters] + "\">")
					switch tmpstr[i+1] {
					case 'u':
						fmt.Printf("%s", tmpstr[i:i+6])
						i += 5
					case '"':
						fmt.Printf("\\&quot;")
						i++
					default:
						fmt.Printf("%s", tmpstr[i:i+2])
						i++
					}
					fmt.Printf("</span>")
				} else {
					fmt.Printf("<span style=\"color:" + colorMap[String] + "\">")
					switch tmpstr[i] {
					case '<':
						fmt.Printf("&lt;")
					case '>':
						fmt.Printf("&gt;")
					case '&':
						fmt.Printf("&amp;")
					case '"':
						fmt.Printf("&quot;")
					case '\'':
						fmt.Printf("&apos;")
					default:
						fmt.Printf("%c", tmpstr[i])
					}
					fmt.Printf("</span>")
				}
			}
		case Number:
			fmt.Printf("<span style=\"color:" + colorMap[Number] + "\">")
			fmt.Printf(tok.v + "</span>")
		}

	}
	fmt.Printf("</span>")
}

func main() {
	jsonFileName := os.Args[1]
	tokens := jsonScanner(jsonFileName)
	// for _, tok := range tokens {
	// 	fmt.Printf("%v", tok.v)
	// }
	writeToHTML(tokens)
}
