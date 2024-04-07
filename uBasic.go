package main

import (
	"log"
	"os"

	"uBasic/exec"
	"uBasic/ide"
	"uBasic/irgen"

	"gioui.org/app"
)

// to call rtllib as a dynamic library
// https://www.ardanlabs.com/blog/2013/08/using-c-dynamic-libraries-in-go-programs.html
// remove the first level of comment to use the following code
//
// /* ----- add # before each three lines below -------
// cgo CFLAGS: -I../DyLib
// cgo LDFLAGS: -L. -lname of the library
// include <name of the library.h>
// */
// import "C"
// functions from the library can be called using C.functionName()
// than read:
// 		http://golang.org/cmd/cgo/
// 		http://golang.org/doc/articles/c_go_cgo.html

func main() {
	var err error
	if len(os.Args) > 1 {
		err = exec.LoadInterpreter(os.Args[1])
	} else {
		// err = exec.LoadInterpreter("testdata/noisy/simple/sim07.bas")
		// err = exec.LoadInterpreter("testdata/noisy/advanced/primes.bas")
		// err = exec.LoadInterpreter("testdata/samples/error01.bas")
		// err = exec.LoadInterpreter("testdata/test color.bas")
		// err = exec.LoadInterpreter("testdata/incorrect/parser/pe17.bas")
		err = exec.LoadInterpreter("testdata/compile.bas")

	}
	if err != nil {
		log.Fatal(err)
	}

	if err := irgen.GenToFile(exec.Debug.File, exec.Debug.Info, "irgen/llvm/compile.ll"); err != nil {
		log.Fatal(err)
	}

	go func() {
		w := app.NewWindow()
		err := ide.Run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	go func() {
		if err := exec.Debug.Run(); err != nil {
			exec.Debug.Terminal.WriteString(err.Error())
		}
	}()
	app.Main()
}
