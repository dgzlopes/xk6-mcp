import mcp from 'k6/x/mcp';

const client = new mcp.StdioClient({
  path: 'npx',
  env: [],
  args: ['-y', '@modelcontextprotocol/server-everything', '/tmp'],
});

export default function () {
  console.log('Is MCP server running?', client.ping());

  console.log('Tools available:');
  var tools = client.listTools().tools;
  for (var i = 0; i < tools.length; i++) {
    console.log(`\t`,tools[i].name);
  }

  console.log('Resources available:');
  var resources = client.listResources().resources;
  for (var i = 0; i < resources.length; i++) {
    console.log(`\t`,resources[i].uri);
  }

  console.log('Prompts available:');
  var prompts = client.listPrompts().prompts;
  for (var i = 0; i < prompts.length; i++) {
    console.log(`\t`,prompts[i].name);
  }
}