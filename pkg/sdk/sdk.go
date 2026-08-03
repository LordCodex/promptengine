package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type OutputFormat string

const (
	OutputJSON     OutputFormat = "json"
	OutputYAML     OutputFormat = "yaml"
	OutputMarkdown OutputFormat = "markdown"
)

type Options struct {
	BinaryPath string
	WorkingDir string
	ConfigPath string
}

type Client struct {
	binaryPath string
	workingDir string
	configPath string
	runner     CommandRunner
}

type CommandRunner interface {
	Run(ctx context.Context, workingDir string, name string, args ...string) ([]byte, []byte, error)
}

type osRunner struct{}

func NewClient(opts Options) *Client {
	binary := opts.BinaryPath
	if strings.TrimSpace(binary) == "" {
		binary = "promptengine"
	}
	return &Client{binaryPath: binary, workingDir: opts.WorkingDir, configPath: opts.ConfigPath, runner: osRunner{}}
}

func NewClientWithRunner(opts Options, runner CommandRunner) *Client {
	c := NewClient(opts)
	if runner != nil {
		c.runner = runner
	}
	return c
}

func (c *Client) AnalyzeProject(ctx context.Context) ([]byte, error) {
	return c.run(ctx, "scan", "--json")
}

func (c *Client) GenerateContext(ctx context.Context, req ContextRequest) ([]byte, error) {
	args := []string{"context", "--json", "--task", defaultString(req.Task, "feature")}
	if req.Intent != "" {
		args = append(args, "--intent", req.Intent)
	}
	if req.Budget != "" {
		args = append(args, "--budget", req.Budget)
	}
	if req.MaxBytes > 0 {
		args = append(args, "--max-bytes", fmt.Sprintf("%d", req.MaxBytes))
	}
	return c.run(ctx, args...)
}

func (c *Client) ExportContext(ctx context.Context, req ContextRequest) ([]byte, error) {
	args := []string{"context", "export", "--json", "--task", defaultString(req.Task, "feature"), "--agent", defaultString(req.Agent, "generic"), "--format", defaultString(req.Format, "markdown")}
	if req.Intent != "" {
		args = append(args, "--intent", req.Intent)
	}
	if req.Budget != "" {
		args = append(args, "--budget", req.Budget)
	}
	if req.MaxBytes > 0 {
		args = append(args, "--max-bytes", fmt.Sprintf("%d", req.MaxBytes))
	}
	return c.run(ctx, args...)
}

func (c *Client) GeneratePrompt(ctx context.Context, req PromptRequest) ([]byte, error) {
	args := []string{"prompt", "--json", "--task", defaultString(req.Task, "feature"), "--client", defaultString(req.Agent, "generic"), "--format", defaultString(req.Format, "markdown")}
	if req.Request != "" {
		args = append(args, "--request", req.Request)
	}
	if req.OutputPath != "" {
		args = append(args, "--out", req.OutputPath)
	}
	if req.Budget != "" {
		args = append(args, "--budget", req.Budget)
	}
	if req.MaxBytes > 0 {
		args = append(args, "--max-bytes", fmt.Sprintf("%d", req.MaxBytes))
	}
	return c.run(ctx, args...)
}

func (c *Client) RunWorkflow(ctx context.Context, id string) ([]byte, error) {
	return c.run(ctx, "workflow", "--json", "--id", defaultString(id, "feature-implementation"))
}

func (c *Client) SyncDocumentation(ctx context.Context) ([]byte, error) {
	return c.run(ctx, "docs", "sync", "--json")
}

func (c *Client) CheckHealth(ctx context.Context) ([]byte, error) {
	return c.run(ctx, "health", "--json")
}

type ContextRequest struct {
	Task     string
	Intent   string
	Agent    string
	Format   string
	Budget   string
	MaxBytes int
}

type PromptRequest struct {
	Task       string
	Request    string
	Agent      string
	Format     string
	OutputPath string
	Budget     string
	MaxBytes   int
}

func DecodeJSON[T any](data []byte) (T, error) {
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}

func (c *Client) run(ctx context.Context, args ...string) ([]byte, error) {
	if c.configPath != "" {
		args = append([]string{"--config", c.configPath}, args...)
	}
	stdout, stderr, err := c.runner.Run(ctx, c.workingDir, c.binaryPath, args...)
	if err != nil {
		return nil, fmt.Errorf("promptengine cli failed: %w: %s", err, strings.TrimSpace(string(stderr)))
	}
	return stdout, nil
}

func (osRunner) Run(ctx context.Context, workingDir string, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
