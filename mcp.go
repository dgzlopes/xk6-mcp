package compare

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/sobek"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/mcp", new(MCP))
}

// MCP is the root module struct
type MCP struct{}

// NewModuleInstance initializes a new module instance
func (*MCP) NewModuleInstance(vu modules.VU) modules.Instance {
	logger := vu.InitEnv().Logger.WithField("component", "xk6-mcp")
	return &MCPInstance{
		vu:     vu,
		logger: logger,
	}
}

// MCPInstance represents an instance of the MCP module
type MCPInstance struct {
	vu     modules.VU
	logger logrus.FieldLogger
}

// Exports defines the JavaScript-accessible functions
func (m *MCPInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"StdioClient": m.newStdioClient,
			"SSEClient":   m.newSSEClient,
		},
	}
}

// ClientConfig represents the configuration for the MCP client
type ClientConfig struct {
	// Stdio
	Path string
	Args []string
	Env  map[string]string

	// SSE
	BaseURL string
	Headers map[string]string
	Timeout time.Duration
}

// Client wraps an MCPClient
type Client struct {
	mcp_client client.MCPClient
}

func (m *MCPInstance) newStdioClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	err := rt.ExportTo(c.Argument(0), &cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("invalid config: %w", err))
	}

	stdioClient, err := createStdioClient(cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("Stdio client error: %w", err))
	}

	return m.initializeClient(rt, stdioClient)
}

func (m *MCPInstance) newSSEClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	err := rt.ExportTo(c.Argument(0), &cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("invalid config: %w", err))
	}

	sseClient, err := createSSEClient(cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("SSE client error: %w", err))
	}

	return m.initializeClient(rt, sseClient)
}

func (m *MCPInstance) initializeClient(rt *sobek.Runtime, cl client.MCPClient) *sobek.Object {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "k6",
		Version: "1.0.0",
	}

	if _, err := cl.Initialize(ctx, initReq); err != nil {
		common.Throw(rt, fmt.Errorf("initialize error: %w", err))
	}

	return rt.ToValue(&Client{mcp_client: cl}).ToObject(rt)
}

func createStdioClient(cfg ClientConfig) (*client.StdioMCPClient, error) {
	env := []string{}
	for k, v := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return client.NewStdioMCPClient(cfg.Path, env, cfg.Args...)
}

func createSSEClient(cfg ClientConfig) (*client.SSEMCPClient, error) {
	opts := []client.ClientOption{}
	if cfg.Headers != nil {
		opts = append(opts, client.WithHeaders(cfg.Headers))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, client.WithSSEReadTimeout(cfg.Timeout))
	}
	return client.NewSSEMCPClient(cfg.BaseURL, opts...)
}

func (c *Client) Ping() bool {
	err := c.mcp_client.Ping(context.Background())
	return err == nil
}

func (c *Client) ListTools(r mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	return c.mcp_client.ListTools(context.Background(), r)
}

func (c *Client) CallTool(r mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return c.mcp_client.CallTool(context.Background(), r)
}

func (c *Client) ListResources(r mcp.ListResourcesRequest) (*mcp.ListResourcesResult, error) {
	return c.mcp_client.ListResources(context.Background(), r)
}

func (c *Client) ReadResource(r mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return c.mcp_client.ReadResource(context.Background(), r)
}

func (c *Client) ListPrompts(r mcp.ListPromptsRequest) (*mcp.ListPromptsResult, error) {
	return c.mcp_client.ListPrompts(context.Background(), r)
}

func (c *Client) GetPrompt(r mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return c.mcp_client.GetPrompt(context.Background(), r)
}
