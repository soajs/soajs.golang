package soajsgo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// MCPToolsPath is the route the SOAJS gateway calls to list the service's APIs as MCP tools.
	MCPToolsPath = "/mcp/tools"
	// DefaultSwaggerFilename is the swagger file used when none is set, relative to the folder the service runs from.
	DefaultSwaggerFilename = "/swagger.yml"
)

// MCPTools answers GET /mcp/tools with the service's Swagger 2.0 or OpenAPI 3 document. The SOAJS gateway
// turns the document into MCP tools for AI agents, filtered by each tenant's package and each user's ACL.
//
// NewFromConfig sets it on the registry from the config (swaggerFilename); with New, set Registry.MCP yourself.
// The document is loaded once, on the first call: a JSON document is sent parsed, a YAML one as text, and with
// no document the answer has no tools.
type MCPTools struct {
	ServiceName    string
	ServiceVersion string
	// SwaggerFilename is the document file, relative to the folder the service runs from. Default /swagger.yml.
	SwaggerFilename string
	// SwaggerDocument is the document content (JSON or YAML); it wins over SwaggerFilename.
	SwaggerDocument []byte

	once    sync.Once
	swagger interface{}
}

type mcpToolsResponse struct {
	Result  bool            `json:"result"`
	TS      int64           `json:"ts"`
	Service mcpToolsService `json:"service"`
	Data    mcpToolsData    `json:"data"`
}

type mcpToolsService struct {
	Service string `json:"service"`
	Type    string `json:"type"`
	Route   string `json:"route"`
}

type mcpToolsData struct {
	Name    string         `json:"name"`
	Version string         `json:"version"`
	Swagger interface{}    `json:"swagger,omitempty"`
	Tools   *[]interface{} `json:"tools,omitempty"`
}

// document turns raw content into what is sent: a JSON object parsed, anything else as text, nil when empty or
// when JSON was expected and the content is not a JSON object.
func document(content []byte, isJSON bool) interface{} {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return nil
	}
	if trimmed[0] == '{' && json.Valid(trimmed) {
		return json.RawMessage(trimmed)
	}
	if isJSON {
		return nil
	}
	return string(content)
}

func (m *MCPTools) load() interface{} {
	m.once.Do(func() {
		if len(m.SwaggerDocument) > 0 {
			m.swagger = document(m.SwaggerDocument, false)
			return
		}
		filename := m.SwaggerFilename
		if filename == "" {
			filename = DefaultSwaggerFilename
		}
		wd, err := os.Getwd()
		if err != nil {
			return
		}
		content, err := os.ReadFile(filepath.Join(wd, filepath.Clean("/"+filename)))
		if err != nil {
			return
		}
		m.swagger = document(content, strings.EqualFold(filepath.Ext(filename), ".json"))
	})
	return m.swagger
}

// ServeHTTP answers with the document. A request that came through the gateway proxy (soajsinjectobj header) is
// a client asking for /<service>/mcp/tools, not the gateway: it gets 404, the document is only for the gateway.
func (m *MCPTools) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(headerDataName) != "" {
		http.NotFound(w, r)
		return
	}
	version := m.ServiceVersion
	if version == "" {
		version = "1"
	}
	data := mcpToolsData{Name: m.ServiceName, Version: version}
	if swagger := m.load(); swagger != nil {
		data.Swagger = swagger
	} else {
		data.Tools = &[]interface{}{}
	}
	body, err := json.Marshal(mcpToolsResponse{
		Result: true,
		TS:     time.Now().UnixNano() / int64(time.Millisecond),
		Service: mcpToolsService{
			Service: strings.ToUpper(m.ServiceName),
			Type:    "rest",
			Route:   MCPToolsPath,
		},
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

// isMCPToolsRequest tells whether the request is the gateway asking for the tools.
func isMCPToolsRequest(r *http.Request) bool {
	return r.Method == http.MethodGet && r.URL.Path == MCPToolsPath && r.Header.Get(headerDataName) == ""
}
