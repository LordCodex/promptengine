package workflows

import (
	"context"

	"github.com/LordCodex/promptengine/internal/domain/workflows"
	"github.com/LordCodex/promptengine/internal/eventbus"
	"github.com/LordCodex/promptengine/internal/filesystem"
)

type State = workflows.State
type FlowContext = workflows.FlowContext

type Engine struct {
	inner *workflows.Engine
}

func NewEngine(fs filesystem.FileSystem, eb *eventbus.EventBus) *Engine {
	reg := workflows.NewRegistry()
	return &Engine{
		inner: workflows.NewEngine(fs, reg, eb),
	}
}

func (e *Engine) RunWorkflow(ctx context.Context, flowName string, flowCtx *FlowContext) (State, error) {
	return e.inner.RunWorkflow(ctx, flowName, flowCtx)
}
