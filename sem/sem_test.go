package sem

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
	"uBasic/source"
)

type TI struct {
	src string
}

var testData = []*TI{
	{`' this is a comment`},
	{`Print Chr(65)`},
	{`Dim t as integer
	For t = 10 To 0 Step -2
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
			_, err := Check(file)
			if err != nil {
				pass = false
				// display enhance message format
				src := source.Source{Input: testCase.src}
				err.(*errors.Error).Src = &src
				t.Log(err.Error())
				// t.Log(fmt.Sprintf("error in %v at line %d, column %d\n",
				// 	err.(*parserError.Error).ErrorToken.Lit,
				// 	err.(*parserError.Error).ErrorToken.Line,
				// 	err.(*parserError.Error).ErrorToken.Column))

			}
		} else if p.Errors() != nil {
			pass = false
			for _, msg := range p.Errors() {
				t.Log(msg)
			}
		}
	}
	if !pass {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	pass := true
	// read files from samples directory
	files, err := os.ReadDir("../testdata/samples")
	if err != nil {
		t.Log(err)
		pass = false
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".bas" {
			t.Log(f.Name())
			// read file
			data, err := os.ReadFile("../testdata/samples/" + f.Name())
			if err != nil {
				t.Log(err)
				pass = false
			} else {
				l := lexer.New(string(data))
				p := parser.New(l)

				file := p.ParseFile()
				if file != nil {
					// t.Logf("%v\n", result)
					f := func(n ast.Node) error {
						t.Log(n.String())
						return nil
					}
					astutil.Walk(file, f)
				} else {
					src := source.Source{Input: string(data)}
					src.Name = f.Name()
					for _, err := range p.Errors() {
						// enhanced error format
						if err, ok := err.(*errors.Error); ok {
							err.Src = &src
						}
						t.Log(err)
					}
					pass = false
				}
			}
		}
	}
	if !pass {
		t.Fail()
	}
}

var directories = []string{
	// "../testdata/samples",
	// "../testdata/incorrect/parser",
	//"../testdata/incorrect/semantic",
	//"../testdata/incorrect",
	"../testdata/noisy",
	"../testdata/noisy/advanced",
	"../testdata/noisy/medium",
	"../testdata/noisy/simple",
	"../testdata/quiet",
	"../testdata/quiet/lexer",
	"../testdata/quiet/parser",
	"../testdata/quiet/rtl",
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
				// t.Log(filename)
				// read file
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Log(err)
					pass = false
				} else {
					l := lexer.New(string(data))
					p := parser.New(l)

					src := source.Source{Input: string(data), Name: f.Name()}
					file := p.ParseFile()
					if file == nil {
						pass = false
						for _, err := range p.Errors() {
							// enhanced error format
							if err, ok := err.(*errors.Error); ok {
								err.Src = &src
							}
							t.Log(err)
						}
					} else {
						file.Name = filename
						_, err := Check(file)
						if err != nil {
							// enhanced error format
							if err, ok := err.(*errors.Error); ok {
								err.Src = &src
							}
							t.Log(err)
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
// 				// read file
// 				data, err := os.ReadFile(filename)
// 				if err != nil {
// 					t.Log(err)
// 					pass = false
// 				} else {
// 					l := lexer.New(string(data))
// 					p := parser.New(l)

// 					src := source.Source{Input: string(data), Name: f.Name()}
// 					file := p.ParseFile()

// 					// to validate success we only keep the file
// 					// analysed and saved successfully
// 					if file != nil {
// 						file.Name = filename
// 						if _, err := Save(file); err == nil {
// 							file2, err := ast.LoadFile(filename + "x")
// 							if err != nil {
// 								pass = false
// 								// enhanced error format
// 								if err, ok := err.(*errors.Error); ok {
// 									err.Src = &src
// 								}
// 								t.Log(filename + ": " + err.Error())
// 							} else {
// 								res := file.EqualsNode(file2)
// 								if res != nil {
// 									pass = false
// 									t.Log(filename + ": Files are not equal in " + res.String() + " at " + res.Token().Position.String())

// 								} else {
// 									t.Log(filename + ": Files are equal.")
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
