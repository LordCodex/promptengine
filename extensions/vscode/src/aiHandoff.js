const path = require("path");

async function readGeneratedPrompt(vscode, workspaceRoot, stdout) {
  let payload;
  try {
    payload = JSON.parse(stdout.trim());
  } catch (error) {
    throw new Error(`PromptEngine returned an invalid prompt result: ${error.message}`);
  }

  if (!payload.file) {
    throw new Error("PromptEngine did not return a generated prompt file.");
  }

  const promptPath = path.isAbsolute(payload.file)
    ? payload.file
    : path.join(workspaceRoot, payload.file);
  const bytes = await vscode.workspace.fs.readFile(vscode.Uri.file(promptPath));
  return Buffer.from(bytes).toString("utf8");
}

async function handoffPrompt(vscode, promptText, preferredClient) {
  const client = String(preferredClient || "generic").toLowerCase();
  const commands = new Set(await vscode.commands.getCommands(true));

  if (client === "claude" && commands.has("claude-vscode.editor.open")) {
    await vscode.commands.executeCommand(
      "claude-vscode.editor.open",
      undefined,
      promptText,
      vscode.ViewColumn.Active
    );
    vscode.window.showInformationMessage(
      "PromptEngine filled the Claude composer. Review it and press Enter when ready."
    );
    return { mode: "prefill", client: "claude" };
  }

  if (client === "codex" && commands.has("chatgpt.addToThread")) {
    const originalEditor = vscode.window.activeTextEditor;
    const document = await vscode.workspace.openTextDocument({
      content: promptText,
      language: "markdown"
    });
    const editor = await vscode.window.showTextDocument(document, { preview: true });
    const start = document.positionAt(0);
    const end = document.positionAt(document.getText().length);
    editor.selection = new vscode.Selection(start, end);

    await vscode.commands.executeCommand("chatgpt.addToThread");
    await vscode.commands.executeCommand("workbench.action.closeActiveEditor");

    if (commands.has("workbench.view.extension.codexSecondaryViewContainer")) {
      await vscode.commands.executeCommand("workbench.view.extension.codexSecondaryViewContainer");
    } else if (commands.has("workbench.view.extension.codexViewContainer")) {
      await vscode.commands.executeCommand("workbench.view.extension.codexViewContainer");
    } else if (originalEditor) {
      await vscode.window.showTextDocument(originalEditor.document, {
        selection: originalEditor.selection,
        preserveFocus: true
      });
    }

    vscode.window.showInformationMessage(
      "PromptEngine added the generated task to the Codex thread. Review it and send when ready."
    );
    return { mode: "thread-context", client: "codex" };
  }

  await vscode.env.clipboard.writeText(promptText);
  vscode.window.showInformationMessage(
    `PromptEngine could not directly prefill ${client}. The generated prompt was copied to your clipboard.`
  );
  return { mode: "clipboard", client };
}

module.exports = { handoffPrompt, readGeneratedPrompt };
