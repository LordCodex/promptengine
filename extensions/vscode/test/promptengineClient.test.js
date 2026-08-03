const assert = require("assert");
const { PromptEngineClient, summarizeResult } = require("../src/promptengineClient");

function fakeExec(calls) {
  return (file, args, options, callback) => {
    calls.push([file, ...args]);
    callback(null, '{"file":"feature-prompt.md"}', "");
    return {};
  };
}

async function run() {
  const calls = [];
  const client = new PromptEngineClient({
    binaryPath: "promptengine",
    configPath: ".promptengine.yaml",
    workspaceRoot: "/repo",
    preferredAIClient: "codex",
    outputFormat: "markdown",
    contextLimitBytes: 1234
  }, fakeExec(calls));

  await client.generatePrompt("feature", "Add billing");
  assert.deepStrictEqual(calls[0], [
    "promptengine",
    "--config",
    ".promptengine.yaml",
    "prompt",
    "--json",
    "--task",
    "feature",
    "--request",
    "Add billing",
    "--client",
    "codex",
    "--format",
    "markdown",
    "--max-bytes",
    "1234"
  ]);

  await client.exportContext("bug_fix", "Current file: src/app.ts", "claude");
  assert.ok(calls[1].includes("context"));
  assert.ok(calls[1].includes("export"));
  assert.ok(calls[1].includes("claude"));
  assert.strictEqual(summarizeResult('{"score":{"overall":88}}'), "PromptEngine health score: 88/100.");
}

run().catch((error) => {
  console.error(error);
  process.exit(1);
});
