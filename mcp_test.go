package soajsgo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const swaggerYAML = `swagger: "2.0"
info:
  title: demo
  version: "1.0.0"
paths:
  /hello:
    get:
      summary: Hello
`

type mcpAnswer struct {
	Result  bool `json:"result"`
	Service struct {
		Service string `json:"service"`
		Type    string `json:"type"`
		Route   string `json:"route"`
	} `json:"service"`
	Data struct {
		Name    string          `json:"name"`
		Version string          `json:"version"`
		Swagger json.RawMessage `json:"swagger"`
		Tools   []interface{}   `json:"tools"`
	} `json:"data"`
	raw map[string]interface{}
}

// inDir runs the test from dir, the way a service finds its swagger file.
func inDir(t *testing.T, dir string) {
	t.Helper()
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
}

func getTools(t *testing.T, h http.Handler, header string) (*httptest.ResponseRecorder, mcpAnswer) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:4000/mcp/tools", nil)
	if header != "" {
		req.Header.Set(headerDataName, header)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	var answer mcpAnswer
	if rec.Code == http.StatusOK {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &answer))
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &answer.raw))
	}
	return rec, answer
}

func TestMCPTools_YAMLFileAsText(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "swagger.yml"), []byte(swaggerYAML), 0o600))
	inDir(t, dir)

	rec, answer := getTools(t, &MCPTools{ServiceName: "demo", ServiceVersion: "2"}, "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.True(t, answer.Result)
	assert.Equal(t, "DEMO", answer.Service.Service)
	assert.Equal(t, "rest", answer.Service.Type)
	assert.Equal(t, MCPToolsPath, answer.Service.Route)
	assert.Equal(t, "demo", answer.Data.Name)
	assert.Equal(t, "2", answer.Data.Version)
	var text string
	require.NoError(t, json.Unmarshal(answer.Data.Swagger, &text))
	assert.Equal(t, swaggerYAML, text)
	assert.NotContains(t, answer.raw["data"], "tools")
}

func TestMCPTools_JSONFileParsed(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "api.json"), []byte(`{"openapi": "3.0.0", "paths": {}}`), 0o600))
	inDir(t, dir)

	_, answer := getTools(t, &MCPTools{ServiceName: "demo", SwaggerFilename: "/api.json"}, "")
	assert.Equal(t, "1", answer.Data.Version)
	assert.JSONEq(t, `{"openapi": "3.0.0", "paths": {}}`, string(answer.Data.Swagger))
}

func TestMCPTools_NoDocument(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "bad.json"), []byte(`{not json`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "list.json"), []byte(`[1, 2]`), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "empty.yml"), []byte("  \n"), 0o600))
	inDir(t, dir)

	for _, filename := range []string{"", "/missing.yml", "bad.json", "list.json", "empty.yml"} {
		t.Run("file "+filename, func(t *testing.T) {
			_, answer := getTools(t, &MCPTools{ServiceName: "demo", SwaggerFilename: filename}, "")
			assert.True(t, answer.Result)
			assert.Equal(t, []interface{}{}, answer.Data.Tools)
			assert.NotContains(t, answer.raw["data"], "swagger")
		})
	}
}

func TestMCPTools_DocumentWinsOverFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "swagger.yml"), []byte(swaggerYAML), 0o600))
	inDir(t, dir)

	_, answer := getTools(t, &MCPTools{ServiceName: "demo", SwaggerDocument: []byte(` {"openapi": "3.1.0"} `)}, "")
	assert.JSONEq(t, `{"openapi": "3.1.0"}`, string(answer.Data.Swagger))

	_, answer = getTools(t, &MCPTools{ServiceName: "demo", SwaggerDocument: []byte("openapi: 3.1.0\n")}, "")
	var text string
	require.NoError(t, json.Unmarshal(answer.Data.Swagger, &text))
	assert.Equal(t, "openapi: 3.1.0\n", text)
}

