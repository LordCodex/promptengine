const { execFile } = require("child_process");

class PromptEngineClient {
  constructor(config, exec = execFile) {
    this.config = config;
    this.exec = exec;
  }

  analyzeProject() {
    return this.run(["scan", "--json"]);
  }

  generateContext(task, intent) {
    return this.run([
      "context",
      "--json",
      "--task",
      task || "feature",
      "--intent",
      intent || "",
      "--max-bytes",
      String(this.config.contextLimitBytes)
    ]);
  }

  exportContext(task, intent, agent) {
    return this.run([
      "context",
      "export",
      "--json",
      "--task",
      task || "feature",
      "--intent",
      intent || "",
      "--agent",
      agent || this.config.preferredAIClient,
      "--format",
      this.config.outputFormat,
      "--max-bytes",
      String(this.config.contextLimitBytes)
    ]);
  }

  generatePrompt(task, request, agent) {
    return this.run([
      "prompt",
      "--json",
      "--task",
      task || "feature",
      "--request",
      request || "",
      "--client",
      agent || this.config.preferredAIClient,
      "--format",
      this.config.outputFormat,
      "--max-bytes",
      String(this.config.contextLimitBytes)
    ]);
  }

  runWorkflow(id) {
    return this.run(["workflow", "--json", "--id", id || "feature-implementation"]);
  }

  checkHealth() {
    return this.run(["health", "--json"]);
  }

  syncDocumentation() {
    return this.run(["docs", "sync", "--json"]);
  }

  run(args) {
    const fullArgs = this.config.configPath ? ["--config", this.config.configPath, ...args] : args;
    return new Promise((resolve, reject) => {
      this.exec(
        this.config.binaryPath,
        fullArgs,
        { cwd: this.config.workspaceRoot, maxBuffer: 10 * 1024 * 1024 },
        (error, stdout, stderr) => {
          if (error) {
            reject(new Error(`${error.message}${stderr ? `\n${stderr}` : ""}`));
            return;
          }
          resolve({ stdout: stdout.toString(), stderr: stderr.toString() });
        }
      );
    });
  }
}

function summarizeResult(stdout) {
  const trimmed = stdout.trim();
  if (!trimmed) {
    return "PromptEngine command completed.";
  }
  try {
    const parsed = JSON.parse(trimmed);
    if (parsed.file) {
      return `PromptEngine generated ${parsed.file}.`;
    }
    if (parsed.project?.root_path || parsed.root_dir) {
      return `PromptEngine analyzed ${parsed.project?.root_path || parsed.root_dir}.`;
    }
    if (parsed.score?.overall !== undefined) {
      return `PromptEngine health score: ${parsed.score.overall}/100.`;
    }
    if (parsed.status) {
      return `PromptEngine status: ${parsed.status}.`;
    }
  } catch {
    // Plain text or YAML output.
  }
  return trimmed.split(/\r?\n/)[0];
}

module.exports = { PromptEngineClient, summarizeResult };
