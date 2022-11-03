package gen

import (
	"fmt"
	"github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av/generator"
	"github.com/Gitforxuyang/proto-validaotr/proto/plugin"
	"github.com/Gitforxuyang/proto-validaotr/utils"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
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
	//if len(file.GetExtension()) > 0 {
	//	m.fields = file.GetExtension()
	//}
	for _, svc := range file.Service {
		m.genClientCode(file.GetPackage(), svc)
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
	p.P(`"google.golang.org/grpc"`)
	p.P(`merr "meta/frame/err"`)
	p.P(`mgrpc "meta/frame/grpc"`)
	p.P(")")
	p.P(`var (
		_ merr.ErrorCode
	)`)
}

func (p *metaPlugin) genServerCode(packageName string, svc *descriptor.ServiceDescriptorProto, file *generator.FileDescriptor) {
	messageMap := make(map[string]*descriptor.DescriptorProto, 0)
	for _, v := range file.MessageType {
		messageMap[v.GetName()] = v
	}
	serviceName := strings.ReplaceAll(svc.GetName(), "Service", "")
	p.P(fmt.Sprintf(`type %sServiceServerImpl struct {`, serviceName))
	p.P(fmt.Sprintf(`svc %sServiceServer`, serviceName))
	p.P(fmt.Sprintf(`}`))

	p.P(fmt.Sprintf(`func New%sServiceServerImpl(svc %sServiceServer) *%sServiceServerImpl {`, serviceName, serviceName, serviceName))
	p.P(fmt.Sprintf(`return &%sServiceServerImpl{svc: svc}`, serviceName))
	p.P(fmt.Sprintf(`}`))

	for _, method := range svc.Method {
		inputMessageDesc := messageMap[strings.Split(method.GetInputType(), ".")[2]]

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
			//doc, err := proto.GetExtension(opts, plugin.E_Doc)
			_validator, err := proto.GetExtension(opts, plugin.E_Validator)
			var validator *plugin.Validator
			if err != nil {
				validator = &plugin.Validator{}
			} else {
				validator = _validator.(*plugin.Validator)
			}
			//非空盘那段
			_type := field.GetType()
			//如果是数组
			if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				if validator.NotEmpty {
					p.P(fmt.Sprintf(`if len(req.%s) == 0 {
						return nil, merr.ParamsError.WithMsg("array %s len is 0")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) != %d {
						return nil, merr.ParamsError.WithMsg("array %s len must is %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Eq), field.GetName(), int64(validator.Eq)))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) < %d {
						return nil, merr.ParamsError.WithMsg("array %s len must >= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gte), field.GetName(), int64(validator.Gte)))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) <= %d {
						return nil, merr.ParamsError.WithMsg("array %s len must > %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gt), field.GetName(), int64(validator.Gt)))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) > %d {
						return nil, merr.ParamsError.WithMsg("array %s len must <= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lte), field.GetName(), int64(validator.Lte)))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if len(req.%s) >= %d {
						return nil, merr.ParamsError.WithMsg("array %s len must < %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lt), field.GetName(), int64(validator.Lt)))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_STRING {
				if validator.NotEmpty {
					p.P(fmt.Sprintf(`if req.%s == "" {
						return nil, merr.ParamsError.WithMsg("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.StringEq != "" {
					p.P(fmt.Sprintf(`if req.%s != %s {
						return nil, merr.ParamsError.WithMsg("%s must eq %s")
					}`, utils.FirstUpper(field.GetName()), validator.StringEq, field.GetName(), validator.StringEq))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_INT32 ||
				_type == descriptor.FieldDescriptorProto_TYPE_INT64 ||
				_type == descriptor.FieldDescriptorProto_TYPE_UINT64 ||
				_type == descriptor.FieldDescriptorProto_TYPE_UINT32 {
				if validator.NotEmpty {
					p.P(fmt.Sprintf(`if req.%s == 0 {
						return nil, merr.ParamsError.WithMsg("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != 0 {
					p.P(fmt.Sprintf(`if req.%s != %d {
						return nil, merr.ParamsError.WithMsg("%s must is %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Eq), field.GetName(), int64(validator.Eq)))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if req.%s < %d {
						return nil, merr.ParamsError.WithMsg("%s must >= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gte), field.GetName(), int64(validator.Gte)))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if req.%s <= %d {
						return nil, merr.ParamsError.WithMsg("%s must > %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Gt), field.GetName(), int64(validator.Gt)))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if req.%s > %d {
						return nil, merr.ParamsError.WithMsg("%s must <= %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lte), field.GetName(), int64(validator.Lte)))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if req.%s >= %d {
						return nil, merr.ParamsError.WithMsg("%s must < %d")
					}`, utils.FirstUpper(field.GetName()), int64(validator.Lt), field.GetName(), int64(validator.Lt)))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_DOUBLE ||
				_type == descriptor.FieldDescriptorProto_TYPE_FLOAT {
				if validator.NotEmpty {
					p.P(fmt.Sprintf(`if req.%s == 0 {
						return nil, merr.ParamsError.WithMsg("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
				if validator.Eq != 0 {
					p.P(fmt.Sprintf(`if req.%s != %v {
						return nil, merr.ParamsError.WithMsg("%s must is %v")
					}`, utils.FirstUpper(field.GetName()), validator.Eq, field.GetName(), validator.Eq))
				}
				if validator.Gte != 0 {
					p.P(fmt.Sprintf(`if req.%s < %v {
						return nil, merr.ParamsError.WithMsg("%s must >= %v")
					}`, utils.FirstUpper(field.GetName()), validator.Gte, field.GetName(), validator.Gte))
				}
				if validator.Gt != 0 {
					p.P(fmt.Sprintf(`if req.%s <= %v {
						return nil, merr.ParamsError.WithMsg("%s must > %v")
					}`, utils.FirstUpper(field.GetName()), validator.Gt, field.GetName(), validator.Gt))
				}
				if validator.Lte != 0 {
					p.P(fmt.Sprintf(`if req.%s > %v {
						return nil, merr.ParamsError.WithMsg("%s must <= %v")
					}`, utils.FirstUpper(field.GetName()), validator.Lte, field.GetName(), validator.Lte))
				}
				if validator.Lt != 0 {
					p.P(fmt.Sprintf(`if req.%s >= %v {
						return nil, merr.ParamsError.WithMsg("%s must < %v")
					}`, utils.FirstUpper(field.GetName()), validator.Lt, field.GetName(), validator.Lt))
				}
			} else if _type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
				if validator.NotEmpty {
					p.P(fmt.Sprintf(`if req.%s == nil {
						return nil, merr.ParamsError.WithMsg("%s can not empty")
					}`, utils.FirstUpper(field.GetName()), field.GetName()))
				}
			}
		}
		p.P(fmt.Sprintf(`return m.svc.%s(ctx, req)}`, method.GetName()))
	}

}
func (p *metaPlugin) genClientCode(packageName string, svc *descriptor.ServiceDescriptorProto) {
	serviceName := strings.ReplaceAll(svc.GetName(), "Service", "")
	p.P(fmt.Sprintf(`type Grpc%sServiceClient interface {`, serviceName))
	p.P("Close() error")
	for _, method := range svc.Method {
		p.P(fmt.Sprintf(`%s(ctx context.Context, req *%s) (*%s, error)`, method.GetName(),
			strings.Split(method.GetInputType(), ".")[2],
			strings.Split(method.GetOutputType(), ".")[2],
		))
	}
	p.P("}")

	p.P(fmt.Sprintf(`
	type grpc%sServiceClient struct {
		client %sServiceClient
		conn   *grpc.ClientConn
	}`, serviceName, serviceName))

	for _, method := range svc.Method {
		p.P(fmt.Sprintf(`
	func (g *grpc%sServiceClient) %s(ctx context.Context, req *%s) (*%s, error) {
		return g.client.%s(ctx, req)
	}`,
			serviceName,
			method.GetName(),
			strings.Split(method.GetInputType(), ".")[2],
			strings.Split(method.GetOutputType(), ".")[2],
			method.GetName(),
		))
	}
	p.P(fmt.Sprintf(`func (m *grpc%sServiceClient) Close() error {
	return m.conn.Close()
}`, serviceName))
	upper := fmt.Sprintf(`
	func GetGrpc%sServiceClient() (Grpc%sServiceClient, error){
		conn, err := mgrpc.MakeConn("%s")
		if err != nil {
			return nil, err
		}
		client := New%sServiceClient(conn)
		return &grpc%sServiceClient{client: client, conn: conn},nil
	}
	`, serviceName, serviceName, serviceName, serviceName, serviceName)
	p.P(upper)
	//p.P("grpc.WithDefaultServiceConfig(fmt.Sprintf(`" + `{"loadBalancingPolicy":"%s"}` + "`, loadBalancingPolicy)),")
	//lower := fmt.Sprintf(`
	//	grpc.WithChainUnaryInterceptor(mgrpc.ClientMiddleware(opt)...),
	//}
	//for _, v := range mgrpc.DefaultDialOptions {
	//	options = append(options, v)
	//}
	//conn, err := grpc.Dial(target, options...)
	//	if err != nil {
	//		return nil, err
	//	}
	//	client := New%sServiceClient(conn)
	//	return &grpc%sServiceClient{client: client, conn: conn},nil
	//}
	//`, serviceName, serviceName)
	//p.P(lower)
}

// 注册
func Init() {
	generator.RegisterPlugin(&metaPlugin{})
}
