package jwtcentersdk

import (
	resolver "google.golang.org/grpc/resolver"
)

const (
	localScheme      = "component"
	localServiceName = "resolver.jwtcenter.grpc.io"
)

// exampleResolver is a
// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type localResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *localResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*localResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*localResolver) Close()                                  {}

type localResolverBuilder struct {
	BackendAddr []string
}

func NewLocalResolverBuilder(backendAddr []string) *localResolverBuilder {
	rb := new(localResolverBuilder)
	rb.BackendAddr = backendAddr
	return rb
}

func (rb *localResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &localResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			localServiceName: rb.BackendAddr,
		},
	}
	r.start()
	return r, nil
}
func (*localResolverBuilder) Scheme() string { return localScheme }

func (rb *localResolverBuilder) RegistToResolver() {
	resolver.Register(rb)
}
