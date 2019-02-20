package soajsgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tt := []struct {
		name        string
		conf        Config
		expectedErr error
	}{
		{
			name:        "empty type",
			conf:        Config{},
			expectedErr: errors.New("could not find [Type] in your config, Type is <required>"),
		},
		{
			name: "empty ServiceName",
			conf: Config{
				Type: "type",
			},
			expectedErr: errors.New("could not find [ServiceName] in your config, ServiceName is <required>"),
		},
		{
			name: "empty ServicePort",
			conf: Config{
				Type:        "type",
				ServiceName: "service name",
			},
			expectedErr: errors.New("could not find [ServicePort] in your config, ServicePort is <required>"),
		},
		{
			name: "empty ServiceVersion",
			conf: Config{
				Type:        "type",
				ServiceName: "service name",
				ServicePort: 10,
			},
			expectedErr: errors.New("could not find [ServiceVersion] in your config, ServiceVersion is <required>"),
		},
		{
			name: "bad version",
			conf: Config{
				Type:           "type",
				ServiceName:    "service name",
				ServicePort:    10,
				ServiceVersion: "version",
			},
			expectedErr: errors.New("error with [ServiceVersion] in your config, ServiceVersion syntax is [[0-9]+(.[0-9]+)?]"),
		},
		{
			name: "all ok",
			conf: Config{
				Type:           "type",
				ServiceName:    "service name",
				ServicePort:    10,
				ServiceVersion: "10",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedErr, tc.conf.Validate())
		})
	}
}
