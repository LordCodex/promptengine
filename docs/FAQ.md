# FAQ & Troubleshooting

Find answers to common questions and diagnostics steps below.

---

## 1. Why am I seeing `path traversal detected`?

**Cause**: A command or configuration referenced a directory or file outside the current PromptEngine project boundaries.
**Fix**: Ensure all relative paths resolve strictly within the directory containing your `playbook-manifest.json`.

---

## 2. Telemetry opt-out: is my code sent to remote servers?

**No**. PromptEngine collects absolutely zero code, file contents, secrets, or folder names.
- By default, all telemetry is **completely disabled**.
- When enabled, tracking is restricted to CLI execution metrics (command name, duration).
- You can ensure tracking is off by verifying the `user_consent` key is `false` in `~/.promptengine/config.json`.

---

## 3. How do I upgrade my playbook templates?

Run the `update` command to pull the latest conventions:

```bash
promptengine update
```

To review recommended updates before applying:
```bash
promptengine update --dry-run
```

---

## 4. Troubleshooting config drifts

If the `doctor` command reports config drifts:
1. Run `promptengine doctor --fix` to restore default configurations.
2. Verify your `playbook-manifest.json` does not contain syntax errors:
   ```bash
   python3 -m json.tool playbook-manifest.json
   ```
3. Run `promptengine scan` to rebuild the discovery stack settings.
