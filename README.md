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

```javascript
import mcp from 'k6/x/mcp';

// Initialize MCP Client
const client = new mcp.StdioClient({
  path: 'npx',
  env: [],
  args: ['-y', '@modelcontextprotocol/server-everything', '/tmp'],
});

export default function () {
  console.log('Checking MCP server status...');
  console.log('MCP server running:', client.ping());

  // List available tools
  console.log('Tools available:');
  const tools = client.listTools().tools;
  tools.forEach(tool => console.log(`  - ${tool.name}`));

  // List available resources
  console.log('Resources available:');
  const resources = client.listResources().resources;
  resources.forEach(resource => console.log(`  - ${resource.uri}`));

  // List available prompts
  console.log('Prompts available:');
  const prompts = client.listPrompts().prompts;
  prompts.forEach(prompt => console.log(`  - ${prompt.name}`));

  // Call a sample tool
  const toolResult = client.callTool({
    params: { name: 'echo', arguments: { message: 'Hello, world!' } }
  });
  console.log('Echo tool response:', toolResult.content[0].text);

  // Read a sample resource
  const resourceContent = client.readResource({
    params: { uri: 'test://static/resource/1' }
  });
  console.log('Resource content:', resourceContent.contents[0].text);

  // Get a sample prompt
  const prompt = client.getPrompt({
    params: { name: 'simple_prompt' }
  });
  console.log('Prompt:', prompt.messages[0].content.text);
}
```