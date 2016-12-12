package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// The "band" will be like a unilaterally limited turing machine
var (
	bp           int      // Pointer for the band
	pc           int      // Program counter
	band         []int    // unilaterally limited turing machine band  |0|1|-1|2|-2|3|-3|...|
	loopStack    []int    // saves points of opening '[' to jump again
	instructions []string // in case brainfuck get's an expansion ;-)
)

func initBand() {
	bp = 0
	pc = 0
	band = make([]int, 50)
	loopStack = make([]int, 0)
	instructions = []string{"<", ">", "[", "]", "+", "-", ".", ","}
}

func contains(arr []string, str string) bool {
	for _, val := range arr {
		if str == val {
			return true
		}
	}
	return false
}

func parseArgs() int {
	switch foo := len(os.Args); {
	case foo > 2:
		log.Fatal("Usage: ./gofuck [file]")
	case foo == 2:
		_, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return 2
	default:
		return 1
	}
	return 0
}

func onlyInstr(text string) string {
	ret := ""
	for _, val := range text {
		if contains(instructions, string(val)) {
			ret += string(val)
		}
	}
	return ret
}

func convStdToStr() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your program. End it with a ';'\n")
	text, _ := reader.ReadString(';')
	return onlyInstr(text)
}

func handleStdin() {
	program := convStdToStr()
	execute(program)
}

func convFileToStr(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return onlyInstr(string(content))
}

func handleFile(filename string) {
	program := convFileToStr(filename)
	execute(program)
}

func checkSyntax(program string) {
	count := 0
	for key, val := range program {
		if string(val) == "[" {
			count++
		} else if string(val) == "]" {
			count--
		}
		if count < 0 {
			log.Fatal("Wrong parenthesis at character ", key+1)
		}
	}
	if count > 0 {
		log.Fatal("Too many or missing parenthesis!")
	}
}

func execute(program string) {
	checkSyntax(program)
	for pc < len(program) {
		switch instr := program[pc]; {
		case string(instr) == "<":
			if bp == 1 {
				bp--
			} else if bp%2 == 1 {
				bp -= 2
			} else {
				bp += 2
			}
			if bp > len(band) {
				newBand := make([]int, len(band)*2)
				band = append(newBand, band[:]...)
			}
			pc++
		case string(instr) == ">":
			if bp == 0 {
				bp++
			} else if bp%2 == 1 {
				bp += 2
			} else {
				bp -= 2
			}
			if bp > len(band) {
				newBand := make([]int, len(band)*2)
				band = append(newBand, band[:]...)
			}
			pc++
		case string(instr) == "[":
			loopStack = append(loopStack, pc)
			pc++
		case string(instr) == "]":
			if band[bp] != 0 {
				pc = loopStack[len(loopStack)-1] + 1
			} else {
				loopStack = loopStack[:len(loopStack)-1] // basiclally popping the last element
				pc++
			}
		case string(instr) == "+":
			band[bp]++
			pc++
		case string(instr) == "-":
			band[bp]--
			pc++
		case string(instr) == ".":
			fmt.Printf("%c", band[bp])
			pc++
		case string(instr) == ",":
			var i rune
			_, err := fmt.Scanf("%c", &i)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(i)
			pc++
		}
	}
}

func main() {
	initBand()
	switch args := parseArgs(); {
	case args == 1:
		handleStdin()
	case args == 2:
		handleFile(os.Args[1])
	default:
		log.Fatal("There was a bad mistake")
	}
}
