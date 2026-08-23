const assert = require("assert");
const { handoffPrompt } = require("../src/aiHandoff");

function fakeVscode(commands = []) {
  const calls = [];
  let clipboard = "";
  return {
    calls,
    ViewColumn: { Active: 1 },
    commands: {
      getCommands: async () => commands,
      executeCommand: async (...args) => {
        calls.push(args);
      }
    },
    window: {
      showInformationMessage: () => {},
      activeTextEditor: undefined,
      showTextDocument: async () => {
        throw new Error("showTextDocument should not be called in this test");
      }
    },
    workspace: {
      openTextDocument: async () => {
        throw new Error("openTextDocument should not be called in this test");
      }
    },
    env: {
      clipboard: {
        writeText: async (value) => {
          clipboard = value;
        }
      }
    },
    clipboard: () => clipboard
  };
}

async function run() {
  const claude = fakeVscode(["claude-vscode.editor.open"]);
  const claudeResult = await handoffPrompt(claude, "Implement billing", "claude");
  assert.strictEqual(claudeResult.mode, "prefill");
  assert.strictEqual(claude.calls.length, 1);
  assert.strictEqual(claude.calls[0][0], "claude-vscode.editor.open");
  assert.strictEqual(claude.calls[0][2], "Implement billing");

  const generic = fakeVscode([]);
  const fallbackResult = await handoffPrompt(generic, "Fix webhook", "codex");
  assert.strictEqual(fallbackResult.mode, "clipboard");
  assert.strictEqual(generic.clipboard(), "Fix webhook");
}

run().catch((error) => {
  console.error(error);
  process.exit(1);
});
