### proto-validator
proto-validator是我们在平时开发中为了解决业务痛点而创造的。  为了实现proto dto的定义和参数验证
，原来我们只能在写业务代码时自己再去做参数验证，会造成一定量的重复劳动。而proto-validator可以实现
在定义proto时就可以声明参数验证规则，通过proto-gen-av生成自动验证代码，把参数验证自动处理掉。

### protoc-gen-av  
全称 protoc-gen-auto-validator 简称protoc-gen-av

### 最佳实践
先安装protoc-gen-av插件

go install github.com/Gitforxuyang/proto-validaotr/cmd/protoc-gen-av



