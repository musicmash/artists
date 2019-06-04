package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	// arrange
	data := `
---
db:
  type:  'mysql'
  host:  'mariadb'
  name:  'musicmash-artists'
  login: 'musicmash-artists'
  pass:  'musicmash-artists'
  log: false

log:
  level: DEBUG
  file: 'musicmash-artists.log'

http:
  port: 5566
`
	expected := &AppConfig{
		DB: DBConfig{
			Type:  "mysql",
			Host:  "mariadb",
			Name:  "musicmash-artists",
			Login: "musicmash-artists",
			Pass:  "musicmash-artists",
			Log:   false,
		},
		Log: LogConfig{
			Level:         "DEBUG",
			File:          "musicmash-artists.log",
		},
		HTTP: HTTPConfig{
			Port: 5566,
		},
	}

	// action
	err := Load([]byte(data))

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expected, Config)
}
