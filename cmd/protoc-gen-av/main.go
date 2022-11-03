package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"go/format"
	"io/ioutil"
	"log"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/gen"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/generator"
	"net/http"
	"os"
)

var (
	version = "1.7.3"
)

func main() {
	gen.Init()
	resp, err := http.Get("http://plugin-check.neoclub.cn/v1/plugin/version")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	plugin := map[string]string{}
	err = json.Unmarshal(body, &plugin)
	g := generator.New()
	log.Println("protoc-gen-meta 组件更新检查开始")
	if err != nil {
		log.Println("warning: protoc-gen-meta 组件更新检查失败")
	} else {
		log.Printf("protoc-gen-meta 组件当前版本:%s 组件最新版本：%s", version, plugin["now"])
	}
	if plugin["min"] > version {
		g.Fail("\n 当前插件已经过于老旧，请马上更新，安装命令：\n go install meta/scripts/cmd/protoc-gen-meta \n")
	}
	if plugin["rec"] > version {
		g.Fail("\n 当前插件已经过于老旧，推荐更新，安装命令：\n go install meta/scripts/cmd/protoc-gen-meta \n")
	}

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
	log.Println("protoc-gen-meta 编译proto完成")
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
