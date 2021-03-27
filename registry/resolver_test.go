package registry_test

import (
	"testing"

	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	"github.com/moderntv/cadre/registry"
	"github.com/moderntv/cadre/registry/static"
)

func Test_resolverBuilder_Build(t *testing.T) {
	type fields struct {
		registry registry.Registry
	}
	type args struct {
		target resolver.Target
		cc     resolver.ClientConn
		opts   resolver.BuildOptions
	}

	sReg, err := static.NewRegistry(map[string][]string{
		"foosvc": {"localhost:5000"},
	})
	if err != nil || sReg == nil {
		t.Errorf("failed to create static registry for test: %v", err)
		return
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "build",
			fields: fields{
				registry: sReg,
			},
			args: args{
				target: resolver.Target{
					Scheme:    "registry",
					Authority: "",
					Endpoint:  "foosvc",
				},
				cc: &clientConn{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := registry.NewResolverBuilder(tt.fields.registry)
			got, err := this.Build(tt.args.target, tt.args.cc, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolverBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("resolverBuilder.Build() = nil, want not nil")
			}
		})
	}
}

// implement clientconn for tests
type clientConn struct{}

// UpdateState updates the state of the ClientConn appropriately.
func (this *clientConn) UpdateState(s resolver.State) {}

// ReportError notifies the ClientConn that the Resolver encountered an
// error.  The ClientConn will notify the load balancer and begin calling
// ResolveNow on the Resolver with exponential backoff.
func (this *clientConn) ReportError(e error) {}

// NewAddress is called by resolver to notify ClientConn a new list
// of resolved addresses.
// The address list should be the complete list of resolved addresses.
//
// Deprecated: Use UpdateState instead.
func (this *clientConn) NewAddress(addresses []resolver.Address) {}

// NewServiceConfig is called by resolver to notify ClientConn a new
// service config. The service config should be provided as a json string.
//
// Deprecated: Use UpdateState instead.
func (this *clientConn) NewServiceConfig(serviceConfig string) {}

// ParseServiceConfig parses the provided service config and returns an
// object that provides the parsed config.
func (this *clientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}
