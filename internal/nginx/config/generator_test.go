package config_test

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/types"

	"github.com/nginxinc/nginx-kubernetes-gateway/internal/nginx/config"
	"github.com/nginxinc/nginx-kubernetes-gateway/internal/state/dataplane"
)

// Note: this test only verifies that Generate() returns a byte array with upstream, server, and split_client blocks.
// It does not test the correctness of those blocks. That functionality is covered by other tests in this package.
func TestGenerate(t *testing.T) {
	bg := dataplane.BackendGroup{
		Source:  types.NamespacedName{Namespace: "test", Name: "hr"},
		RuleIdx: 0,
		Backends: []dataplane.Backend{
			{UpstreamName: "test", Valid: true, Weight: 1},
			{UpstreamName: "test2", Valid: true, Weight: 1},
		},
	}

	conf := dataplane.Configuration{
		HTTPServers: []dataplane.VirtualServer{
			{
				IsDefault: true,
			},
			{
				Hostname: "example.com",
			},
		},
		SSLServers: []dataplane.VirtualServer{
			{
				IsDefault: true,
			},
			{
				Hostname: "example.com",
				SSL: &dataplane.SSL{
					CertificatePath: "/etc/nginx/secrets/default",
				},
			},
		},
		Upstreams: []dataplane.Upstream{
			{
				Name:      "up",
				Endpoints: nil,
			},
		},
		BackendGroups: []dataplane.BackendGroup{bg},
	}
	generator := config.NewGeneratorImpl()
	cfg := string(generator.Generate(conf))

	if !strings.Contains(cfg, "listen 80") {
		t.Errorf("Generate() did not generate a config with a default HTTP server; config: %s", cfg)
	}

	if !strings.Contains(cfg, "listen 443") {
		t.Errorf("Generate() did not generate a config with an SSL server; config: %s", cfg)
	}

	if !strings.Contains(cfg, "upstream") {
		t.Errorf("Generate() did not generate a config with an upstream block; config: %s", cfg)
	}

	if !strings.Contains(cfg, "split_clients") {
		t.Errorf("Generate() did not generate a config with an split_clients block; config: %s", cfg)
	}
}
