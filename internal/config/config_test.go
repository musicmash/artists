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
  name:  'ƒartists'
  login: 'artists'
  pass:  'artists'
  log: false

log:
  level: DEBUG
  file: 'artists.log'

http:
  port: 5566
`
	expected := &AppConfig{
		DB: DBConfig{
			Type:  "mysql",
			Host:  "mariadb",
			Name:  "artists",
			Login: "artists",
			Pass:  "artists",
			Log:   false,
		},
		Log: LogConfig{
			Level:         "DEBUG",
			File:          "artists.log",
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
