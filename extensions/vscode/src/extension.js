const vscode = require("vscode");
const { PromptEngineClient, summarizeResult } = require("./promptengineClient");

function activate(context) {
  const register = (command, handler) => {
    context.subscriptions.push(vscode.commands.registerCommand(command, () => runWithErrors(handler)));
  };

  register("promptengine.analyzeProject", async () => {
    await showResult((await client().analyzeProject()).stdout);
  });

  register("promptengine.generateContext", async () => {
    const intent = await prompt("Describe the task for context generation");
    await showResult((await client().generateContext("feature", intent || "")).stdout);
  });

  register("promptengine.generatePrompt", async () => {
    const request = await prompt("Describe the implementation request");
    await showResult((await client().generatePrompt("feature", request || "")).stdout);
  });

  register("promptengine.runWorkflow", async () => {
    const id = await prompt("Workflow id", "feature-implementation");
    await showResult((await client().runWorkflow(id || "feature-implementation")).stdout);
  });

  register("promptengine.checkHealth", async () => {
    await showResult((await client().checkHealth()).stdout);
  });

  register("promptengine.syncDocumentation", async () => {
    await showResult((await client().syncDocumentation()).stdout);
  });

  register("promptengine.exportSelectedCodeContext", async () => {
    const text = selectedText();
    const intent = text ? `Selected code:\n${text}` : "Generate context for selected code.";
    await showResult((await client().exportContext("feature", intent)).stdout);
  });

  register("promptengine.exportCurrentFileContext", async () => {
    const file = activeFilePath();
    await showResult((await client().exportContext("bug_fix", file ? `Current file: ${file}` : "Analyze current file.")).stdout);
  });
}

function deactivate() {}

function client() {
  return new PromptEngineClient(readConfig());
}

function readConfig() {
  const cfg = vscode.workspace.getConfiguration("promptengine");
  return {
    binaryPath: cfg.get("path", "promptengine"),
    configPath: cfg.get("configPath", ""),
    workspaceRoot: workspaceRoot(),
    preferredAIClient: cfg.get("preferredAIClient", "codex"),
    outputFormat: cfg.get("outputFormat", "markdown"),
    contextLimitBytes: cfg.get("contextLimitBytes", 100000)
  };
}

function workspaceRoot() {
  const folders = vscode.workspace.workspaceFolders;
  if (!folders || folders.length === 0) {
    throw new Error("Open a workspace before running PromptEngine.");
  }
  return folders[0].uri.fsPath;
}

function selectedText() {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    return "";
  }
  return editor.document.getText(editor.selection);
}

function activeFilePath() {
  return vscode.window.activeTextEditor?.document.uri.fsPath || "";
}

function prompt(placeHolder, value = "") {
  return vscode.window.showInputBox({ placeHolder, value });
}

async function showResult(stdout) {
  const summary = summarizeResult(stdout);
  vscode.window.showInformationMessage(summary);
  const doc = await vscode.workspace.openTextDocument({ content: stdout, language: "json" });
  await vscode.window.showTextDocument(doc, { preview: true });
}

async function runWithErrors(handler) {
  try {
    await handler();
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    vscode.window.showErrorMessage(`PromptEngine failed: ${message}`);
  }
}

module.exports = { activate, deactivate };
