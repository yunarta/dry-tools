package dry_config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDryConfig_Resolve(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	resolve, err := config.Resolve("nexus", "one")
	assert.NotNil(t, resolve)

	assert.Equal(t, resolve["username"], "one")
	assert.Equal(t, resolve["password"], "secret")
}
