package gen

import (
	"errors"
	"fmt"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/generator"
	"github.com/Gitforxuyang/proto-validaotr/proto/plugin"
	"github.com/Gitforxuyang/proto-validaotr/utils"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strconv"
	"strings"
)

type metaPlugin struct {
	*generator.Generator
}

func (m *metaPlugin) Name() string {
	return "av"
}

func (m *metaPlugin) Init(g *generator.Generator) {
	m.Generator = g
}

func (m *metaPlugin) Generate(file *generator.FileDescriptor) {

	for _, svc := range file.Service {
		m.genServerCode(file.GetPackage(), svc, file)
	}
}

func (m *metaPlugin) GenerateImports(file *generator.FileDescriptor) {
	if len(file.Service) == 0 {
		return
	}
	m.genImportCode(file)
}

func (p *metaPlugin) genImportCode(file *generator.FileDescriptor) {
	p.P("import (")
	p.P(`"context"`)
	p.P(`"regexp"`)
	p.P(")")

}

func (p *metaPlugin) genServerCode(packageName string, svc *descriptor.ServiceDescriptorProto, file *generator.FileDescriptor) {
	messageMap := make(map[string]*descriptor.DescriptorProto, 0)
	for _, v := range file.MessageType {
		messageMap[v.GetName()] = v
	}
	serviceName := strings.ReplaceAll(svc.GetName(), "Service", "")
	p.P("type CreateErrFunc func(string) error")
	p.P(fmt.Sprintf(`type %sServiceServerImpl struct {`, serviceName))
	p.P(fmt.Sprintf(`svc %sServiceServer`, serviceName))
	p.P(fmt.Sprintf(`	cef  CreateErrFunc`))
	p.P(fmt.Sprintf(`}`))

	p.P(fmt.Sprintf(`func New%sServiceServerImpl(svc %sServiceServer, cef CreateErrFunc) *%sServiceServerImpl {`, serviceName, serviceName, serviceName))
	p.P(fmt.Sprintf(`return &%sServiceServerImpl{svc: svc, cef: cef}`, serviceName))
	p.P(fmt.Sprintf(`}`))

	for _, method := range svc.Method {
		inputMessageDesc := messageMap[strings.Split(method.GetInputType(), ".")[2]]

		for _, field := range inputMessageDesc.GetField() {
			opts := field.GetOptions()
			_validator, err := proto.GetExtension(opts, plugin.E_Validator)
			var validator *plugin.Validator
			if err != nil {
				validator = &plugin.Validator{}
			} else {
				validator = _validator.(*plugin.Validator)
			}
			if validator.Regexp != "" {
				p.P(fmt.Sprintf(`var %s%sRegexp=regexp.MustCompile("%s")`, utils.FirstLower(method.GetName()), utils.FirstUpper(field.GetName()), validator.Regexp))
			}

		}

		p.P(fmt.Sprintf(`
	func (m *%sServiceServerImpl) %s(ctx context.Context, req *%s) (*%s, error) {`,
			serviceName,
			method.GetName(),
			strings.Split(method.GetInputType(), ".")[2],
			strings.Split(method.GetOutputType(), ".")[2],
		))

		//加载各个参数
		for _, field := range inputMessageDesc.GetField() {
			opts := field.GetOptions()
			_validator, err := proto.GetExtension(opts, plugin.E_Validator)
			var validator *plugin.Validator
			if err != nil {
				validator = &plugin.Validator{}
			} else {
				validator = _validator.(*plugin.Validator)
			}
			//类型判断
			_type := field.GetType()
			//如果是数组
			if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				if !validator.Omitempty {
					p.P(fmt.Sprintf(`if len(req.%s) == 0 {
						return nil, m.cef("array %s len is 0")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != "" {
					n, _ := strconv.ParseInt(validator.Eq, 10, 64)
					p.P(fmt.Sprintf(`if len(req.%s) != %d {
						return nil, m.cef("array %s len must is %d")
					}`, utils.FirstUpper(field.GetName()), n, field.GetName(), n))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) < %d {
						return nil, m.cef("array %s len must >= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gte), field.GetName(), int64(validator.Gte)))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) <= %d {
						return nil, m.cef("array %s len must > %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gt), field.GetName(), int64(validator.Gt)))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) > %d {
						return nil, m.cef("array %s len must <= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lte), field.GetName(), int64(validator.Lte)))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) >= %d {
						return nil, m.cef("array %s len must < %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lt), field.GetName(), int64(validator.Lt)))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_STRING {
				if !validator.Omitempty {
					p.P(fmt.Sprintf(`if req.%s == "" {
						return nil, m.cef("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != "" {
					p.P(fmt.Sprintf(`if req.%s != "%s" {
						return nil, m.cef("%s must eq %s")
					}`, utils.FirstUpper(field.GetName()), validator.Eq, field.GetName(), validator.Eq))
				}
				if validator.In != "" {
					if validator.In[0:1] != "[" || validator.In[len(validator.In)-1:] != "]" {
						p.Error(errors.New(fmt.Sprintf("参数%s的in验证必须符合[]格式", field.GetName())))
					}
					in := validator.In
					in = strings.ReplaceAll(in, "[", "")
					in = strings.ReplaceAll(in, "]", "")
					arr := strings.Split(in, ",")
					str := ""
					for index, v := range arr {
						str = str + fmt.Sprintf(`req.%s != "%s"`, utils.FirstUpper(field.GetName()), v)
						if index != len(arr)-1 {
							str = str + "&&"
						}
					}
					p.P(fmt.Sprintf(`if %s {
						return nil, m.cef("%s must in %s")
					}`, str, field.GetName(), validator.In))
				}
				if validator.Regexp != "" {
					p.P(fmt.Sprintf(`if !%s%sRegexp.MatchString(req.%s)  {
						return nil, m.cef("%s regexp verification failed")
					}`, utils.FirstLower(method.GetName()), utils.FirstUpper(field.GetName()), utils.FirstUpper(field.GetName()), field.GetName()))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_INT32 ||
				_type == descriptor.FieldDescriptorProto_TYPE_INT64 ||
				_type == descriptor.FieldDescriptorProto_TYPE_UINT64 ||
				_type == descriptor.FieldDescriptorProto_TYPE_UINT32 {
				if !validator.Omitempty {
					p.P(fmt.Sprintf(`if req.%s == 0 {
						return nil, m.cef("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != "" {
					n, _ := strconv.ParseInt(validator.Eq, 10, 64)
					p.P(fmt.Sprintf(`if req.%s != %d {
						return nil, m.cef("%s must is %d")
					}`, utils.FirstUpper(field.GetName()), n, field.GetName(), n))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if req.%s < %d {
						return nil, m.cef("%s must >= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gte), field.GetName(), int64(validator.Gte)))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if req.%s <= %d {
						return nil, m.cef("%s must > %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gt), field.GetName(), int64(validator.Gt)))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if req.%s > %d {
						return nil, m.cef("%s must <= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lte), field.GetName(), int64(validator.Lte)))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if req.%s >= %d {
						return nil, m.cef("%s must < %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lt), field.GetName(), int64(validator.Lt)))
				}
				if validator.In != "" {
					if validator.In[0:1] != "[" || validator.In[len(validator.In)-1:] != "]" {
						p.Error(errors.New(fmt.Sprintf("参数%s的in验证必须符合[]格式", field.GetName())))
					}
					in := validator.In
					in = strings.ReplaceAll(in, "[", "")
					in = strings.ReplaceAll(in, "]", "")
					arr := strings.Split(in, ",")
					str := ""
					for index, v := range arr {
						str = str + fmt.Sprintf(`req.%s != %s`, utils.FirstUpper(field.GetName()), v)
						if index != len(arr)-1 {
							str = str + "&&"
						}
					}
					p.P(fmt.Sprintf(`if %s {
						return nil, m.cef("%s must in %s")
					}`, str, field.GetName(), validator.In))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_DOUBLE ||
				_type == descriptor.FieldDescriptorProto_TYPE_FLOAT {
				if !validator.Omitempty {
					p.P(fmt.Sprintf(`if req.%s == 0 {
						return nil, m.cef("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != "" {
					n, _ := strconv.ParseFloat(validator.Eq, 10)
					p.P(fmt.Sprintf(`if req.%s != %v {
						return nil, m.cef("%s must is %v")
					}`, utils.FirstUpper(field.GetName()), n, field.GetName(), n))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if req.%s < %v {
						return nil, m.cef("%s must >= %v")
					}`, utils.FirstUpper(field.GetName()), validator.Gte, field.GetName(), validator.Gte))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if req.%s <= %v {
						return nil, m.cef("%s must > %v")
					}`, utils.FirstUpper(field.GetName()), validator.Gt, field.GetName(), validator.Gt))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if req.%s > %v {
						return nil, m.cef("%s must <= %v")
					}`, utils.FirstUpper(field.GetName()), validator.Lte, field.GetName(), validator.Lte))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if req.%s >= %v {
						return nil, m.cef("%s must < %v")
					}`, utils.FirstUpper(field.GetName()), validator.Lt, field.GetName(), validator.Lt))
				}
				if validator.In != "" {
					if validator.In[0:1] != "[" || validator.In[len(validator.In)-1:] != "]" {
						p.Error(errors.New(fmt.Sprintf("参数%s的in验证必须符合[]格式", field.GetName())))
					}
					in := validator.In
					in = strings.ReplaceAll(in, "[", "")
					in = strings.ReplaceAll(in, "]", "")
					arr := strings.Split(in, ",")
					str := ""
					for index, v := range arr {
						str = str + fmt.Sprintf(`req.%s != %s`, utils.FirstUpper(field.GetName()), v)
						if index != len(arr)-1 {
							str = str + "&&"
						}
					}
					p.P(fmt.Sprintf(`if %s {
						return nil, m.cef("%s must in %s")
					}`, str, field.GetName(), validator.In))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
				if !validator.Omitempty {
					p.P(fmt.Sprintf(`if req.%s == nil {
						return nil, m.cef("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
			}
		}
		p.P(fmt.Sprintf(`return m.svc.%s(ctx, req)}`, method.GetName()))
	}

}

// 注册
func Init() {
	generator.RegisterPlugin(&metaPlugin{})
}
