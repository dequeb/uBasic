package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/lexer"
	"uBasic/parser"
	"uBasic/sem"
	"uBasic/source"
)

type TI struct {
	src string
}

var testData = []*TI{
	{`' this is a comment	`},
	{`Print Chr(65)`},
	{`For t = 10 To 0 Step -2
		Debug.Print t
	Next t`},
	{`Dim b As String`},
	{`Dim b As Boolean`},
	{`Dim aa As Long`},
	{`Debug.Print ""`},
}

func PrintNode(n ast.Node) error {
	fmt.Println(n.String())
	return nil
}

func Test1(t *testing.T) {
	pass := true
	for _, testCase := range testData {
		l := lexer.New(testCase.src)
		p := parser.New(l)
		file := p.ParseFile()
		if file != nil {
			f := func(n ast.Node) error {
				t.Log(n.String())
				return nil
			}
			astutil.Walk(file, f)
		}
	}
	if !pass {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	pass := true
	// read files from samples directory
	files, err := os.ReadDir("./testdata/samples")
	if err != nil {
		t.Log(err)
		pass = false
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".bas" {
			t.Log(file.Name())
			// read file
			data, err := os.ReadFile("./testdata/samples/" + file.Name())
			if err != nil {
				t.Log(err)
				pass = false
			} else {
				l := lexer.New(string(data))
				p := parser.New(l)
				file := p.ParseFile()
				if file != nil {
					f := func(n ast.Node) error {
						t.Log(n.String())
						return nil
					}
					astutil.Walk(file, f)
				}
			}
		}
	}
	if !pass {
		t.Fail()
	}
}

var directories = []string{
	"testdata/samples",
	// "testdata/incorrect/parser",
	"testdata/incorrect/semantic",
	// "testdata/incorrect",
	"testdata/noisy",
	"testdata/noisy/advanced",
	"testdata/noisy/medium",
	"testdata/noisy/simple",
	"testdata/quiet",
	// "testdata/quiet/lexer",
	// "testdata/quiet/parser",
	"testdata/quiet/rtl",
}

func Test3(t *testing.T) {
	pass := true
	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			t.Log(err)
			pass = false
		}

		for _, f := range files {
			if filepath.Ext(f.Name()) == ".bas" {
				filename := path.Join(dir, f.Name())
				t.Log(filename)
				// read file
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Log(err)
					pass = false
				} else {
					l := lexer.New(string(data))
					p := parser.New(l)
					file := p.ParseFile()
					if file == nil {
						for _, err := range p.Errors() {
							e := err.(*errors.Error)
							e.Source = &source.Source{Input: string(data), Name: filename}
							t.Log(err.Error())
						}
					} else {
						file.Name = filename
						if _, err := sem.Check(file); err != nil {
							pass = false
							e := err.(*errors.Error)
							e.Source = &source.Source{Input: string(data), Name: filename}
							t.Log(err.Error())
						} else {
							t.Log("File: ", filename, " passed semantic analysis.")
						}
					}
				}
			}
		}
	}
	if !pass {
		t.Fail()
	}
}

// func Test4(t *testing.T) {
// 	pass := true
// 	for _, dir := range directories {
// 		files, err := os.ReadDir(dir)
// 		if err != nil {
// 			t.Log(err)
// 			pass = false
// 		}

// 		for _, f := range files {
// 			if filepath.Ext(f.Name()) == ".bas" {
// 				filename := path.Join(dir, f.Name())
// 				t.Log(filename)
// 				// read file
// 				data, err := os.ReadFile(filename)
// 				if err != nil {
// 					t.Log(err)
// 					pass = false
// 				} else {
// 					s := &source.Source{Input: string(data), Name: filename}
// 					l := lexer.New(string(data))
// 					p := parser.New(l)
// 					file := p.ParseFile()
// 					if file == nil {
// 						for _, err := range p.Errors() {
// 							e := err.(*errors.Error)
// 							e.Src = s
// 							t.Log(err.Error())
// 						}
// 					} else {
// 						file.Name = filename
// 						if _, err := sem.Save(file); err != nil {
// 							pass = false
// 							e := err.(*errors.Error)
// 							e.Src = s
// 							t.Log(err.Error())
// 						} else {
// 							file2, err := ast.LoadFile(filename + "x")
// 							if err != nil {
// 								pass = false
// 								t.Log(err.Error())
// 							} else {

// 								errNode := file.EqualsNode(file2)
// 								if errNode != nil {
// 									pass = false
// 									t.Log("Files are not equal : ", filename, " at ", errNode.Token().Position.String()+" : "+errNode.String())
// 								} else {
// 									t.Log("Files are equal: " + filename)
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	if !pass {
// 		t.Fail()
// 	}
// }
