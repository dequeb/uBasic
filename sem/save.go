package sem

import (
	"os"
	"uBasic/ast"
)

func SaveFile(file *ast.File, types infoTypes) error {
	var dict ast.Dict

	// get filename
	filename := file.Name[:len(file.Name)-4] + ".basx" // replace .bas by .basx
	// dict is used to compress file by replacing
	// the types and expressions with a number
	dict = make(ast.Dict)
	// create file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file in compact format
	_, err = f.WriteString(file.WriteCompact(&dict))
	return err
}
