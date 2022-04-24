// parseLib.go
// go library to parse a go file and list all types, constants, functions and methods
// author: prr
// date 24/4/2022
// copyright 2022 prr azul software
//

package parseLib

import (
	"fmt"
	"os"
	util "utility/utilLib"
)

func parseFil(infil, outfil *os.File)(err error) {

	filinfo, err := infil.Stat()
	if err != nil {
		return fmt.Errorf("infil does not exist! %v", err)
	}
	filinfo, err = outfil.Stat()
	if err != nil {
		return fmt.Errorf("outfil does not exist! %v", err)
	}

	return nil
}
