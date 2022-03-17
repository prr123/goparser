# goparser
parser of go files to identify functions and methods in a go source file!

The purpose of the program is to list all functions and methods with their line locations in a file.
I found this to be useful when creating a larger library with dozens of functions and methods.

usage: goparser inputfile

output:   
  >>file size is 2493!    
  >>out file name: goparser.gdat    
  >>line:   17 funk found: func fatErr(fs string, msg string, err error)    
  line:   26 funk found: func main()    
  total lines: 130    
  success!   
  
  this sample is the output of parsing goparser.go    
  more tests will be added
  
  
