package main
//
// goparser
// program that parses a go file and lists functions and types
//
// author prr
// created 28/2/2022
//
// copyright 2022 prr
//

import (
	"os"
	"fmt"
)

type CmdOpt struct {
    Verb bool
}

func isAlpha(let byte)(res bool) {
// function that test whether byte is alpha
    res = false
    if (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z') { res = true}
    return res
}

func isAlphaNumeric(let byte)(res bool) {
// function that test whether byte is aphanumeric
    res = false
    tbool := (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z')
    if tbool || (let >= '0' && let <= '9') { res = true }
    return res
}

func fatErr(fs string, msg string, err error) {
// function that displays a console error message and exits program
	if err != nil {
		fmt.Printf("error %s:: %s!%v\n", fs, msg, err)
	} else {
		fmt.Printf("error %s:: %s!\n", fs, msg)
	}
	os.Exit(2)
}

func main() {
// main program
// cmd line: goparser file [-v]
// program creates a consol log of all: type, func, method
    var opt CmdOpt
	var i, ist int64

	numArg := len(os.Args)
	if numArg < 2 {fatErr("main", "insufficient arguments", nil)}

	infilnam := os.Args[1]
	if infilnam[len(infilnam)-3:] != ".go" {
		fatErr("main", "not a go file", nil)
	}

    opt.Verb = false
    if numArg > 2 {
        if os.Args[2] == "-v" {
            opt.Verb = true
        }
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

//	outfilnam := infilnam[:len(infilnam) - 3]+ ".gdat"
//	fmt.Printf("out file name: %s \n", outfilnam)
//	outfil, err := os.Create(outfilnam)
//	if err != nil {fatErr("Output file", "could not create",err)}
//	defer outfil.Close()

	istate := 0
	ilin:= 0
//	linst =0
	ist = 0
//	outstr := ""
	for i=0; i< nb - 4; i++ {
		switch istate {
		case 0:
           // new line no history
			if string(buf[i: i+5]) == "func " {
//				fmt.Printf("dbg: %s line: %d\n", string(buf[i: i+10]), ilin)
                // remember start
				ist =  i
				istate = 3
				i = ist + 5
			}
			if string(buf[i: i+5]) == "func(" {
            // must be a method
                istate = 4
                ist = i
                i = ist + 5
            }
            if string(buf[i:i+5]) == "type " {
                istate = 11
                ist = i
                i = ist + 5
            }
			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i
			}

		case 1:
        // found comment looking for EOL
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i + 1
			}
		case 2:
        // found comment looking for EOL
			if buf[i] == '\n' {
                fmt.Printf("line: %4d comment: %s\n", ilin, string(buf[ist:i-1]))
				istate = 0
				ilin++
                ist = i+1
			}

		case 3:
        // found a func  now we need to decide whether it is a method or a function
            // need to parse string between func and {
			if buf[i] == ' ' {break}
            // if found alpha char first then is a func
   			if isAlpha(buf[i]) {
                istate = 8
                break
            }
            // if we find a ( first then it is a method
			if buf[i] == '(' {
                istate = 4
                break
			}
            istate = 10
        case 4:
        // found ( before alpha must be a method looking for closing par and subsequent string
            if isAlpha(buf[i]) {
//                fmt.Printf("method type: %s\n", string(buf[ist: i]))
                istate = 5
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 5:
        // method start alpha now white space and (
            if buf[i] == ' ' {
                istate = 6
            }
            if buf[i] == '(' {
                istate = 6
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i + 1
			}

        case 6:
        // method has to have parameter parenthesis
           if buf[i] == '(' {
                istate = 7
           }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 7:
        // found method need to find {
			if buf[i] == '{' {
				fmt.Printf("line: %4d method found:   %s \n", ilin, string(buf[ist: i]))
                istate = 10
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 8:
        // func found alpha string before
           if buf[i] == '(' {
                istate = 9
           }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 9:
        // function need to find { next and thereafter EOL
			if buf[i] == '{' {
                fmt.Printf("line: %4d function found: %s \n", ilin, string(buf[ist: i]))
                istate = 10            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 10:
        // method and function found { now need to find EOL
           if buf[i] == '\n' {
				ilin++
                istate = 20
                ist = i+1
           }

        case 11:
        // found type need to find alpha
            if isAlpha(buf[i]) {
                istate = 12
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}
        case 12:
        // found alpha need to find whitespace
            if buf[i] == ' ' {
                istate = 13
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}
        case 13:
        // found whitespace need to find word
            if isAlpha(buf[i]) {
                istate = 14
            }
   			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

        case 14:
   			if buf[i] == ' ' {
				istate = 15
			}
   			if buf[i] == '{' {
				ilin++
				istate = 15
			}
            // this is a type
   			if buf[i] == '\n' {
                fmt.Printf("line: %4d type found: %s \n", ilin, string(buf[ist: i]))
				ilin++
				istate = 0
                ist = i+1
			}

        case 15:
            if isAlpha(buf[i]) {
                istate = 16
            }
   			if buf[i] == '{' {
				istate = 16
			}

        case 16:
        // found type
   			if buf[i] == '\n' {
                fmt.Printf("line: %4d type found: %s \n", ilin, string(buf[ist: i]))
				ilin++
				istate = 20
                ist = i+1
			}

        case 20:
           // new line no history
			if string(buf[i: i+5]) == "func " {
//				fmt.Printf("dbg: %s line: %d\n", string(buf[i: i+10]), ilin)
                // remember start
				ist =  i
				istate = 3
				i = ist + 5
			}
			if string(buf[i: i+5]) == "func(" {
            // must be a method
                istate = 4
                ist = i
                i = ist + 5
            }
            if string(buf[i:i+5]) == "type " {
                istate = 11
                ist = i
                i = ist + 5
            }
			if buf[i] == '\n' {
				ilin++
				istate = 0
                ist = i+1
			}

            if string(buf[i:i+2]) == "//" {
				istate = 2
                ist = i
                i = ist + 2
			}

		default:

		}

	}

//	fmt.Printf("%s\n", outstr)
	fmt.Printf("total lines: %d\n", ilin)
	fmt.Println("success!")
}
