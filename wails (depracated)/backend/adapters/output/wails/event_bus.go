package wails

import (
	"context"
	"pdr/backend/core/renderer"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type eventBus struct{}

func NewEventBus() *eventBus {
	return &eventBus{}
}

const (
	RenderPorgessEventKey = "render_progress"
)

func (*eventBus) EmitRenderProgress(ctx context.Context, payload renderer.ProgressEventPayload) {
	runtime.EventsEmit(ctx, RenderPorgessEventKey, payload)
}
