package soajsgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type (
	key int
)

const (
	// headerDataName is the SOAJS Gateway injected object attached to the header of each request
	// between the gateway and tech service.
	headerDataName = "soajsinjectobj"
	// SoajsKey use this key to init soajs data from context.
	SoajsKey = key(1)
)

// Middleware is http middleware that gets triggered per request.
func (reg *Registry) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d, err := headerData(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		out := ContextData{
			Tenant:         d.Tenant,
			Urac:           d.Urac,
			ServicesConfig: d.Key.Config,
			Device:         d.Device,
			Geo:            d.Geo,
			Awareness:      d.Awareness,
			Reg:            reg,
		}
		out.Tenant.Key.IKey = d.Key.IKey
		out.Tenant.Key.EKey = d.Key.EKey

		out.Tenant.Application = d.Application
		out.Tenant.Application.PackageACL = d.Package.ACL
		out.Tenant.Application.PackageACLAllEnv = d.Package.ACLAllEnv

		soajs := context.WithValue(r.Context(), SoajsKey, out)
		next.ServeHTTP(w, r.WithContext(soajs))
	})
}

func headerData(r *http.Request) (*headerInfo, error) {
	info := strings.NewReader(r.Header.Get(headerDataName))
	var d *headerInfo
	if err := json.NewDecoder(info).Decode(&d); err != nil {
		return nil, fmt.Errorf("unable to parse SOAJS header: %v", err)
	}
	if d == nil {
		return nil, fmt.Errorf("SOAJS header is empty or null")
	}
	return d, nil
}

// Path returns compiled service path.
func (a Host) Path(args ...string) string {
	var serviceName, version string
	switch len(args) {
	case 1: // controller
		serviceName = args[0]
	case 2, 3: // controller, 1, dash [dash is ignored]
		serviceName = args[0]
		version = args[1]
	}
	host := fmt.Sprintf("%s:%d/", a.Host, a.Port)
	if strings.EqualFold(serviceName, "controller") {
		host = fmt.Sprintf("%s%s/", host, serviceName)
		if _, err := strconv.Atoi(version); err == nil {
			host = fmt.Sprintf("%sv%s/", host, version)
		}
	}
	return host
}

// Connect handles the Mesh inter connect between micro service
func (c ContextData) Connect(args ...string) Connect {
	var connectResponse Connect
	var serviceName, version string
	switch len(args) {
	case 1: // controller
		serviceName = args[0]
	case 2, 3: // controller, 1, dash [dash is ignored]
		serviceName = args[0]
		version = args[1]
	}

	// Try to find service in InterConnect mesh for direct service-to-service communication
	foundInMesh := false
	if c.Awareness.InterConnect != nil {
		for i := 0; i < len(c.Awareness.InterConnect); i++ {
			service := c.Awareness.InterConnect[i]
			if serviceName == service.Name {
				// Match by version: either request latest version or specific version
				isLatestVersionMatch := version == "" && service.Version == service.Latest
				isSpecificVersionMatch := version != "" && version == service.Version

				if isLatestVersionMatch || isSpecificVersionMatch {
					connectResponse.Host = fmt.Sprintf("%s:%d", service.Host, service.Port)
					foundInMesh = true
					break
				}
			}
		}
	}

	// Service found in mesh: use direct connection with full SOAJS context
	if foundInMesh {
		connectResponse.Headers.SoajsInjectobj = headerInfo{
			Tenant: c.Tenant,
			Key: Key{
				IKey:   c.Tenant.Key.IKey,
				EKey:   c.Tenant.Key.EKey,
				Config: c.ServicesConfig,
			},
			Application: c.Tenant.Application,
			Package: Package{
				ACL:       c.Tenant.Application.PackageACL,
				ACLAllEnv: c.Tenant.Application.PackageACLAllEnv,
			},
			Device:    c.Device,
			Geo:       c.Geo,
			Urac:      c.Urac,
			Awareness: c.Awareness,
		}
	} else {
		// Service not found in mesh: fallback to gateway routing via Awareness.Path
		connectResponse.Host = c.Awareness.Path(serviceName, version)
		connectResponse.Headers.Key = c.Tenant.Key.EKey
		// connectResponse.Headers.AccessToken =
	}

	return connectResponse
}
