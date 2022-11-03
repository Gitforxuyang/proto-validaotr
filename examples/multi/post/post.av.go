// Code generated by protoc-gen-av. DO NOT EDIT.
// source: post.proto

package post

import (
	"context"
	"regexp"
)

type CreateErrFunc func(string) error
type PostServiceServerImpl struct {
	svc PostServiceServer
	cef CreateErrFunc
}

func NewPostServiceServerImpl(svc PostServiceServer, cef CreateErrFunc) *PostServiceServerImpl {
	return &PostServiceServerImpl{svc: svc, cef: cef}
}

var pingPongNameRegexp = regexp.MustCompile("123")

func (m *PostServiceServerImpl) PingPong(ctx context.Context, req *Ping) (*Pong, error) {
	if req.Name == "" {
		return nil, m.cef("name can not empty")
	}
	if req.Name != "1" && req.Name != "2" && req.Name != "3" {
		return nil, m.cef("name must in [1,2,3]")
	}
	if !pingPongNameRegexp.MatchString(req.Name) {
		return nil, m.cef("name regexp verification failed")
	}
	return m.svc.PingPong(ctx, req)
}