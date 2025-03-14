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

// Register the module with k6
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
			"StdioClient": m.stdioClient,
		},
	}
}

// Client represents the MCP client
type Client struct {
	mcp_client *client.StdioMCPClient
}

// ClientConfig represents the configuration for the MCP client
type ClientConfig struct {
	Path string
	Args []string
	Env  map[string]string
}

// client constructor for JavaScript runtime
// Usage in JS: `const client = new mcp.StdioClient();`
func (m *MCPInstance) stdioClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	err := rt.ExportTo(c.Argument(0), &cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("unable to create client: constructor expects first argument to be ClientConfig: %w", err))
	}

	if cfg.Path == "" {
		common.Throw(rt, fmt.Errorf("unable to create client: path is required"))
	}

	if cfg.Args == nil {
		cfg.Args = []string{}
	}

	if cfg.Env == nil {
		cfg.Env = map[string]string{}
	}

	env := []string{}
	for k, v := range cfg.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	mcp_client, err := client.NewStdioMCPClient(cfg.Path, env, cfg.Args...)
	if err != nil {
		common.Throw(rt, fmt.Errorf("unable to create MCP client: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "k6",
		Version: "1.0.0",
	}

	_, err = mcp_client.Initialize(ctx, initRequest)
	if err != nil {
		common.Throw(rt, fmt.Errorf("unable to initialize MCP client: %w", err))
	}

	client := &Client{
		mcp_client: mcp_client,
	}

	return rt.ToValue(client).ToObject(rt)
}

// Ping checks if the MCP server is alive
func (c *Client) Ping() bool {
	err := c.mcp_client.Ping(context.Background())
	if err != nil {
		return false
	}
	return true
}

// ListTools returns available tools in MCP
func (c *Client) ListTools(r mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	tools, err := c.mcp_client.ListTools(context.Background(), r)
	if err != nil {
		return &mcp.ListToolsResult{}, err
	}
	return tools, nil
}

// ListResources returns available resources in MCP
func (c *Client) ListResources(r mcp.ListResourcesRequest) (*mcp.ListResourcesResult, error) {
	resources, err := c.mcp_client.ListResources(context.Background(), r)
	if err != nil {
		return &mcp.ListResourcesResult{}, err
	}
	return resources, nil
}

// ListPrompts returns available prompts in MCP
func (c *Client) ListPrompts(r mcp.ListPromptsRequest) (*mcp.ListPromptsResult, error) {
	prompts, err := c.mcp_client.ListPrompts(context.Background(), r)
	if err != nil {
		return &mcp.ListPromptsResult{}, err
	}
	return prompts, nil
}
