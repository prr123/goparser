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
	var i, ist, j, opst, linst int64
	linst =0
	ist = 0
//	outstr := ""
	for i=0; i< nb - 4; i++ {
		switch istate {
		case 0:
			if string(buf[i: i+4]) == "func" {
			// need to find function name
//				fmt.Printf("dbg: %s line: %d\n", string(buf[i: i+10]), ilin)
					ist =  i
					opst = 0
					istate = 3
					i = i+ 4
			} else {
				if buf[i] == '\n' {
					ilin++
					linst = i+1
					istate = 0
					opst = 0
				} else {
					istate = 1
				}
			}
		case 1:
			if buf[i] == '\n' {
				linst = i+1
				ilin++
				opst = 0
				istate = 0
			}
		case 2:
			if string(buf[i: i+2]) == "//" {
				istate =1
				i= i+2
			}

		case 3:
			if buf[i] == '{' {
				method := false
				istate = 1
				if linst < ist {
					for j=linst; j< ist; j++ {
						if buf[j] == '(' {
							opst = j
							break
						}
					}
					if opst > 0 {
						for j= opst+1; j< ist; j++ {
							if buf[j] == ')' {
								method = true
								break
							}
						}
					}
				}
				if method {
					fmt.Printf("line %4d method found: %s \n", ilin+1, string(buf[opst: i]))
				} else {
					fmt.Printf("line: %4d funk found: %s \n", ilin+1, string(buf[ist: i]))
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
