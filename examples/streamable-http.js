import mcp from 'k6/x/mcp';

export default function () {
    const client = new mcp.StreamableHTTPClient({
        base_url: 'http://localhost:3001',
    });

    console.log('MCP (http) server running:', client.ping());

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
