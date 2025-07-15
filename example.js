import mcp from 'k6/x/mcp';

const clientType = __ENV.MCP_CLIENT || 'stdio';

function createClient() {
    if (clientType === 'stdio') {
        return new mcp.StdioClient({
            path: './mcp-example-server',
            args: ['--stdio'],
        });
    }
    if (clientType === 'sse') {
        return new mcp.SSEClient({
            base_url: 'http://localhost:3002',
        });
    }
    if (clientType === 'http') {
        return new mcp.StreamableHTTPClient({
            base_url: 'http://localhost:3001',
        });
    }
    throw new Error(`Unknown MCP_CLIENT: ${clientType}`);
}

export default function () {
    const client = createClient();

    console.log(`MCP (${clientType}) server running:`, client.ping());

    console.log('Tools available:');
    const tools = client.listAllTools().tools;
    tools.forEach(tool => console.log(`  - ${tool.name}`));

    console.log('Resources available:');
    const resources = client.listAllResources().resources;
    resources.forEach(resource => console.log(`  - ${resource.uri}`));

    console.log('Prompts available:');
    const prompts = client.listAllPrompts().prompts;
    prompts.forEach(prompt => console.log(`  - ${prompt.name}`));

    const toolResult = client.callTool({ name: 'greet', arguments: { name: 'Grafana k6' } });
    console.log(`Greet tool response: "${toolResult.content[0].text}"`);

    const resourceContent = client.readResource({ uri: 'embedded:info' });
    console.log(`Resource content: ${resourceContent.contents[0].text}`);

    const prompt = client.getPrompt({ name: 'greet' });
    console.log(`Prompt: ${prompt.messages[0].content.text}`);
}