func TestMCPTools_LoadedOnce(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "swagger.yml")
	require.NoError(t, os.WriteFile(file, []byte(swaggerYAML), 0o600))
	inDir(t, dir)

	tools := &MCPTools{ServiceName: "demo"}
	getTools(t, tools, "")
	require.NoError(t, os.Remove(file))
	_, answer := getTools(t, tools, "")
	assert.NotEmpty(t, answer.Data.Swagger)
}

func TestMCPTools_PathCannotLeaveTheServiceFolder(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "service")
	require.NoError(t, os.Mkdir(dir, 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(parent, "secret.yml"), []byte("secret"), 0o600))
	inDir(t, dir)

	_, answer := getTools(t, &MCPTools{ServiceName: "demo", SwaggerFilename: "../secret.yml"}, "")
	assert.Equal(t, []interface{}{}, answer.Data.Tools)
}

func TestMCPTools_ThroughGatewayIsNotFound(t *testing.T) {
	rec, _ := getTools(t, &MCPTools{ServiceName: "demo"}, "{}")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRegistry_MiddlewareMCPTools(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "swagger.yml"), []byte(swaggerYAML), 0o600))
	inDir(t, dir)

	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	reg := &Registry{MCP: &MCPTools{ServiceName: "demo"}}
	h := reg.Middleware(app)

	rec, answer := getTools(t, h, "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "demo", answer.Data.Name)

	rec, _ = getTools(t, h, "{}")
	assert.Equal(t, http.StatusTeapot, rec.Code, "a request through the gateway goes to the app")

	for _, req := range []*http.Request{
		httptest.NewRequest(http.MethodPost, "http://localhost:4000/mcp/tools", nil),
		httptest.NewRequest(http.MethodGet, "http://localhost:4000/mcp/tools/x", nil),
	} {
		rec = httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code, req.Method+" "+req.URL.Path)
	}

	rec, _ = getTools(t, (&Registry{}).Middleware(app), "")
	assert.Equal(t, http.StatusTeapot, rec.Code, "without MCP the app handles the route")
}

func TestManualDeploy_NotSet(t *testing.T) {
	last, had := os.LookupEnv(EnvDeployManual)
	t.Cleanup(func() {
		if had {
			_ = os.Setenv(EnvDeployManual, last)
		} else {
			_ = os.Unsetenv(EnvDeployManual)
		}
	})
	require.NoError(t, os.Unsetenv(EnvDeployManual))
	assert.NoError(t, manualDeploy(Config{ServiceName: "test"}, registryPath("localhost")))

	require.NoError(t, os.Setenv(EnvDeployManual, "false"))
	assert.NoError(t, manualDeploy(Config{ServiceName: "test"}, registryPath("localhost")))
}

func TestNewFromConfig_SetsMCP(t *testing.T) {
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"result": true, "data": {"name": "name", "environment": "env"}}`))
	}))
	defer gateway.Close()

	for _, env := range []string{EnvRegistryAPIAddress, EnvSoajsEnv, EnvDeployManual} {
		last, had := os.LookupEnv(env)
		name := env
		t.Cleanup(func() {
			if had {
				_ = os.Setenv(name, last)
			} else {
				_ = os.Unsetenv(name)
			}
		})
	}
	require.NoError(t, os.Setenv(EnvRegistryAPIAddress, strings.TrimPrefix(gateway.URL, "http://")))
	require.NoError(t, os.Setenv(EnvSoajsEnv, "dev"))
	require.NoError(t, os.Unsetenv(EnvDeployManual))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reg, err := NewFromConfig(ctx, Config{
		Type:            "service",
		ServiceName:     "name",
		ServiceVersion:  "2",
		ServicePort:     4000,
		ServiceGroup:    "group-a",
		SwaggerFilename: "/api/swagger.json",
		Maintenance: maintenance{
			Port:      maintenancePort{Type: "inherit"},
			Readiness: "/heartbeat",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, reg.MCP)
	assert.Equal(t, "name", reg.MCP.ServiceName)
	assert.Equal(t, "2", reg.MCP.ServiceVersion)
	assert.Equal(t, "/api/swagger.json", reg.MCP.SwaggerFilename)
}
