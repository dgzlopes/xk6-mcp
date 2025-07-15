package compare

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/grafana/sobek"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

// ClientConfig represents the configuration for the MCP client
type ClientConfig struct {
	// Stdio
	Path string
	Args []string
	Env  map[string]string

	// SSE
	BaseURL string
	Timeout time.Duration
}

// Client wraps an MCP client session
type Client struct {
	session *mcp.ClientSession
}

// Exports defines the JavaScript-accessible functions
func (m *MCPInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"StdioClient":          m.newStdioClient,
			"SSEClient":            m.newSSEClient,
			"StreamableHTTPClient": m.newStreamableHTTPClient,
		},
	}
}

func (m *MCPInstance) newStdioClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	if err := rt.ExportTo(c.Argument(0), &cfg); err != nil {
		common.Throw(rt, fmt.Errorf("invalid config: %w", err))
	}

	cmd := exec.Command(cfg.Path, cfg.Args...)
	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	//cmd.Stderr = os.Stderr

	transport := mcp.NewCommandTransport(cmd)

	return m.connect(rt, transport, false)
}

func (m *MCPInstance) newSSEClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	if err := rt.ExportTo(c.Argument(0), &cfg); err != nil {
		common.Throw(rt, fmt.Errorf("invalid config: %w", err))
	}

	transport := mcp.NewSSEClientTransport(cfg.BaseURL, &mcp.SSEClientTransportOptions{
		HTTPClient: m.newk6HTTPClient(),
	})

	return m.connect(rt, transport, true)
}

func (m *MCPInstance) newStreamableHTTPClient(c sobek.ConstructorCall, rt *sobek.Runtime) *sobek.Object {
	var cfg ClientConfig
	if err := rt.ExportTo(c.Argument(0), &cfg); err != nil {
		common.Throw(rt, fmt.Errorf("invalid config: %w", err))
	}

	transport := mcp.NewStreamableClientTransport(cfg.BaseURL, &mcp.StreamableClientTransportOptions{
		HTTPClient: m.newk6HTTPClient(),
	})

	return m.connect(rt, transport, true)
}

func (m *MCPInstance) newk6HTTPClient() *http.Client {
	var tlsConfig *tls.Config
	if m.vu.State().TLSConfig != nil {
		tlsConfig = m.vu.State().TLSConfig.Clone()
		tlsConfig.NextProtos = []string{"http/1.1"}
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext:       m.vu.State().Dialer.DialContext,
			Proxy:             http.ProxyFromEnvironment,
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: m.vu.State().Options.NoConnectionReuse.ValueOrZero() || m.vu.State().Options.NoVUConnectionReuse.ValueOrZero(),
		},
	}

	return httpClient
}

func (m *MCPInstance) connect(rt *sobek.Runtime, transport mcp.Transport, isSSE bool) *sobek.Object {
	var ctx context.Context
	var cancel context.CancelFunc
	if isSSE {
		ctx = context.Background()
		cancel = func() {}
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	}
	defer cancel()

	client := mcp.NewClient("k6", "1.0.0", nil)
	session, err := client.Connect(ctx, transport)
	if err != nil {
		common.Throw(rt, fmt.Errorf("connection error: %w", err))
	}

	return rt.ToValue(&Client{session: session}).ToObject(rt)
}

func (c *Client) Ping() bool {
	err := c.session.Ping(context.Background(), &mcp.PingParams{})
	return err == nil
}

func (c *Client) ListTools(r mcp.ListToolsParams) (*mcp.ListToolsResult, error) {
	return c.session.ListTools(context.Background(), &r)
}

type ListAllToolsParams struct {
	Meta mcp.Meta
}

type ListAllToolsResult struct {
	Tools []mcp.Tool
}

func (c *Client) ListAllTools(r ListAllToolsParams) (*ListAllToolsResult, error) {
	if r.Meta == nil {
		r.Meta = mcp.Meta{}
	}

	var allTools []mcp.Tool
	cursor := ""
	for {
		params := &mcp.ListToolsParams{Meta: r.Meta}
		if cursor != "" {
			params.Cursor = cursor
		}
		result, err := c.session.ListTools(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools: %w", err)
		}

		for _, t := range result.Tools {
			if t != nil {
				allTools = append(allTools, *t)
			}
		}

		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	return &ListAllToolsResult{
		Tools: allTools,
	}, nil
}

func (c *Client) CallTool(r mcp.CallToolParams) (*mcp.CallToolResult, error) {
	return c.session.CallTool(context.Background(), &r)
}

func (c *Client) ListResources(r mcp.ListResourcesParams) (*mcp.ListResourcesResult, error) {
	return c.session.ListResources(context.Background(), &r)
}

func (c *Client) ReadResource(r mcp.ReadResourceParams) (*mcp.ReadResourceResult, error) {
	return c.session.ReadResource(context.Background(), &r)
}

func (c *Client) ListPrompts(r mcp.ListPromptsParams) (*mcp.ListPromptsResult, error) {
	return c.session.ListPrompts(context.Background(), &r)
}

func (c *Client) GetPrompt(r mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return c.session.GetPrompt(context.Background(), &r)
}

type ListAllResourcesParams struct {
	Meta mcp.Meta
}

type ListAllResourcesResult struct {
	Resources []mcp.Resource
}

func (c *Client) ListAllResources(r ListAllResourcesParams) (*ListAllResourcesResult, error) {
	if r.Meta == nil {
		r.Meta = mcp.Meta{}
	}

	var allResources []mcp.Resource
	cursor := ""
	for {
		params := &mcp.ListResourcesParams{Meta: r.Meta}
		if cursor != "" {
			params.Cursor = cursor
		}
		result, err := c.session.ListResources(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("failed to list resources: %w", err)
		}

		for _, res := range result.Resources {
			if res != nil {
				allResources = append(allResources, *res)
			}
		}

		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	return &ListAllResourcesResult{
		Resources: allResources,
	}, nil
}

type ListAllPromptsParams struct {
	Meta mcp.Meta
}

type ListAllPromptsResult struct {
	Prompts []mcp.Prompt
}

func (c *Client) ListAllPrompts(r ListAllPromptsParams) (*ListAllPromptsResult, error) {
	if r.Meta == nil {
		r.Meta = mcp.Meta{}
	}

	var allPrompts []mcp.Prompt
	cursor := ""
	for {
		params := &mcp.ListPromptsParams{Meta: r.Meta}
		if cursor != "" {
			params.Cursor = cursor
		}
		result, err := c.session.ListPrompts(context.Background(), params)
		if err != nil {
			return nil, fmt.Errorf("failed to list prompts: %w", err)
		}

		for _, p := range result.Prompts {
			if p != nil {
				allPrompts = append(allPrompts, *p)
			}
		}

		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}

	return &ListAllPromptsResult{
		Prompts: allPrompts,
	}, nil
}
