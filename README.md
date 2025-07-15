# xk6-mcp

A k6 extension for testing [Model Context Protocol (MCP)](https://modelcontextprotocol.io/introduction) servers.

## Installation

1. First, ensure you have [xk6](https://github.com/grafana/xk6) installed:
```bash
go install go.k6.io/xk6/cmd/xk6@latest
```

2. Build a k6 binary with the xk6-mcp extension:
```bash
xk6 build --with github.com/dgzlopes/xk6-mcp
```

3. Import the mcp module in your script, at the top of your test script:
```javascript
import mcp from 'k6/x/mcp';
```

4. The built binary will be in your current directory. You can move it to your PATH or use it directly:
```bash
./k6 run script.js
```

## Example

> ⚠️ This example depends on [mcp-example-server](https://github.com/dgzlopes/mcp-example-server). You can download the latest version from the releases page.

```javascript
import mcp from 'k6/x/mcp';

export default function () {
  // Initialize MCP Client with stdio transport
  const client = new mcp.StdioClient({
    path: './mcp-example-server',
  });

  // Check connection to MCP server
  console.log('MCP server running:', client.ping());

  // List all available tools
  console.log('Tools available:');
  const tools = client.listAllTools().tools;
  tools.forEach(tool => console.log(`  - ${tool.name}`));

  // List all available resources
  console.log('Resources available:');
  const resources = client.listAllResources().resources;
  resources.forEach(resource => console.log(`  - ${resource.uri}`));

  // List all available prompts
  console.log('Prompts available:');
  const prompts = client.listAllPrompts().prompts;
  prompts.forEach(prompt => console.log(`  - ${prompt.name}`));

  // Call a sample tool
  const toolResult = client.callTool({ name: 'greet', arguments: { name: 'Grafana k6' } });
  console.log(`Greet tool response: "${toolResult.content[0].text}"`);

  // Read a sample resource
  const resourceContent = client.readResource({ uri: 'embedded:info' });
  console.log(`Resource content: ${resourceContent.contents[0].text}`);

  // Get a sample prompt
  const prompt = client.getPrompt({ name: 'greet' });
  console.log(`Prompt: ${prompt.messages[0].content.text}`);
}
```

### What about non-stdio transports?

Yeah, SSE and Streamable HTTP are supported too! 

You can use them like this:

```javascript
// SSE
const client = new mcp.SSEClient({
    base_url: 'http://localhost:3002',
});

// Streamable HTTP
const client = new mcp.StreamableHTTPClient({
    base_url: 'http://localhost:3001',
});
```