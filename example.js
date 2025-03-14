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
  const result = client.callTool({
    params: { name: 'echo', arguments: { message: 'Hello, world!' } }
  });
  console.log('Echo tool response:', result.content[0].text);
}