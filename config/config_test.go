package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		cfg, err := LoadConfig("./", "config-local")
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})

	t.Run("failed load config", func(t *testing.T) {
		cfg, err := LoadConfig("./wrong-path", "config-wrong")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	// 	t.Run("failed to unmarshal", func(t *testing.T) {
	// 		// buat temporary file config
	// 		invalidContent := []byte(`server:
	// 	AppVersion: 12
	// 		Port: {}
	// 	Mode: development
	// `)

	// 		tmpFile := "config-test.yml"
	// 		err := os.WriteFile(tmpFile, invalidContent, 0o644)
	// 		assert.NoError(t, err)
	// 		// defer os.Remove(tmpFile)

	//		cfg, err := LoadConfig(".", "config-test")
	//		assert.Error(t, err)
	//		assert.Nil(t, cfg)
	//		log.Info().Any("cfg", cfg).Msg("cfg")
	//	})
}
