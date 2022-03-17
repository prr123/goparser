//
// goparser
// program that parses a go file and lists functions and types
//
// author prr
// created 28/2/2022
//
// copyright 2022 prr
//
package main

import (
	"os"
	"fmt"
)

func fatErr(fs string, msg string, err error) {
	if err != nil {
		fmt.Printf("error %s:: %s!%v\n", fs, msg, err)
	} else {
		fmt.Printf("error %s:: %s!\n", fs, msg)
	}
	os.Exit(2)
}

func main() {

	numArg := len(os.Args)
	if numArg < 2 {fatErr("main", "insufficient arguments", nil)}

	infilnam := os.Args[1]
	if infilnam[len(infilnam)-3:] != ".go" {
		fatErr("main", "not a go file", nil)
	}

	infil, err := os.Open(infilnam)
	if err != nil {fatErr("Open File", "could not open infil",err)}
	defer infil.Close()

	filInfo, err := infil.Stat()
	if err != nil {fatErr("File Stat", "could get fileinfo",err)}

	nb := filInfo.Size()
	fmt.Printf("file size is %d!\n", nb)

	buf := make([]byte, nb)
	_, err = infil.Read(buf)
	if err != nil {fatErr("File read", "could not read",err)}

	outfilnam := infilnam[:len(infilnam) - 3]+ ".gdat"
	fmt.Printf("out file name: %s \n", outfilnam)

	outfil, err := os.Create(outfilnam)
	if err != nil {fatErr("Outpput file", "could not create",err)}
	defer outfil.Close()

	istate := 0
	ilin:= 0
	var i, ist, j int64
	ist = 0
//	outstr := ""

	for i=0; i< nb - 4; i++ {
		switch istate {
		case 0:
		// parsing a line
			pstr := string(buf[i: i+5])
			switch pstr {
				case "type ":
					ist =  i
//					opst = 0
					istate = 2
					i = i+ 4
					break

				case "func ":
					ist =  i
//					opst = 0
					istate = 3
					i = i+ 4
					break

				default:

			}

			if string(buf[i: i+2]) == "//" {
				istate =1
				i= i+2
			}

			// did not find func or type but EOL
			// start parsing next line
			if buf[i] == '\n' {
				ilin++
				istate = 0
//				opst = 0
			}


		case 1:
			// ignoring everything until there is a new line
			if buf[i] == '\n' {
				ilin++
//				opst = 0
				istate = 0
			}
		case 2:
			// detected the word type parsing the remaing line
			if buf[i] == '{' {
				istate = 1
			// parsing the letters between type and {
				for j=ist+5; j< i; j++ {
					if buf[j] == ' ' {continue}
					// if parser found a letter after 'type' it could be a struct
					if (buf[j] >= 'a' && buf[j] <= 'z') || (buf[j] >= 'A' && buf[j] <= 'Z') {
//						opst = 0
						break
					}
				}
				fmt.Printf("line: %4d type:    %s \n", ilin+1, string(buf[ist: i]))
			}
			// maybe we should also parse for struct


		case 3:
			// detected the word func parsing the remaining line
			// func has 'func name () {' pattern
			// method has 'func (xx) name () {' pattern
			if buf[i] == '{' {
				method := false
				istate = 1
				for j=ist+5; j< i; j++ {
					if buf[j] == ' ' {continue}
					// if parser found a letter after 'func' it must be a function
					if (buf[j] >= 'a' && buf[j] <= 'z') || (buf[j] >= 'A' && buf[j] <= 'Z') {
//						opst = 0
						break
					}
					//  if parser found '(' after 'func' it is a  method
					if buf[j] == '(' {
//						opst = j
						method = true
						break
					}
				}
/*
				// looking for ')'
				if opst > 0 {
					for j= opst+1; j< i; j++ {
						if buf[j] == ')' {
							method = true
							break
						}
					}
				}
*/
				if method {
					fmt.Printf("line: %4d method:  %s \n", ilin+1, string(buf[ist: i]))
				} else {
					fmt.Printf("line: %4d func:    %s \n", ilin+1, string(buf[ist: i]))
				}
				ist =0
			}
		default:

		}

	}

//	fmt.Printf("%s\n", outstr)
	fmt.Printf("total lines: %d\n", ilin)
	fmt.Println("success!")
}
