package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

func main() {
	fmt.Println("CalDAV fixer for missing UID.")
	inPtr := flag.String("f", "", "the input filename")
	outPtr := flag.String("o", "", "the input filename")
	flag.Parse()
	if *inPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	in := *inPtr
	var out string
	if outPtr == nil || *outPtr == "" {
		fmt.Println(*outPtr)
		out = fmt.Sprintf("%s-%s.ics", path.Base(*inPtr), time.Now())
	} else {
		fmt.Println(*outPtr)
		out = *outPtr
	}
	fmt.Println("Fixing input file: ", in)
	fmt.Println("Writing to output: ", out)
	inF, err := os.Open(in)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outF, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer inF.Close()
	defer outF.Close()

	bufR := bufio.NewScanner(bufio.NewReader(inF))
	bufW := bufio.NewWriter(outF)
	for bufR.Scan() {
		line := bufR.Text()
		bufW.WriteString(line)
		bufW.WriteString("\n")
		if strings.Contains(line, "SUMMARY") {
			uid, err := uuid.NewV4()
			if err != nil {
				fmt.Println(err)
			}
			uidLine := fmt.Sprintf("UID:%s", uid)
			bufW.WriteString(uidLine)
			bufW.WriteString("\n")
		}
	}
	bufW.Flush()
}
