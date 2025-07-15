import mcp from 'k6/x/mcp';

export default function () {
  // Initialize MCP Client with sse transport
  const client = new mcp.SSEClient({
    base_url: 'http://localhost:3002',
  });

  console.log('Checking MCP server status...');
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
  const toolResult = client.callTool({ name: 'echo', arguments: { message: 'Hello, world!' } });
  console.log(`Echo tool response: "${toolResult.content[0].text}"`);

  // Read a sample resource
  const resourceContent = client.readResource({ uri: 'test://static/resource/1' });
  console.log(`Resource content: ${resourceContent.contents[0].text}`);

  // Get a sample prompt
  const prompt = client.getPrompt({ name: 'simple_prompt' });
  console.log(`Prompt: ${prompt.messages[0].content.text}`);
}