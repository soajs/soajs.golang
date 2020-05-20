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
			expectedErr: errors.New("could not find [Type] in your config, type is <required>"),
		},
		{
			name: "empty ServiceName",
			conf: Config{
				Type: "type",
			},
			expectedErr: errors.New("could not find [ServiceName] in your config, name is <required>"),
		},
		{
			name: "bad ServiceName",
			conf: Config{
				Type:        "type",
				ServiceName: "service name",
			},
			expectedErr: errors.New("error with [ServiceName] in your config, name syntax is [^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$]"),
		},
		{
			name: "empty ServicePort",
			conf: Config{
				Type:        "type",
				ServiceName: "servicename",
			},
			expectedErr: errors.New("could not find [ServicePort] in your config, port is <required>"),
		},
		{
			name: "empty ServiceVersion",
			conf: Config{
				Type:        "type",
				ServiceName: "servicename",
				ServicePort: 400,
			},
			expectedErr: errors.New("could not find [ServiceVersion] in your config, version is <required>"),
		},
		{
			name: "bad version",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "version",
			},
			expectedErr: errors.New("error with [ServiceVersion] in your config, version syntax is [[0-9]+(.[0-9]+)?]"),
		},
		{
			name: "empty maintenance readiness",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "1",
			},
			expectedErr: errors.New("could not find [Readiness] in your config, maintenance.readiness is <required>"),
		},
		{
			name: "empty maintenance port type",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "1",
				Maintenance: maintenance{
					Readiness: "/heartbeat",
				},
			},
			expectedErr: errors.New("could not find [Maintenance Port Type] in your config, maintenance.port.type is <required>"),
		},
		{
			name: "empty ServiceGroup",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "1",
				Maintenance: maintenance{
					Port: maintenancePort{
						Type: "inherit",
					},
					Readiness: "/heartbeat",
				},
			},
			expectedErr: errors.New("could not find [ServiceGroup] in your config, group is <required>"),
		},
		{
			name: "bad ServiceGroup",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "1",
				Maintenance: maintenance{
					Port: maintenancePort{
						Type: "inherit",
					},
					Readiness: "/heartbeat",
				},
				ServiceGroup: "group A",
			},
			expectedErr: errors.New("error with [ServiceGroup] in your config, group syntax is [^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$]"),
		},
		{
			name: "all ok",
			conf: Config{
				Type:           "type",
				ServiceName:    "servicename",
				ServicePort:    4000,
				ServiceVersion: "1",
				Maintenance: maintenance{
					Port: maintenancePort{
						Type: "inherit",
					},
					Readiness: "/heartbeat",
				},
				ServiceGroup: "group-a",
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
