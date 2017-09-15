package grpchelper

import (
	"github.com/fresh8/consul"
	"github.com/olivere/grpc/lb/consul"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	defaultOpts = []grpc.DialOption{grpc.WithInsecure()}
)

func DialConsul(service, tag string) (conn *grpc.ClientConn, err error) {
	return DialConsulWithOpts(service, tag)
}

// DialConsulWithOpts dials consul with provided options, including default options
func DialConsulWithOpts(service, tag string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	lbOpts, err := addLBOptionToDefaultOpts(service, tag)
	if err != nil {
		return nil, errors.Wrap(err, "DialConsulWithOpts")
	}
	finOptions := append(lbOpts, opts...)
	conn, err = grpc.Dial("", finOptions...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}
	return
}

// addLBOptionToDefaultOpts sets up the load balancing options within deaultOpts
func addLBOptionToDefaultOpts(service, tag string) ([]grpc.DialOption, error) {
	if consul.Client == nil {
		return nil, ErrNoConsul
	}
	r, err := lb.NewConsulResolver(consul.Client, service, tag)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Consul resolver")
	}
	return append(defaultOpts, grpc.WithBalancer(grpc.RoundRobin(r))), nil
}

var (
	ErrNoConsul = errors.New("initialise Consul first")
)
