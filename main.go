package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/emicklei/proto"
)

func main() {
	var rootDir, filePath string
	flag.StringVar(&rootDir, "dir", "", "root directory to process")
	flag.StringVar(&filePath, "file", "", "file to process")
	flag.Parse()

	if rootDir == "" && filePath == "" {
		fmt.Println("-dir or -file flag is required")
		os.Exit(1)
	}

	if rootDir != "" {
		_, err := os.Stat(rootDir)
		if os.IsNotExist(err) {
			fmt.Printf("directory %s does not exist\n", rootDir)
			os.Exit(1)
		}
	}

	var fileList []string
	if filePath != "" {
		fileList = append(fileList, filePath)
	}

	if rootDir != "" {
		protoFiles := resolveProtoFilesForDir(rootDir)
		fileList = append(fileList, protoFiles...)
	}

	var totalTypes int

	fmt.Printf("processing %d files\n", len(fileList))

	for _, path := range fileList {
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		totalTypes = totalTypes + countFileTypes(f)
	}

	fmt.Printf("total types: %d\n", totalTypes)
	estimatedMonthlyPrice := float64(totalTypes * 5)
	fmt.Printf("estimated monthly list price: $%f\n", estimatedMonthlyPrice)
	fmt.Printf("estimated monthly price @ 25%%: $%f \n", estimatedMonthlyPrice*.75)
	fmt.Printf("estimated monthly price @ 50%%: $%f\n", estimatedMonthlyPrice*.50)
	fmt.Printf("estimated monthly price @ 75%%: $%f\n", estimatedMonthlyPrice*.25)
}

func resolveProtoFilesForDir(dir string) []string {
	var fileList []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".proto" {
			fileList = append(fileList, path)
		}
		return nil
	}

	filepath.WalkDir(dir, walkFunc)
	return fileList
}

func countFileTypes(rdr io.Reader) int {
	parser := proto.NewParser(rdr)
	p, err := parser.Parse()
	if err != nil {
		panic(err.Error())
	}

	var totalTypes int
	proto.Walk(p,
		proto.WithMessage(func(m *proto.Message) {
			totalTypes = totalTypes + 1
		}),
		proto.WithEnum(func(e *proto.Enum) {
			totalTypes = totalTypes + 1
		}),
		proto.WithRPC(func(r *proto.RPC) {
			totalTypes = totalTypes + 1
		}),
	)
	return totalTypes
}
