package client

import (
	"context"
	"testing"
)

func TestClientV2_UsesConfiguredHost(t *testing.T) {
	tests := []struct {
		name string
		host string
		want string
	}{
		{name: "custom host", host: "api.qase.example.com", want: "https://api.qase.example.com/v2"},
		{name: "empty host falls back to default", host: "", want: "https://api.qase.io/v2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClientV2("token", tt.host, NewClientV1("token", tt.host))

			_, apiClient := c.getApiV2Client(context.Background())

			servers := apiClient.GetConfig().Servers
			if len(servers) == 0 {
				t.Fatal("no servers configured")
			}
			if servers[0].URL != tt.want {
				t.Errorf("server URL = %q, want %q", servers[0].URL, tt.want)
			}
		})
	}
}
