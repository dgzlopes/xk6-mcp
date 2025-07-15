import mcp from 'k6/x/mcp';

export default function () {
  // Initialize MCP Client with streamable HTTP transport
  const client = new mcp.StreamableHTTPClient({
    base_url: 'http://localhost:3001', // adjust as needed
  });

  // Check connection to MCP HTTP server
  console.log('MCP HTTP server running:', client.ping());

  // List all available tools (HTTP)
  console.log('Tools available (HTTP):');
  const tools = client.listAllTools().tools;
  tools.forEach(tool => console.log(`  - ${tool.name}`));

  // List all available resources (HTTP)
  console.log('Resources available (HTTP):');
  const resources = client.listAllResources().resources;
  resources.forEach(resource => console.log(`  - ${resource.uri}`));

  // List all available prompts (HTTP)
  console.log('Prompts available (HTTP):');
  const prompts = client.listAllPrompts().prompts;
  prompts.forEach(prompt => console.log(`  - ${prompt.name}`));

  // Call a sample tool (HTTP)
  const toolResult = client.callTool({ name: 'greet', arguments: { name: 'Grafana k6' } });
  console.log(`Greet tool response (HTTP): "${toolResult.content[0].text}"`);

  // Read a sample resource (HTTP)
  const resourceContent = client.readResource({ uri: 'embedded:info' });
  console.log(`Resource content (HTTP): ${resourceContent.contents[0].text}`);

  // Get a sample prompt (HTTP)
  const prompt = client.getPrompt({ name: 'greet' });
  console.log(`Prompt (HTTP): ${prompt.messages[0].content.text}`);
}
