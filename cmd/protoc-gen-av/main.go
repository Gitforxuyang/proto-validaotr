package main

import (
	"fmt"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/gen"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/generator"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"go/format"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	gen.Init()

	g := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())
	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GenerateAllFiles()
	g.Bytes()
	if err := goformat(g.Response); err != nil {
		g.Error(err)
	}
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
	log.Println("protoc-gen-av 编译proto完成")
}

func goformat(resp *plugin.CodeGeneratorResponse) error {
	for i := 0; i < len(resp.File); i++ {
		formatted, err := format.Source([]byte(resp.File[i].GetContent()))
		if err != nil {
			return fmt.Errorf("go format error: %v", err)
		}
		fmts := string(formatted)
		resp.File[i].Content = &fmts
	}
	return nil
}
