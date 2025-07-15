import mcp from 'k6/x/mcp';

export default function () {
  // Initialize MCP Client with stdio transport
  const client = new mcp.StdioClient({
    path: './mcp-example-server',
    args: ['--stdio'],
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