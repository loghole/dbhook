package dbhook

import (
	"context"
	"database/sql/driver"
	"fmt"
)

type HookCall interface {
	Call(ctx context.Context, input *HookInput) (context.Context, error)
}

type Hook interface {
	Before(ctx context.Context, input *HookInput) (context.Context, error)
	After(ctx context.Context, input *HookInput) (context.Context, error)
	Error(ctx context.Context, input *HookInput) (context.Context, error)
}

type HookInput struct {
	Caller CallerType
	Query  string
	Args   []driver.Value
	Error  error
}

type HookOption func(*Hooks)

type Hooks struct {
	before []HookCall
	after  []HookCall
	err    []HookCall
}

func NewHooks(opts ...HookOption) *Hooks {
	hooks := &Hooks{}

	for _, opt := range opts {
		opt(hooks)
	}

	return hooks
}

func WithHooksBefore(hooks ...HookCall) HookOption {
	return func(h *Hooks) {
		switch {
		case len(h.before) == 0:
			h.before = hooks
		default:
			h.before = append(h.before, hooks...)
		}
	}
}

func WithHooksAfter(hooks ...HookCall) HookOption {
	return func(h *Hooks) {
		switch {
		case len(h.after) == 0:
			h.after = hooks
		default:
			h.after = append(h.after, hooks...)
		}
	}
}

func WithHooksError(hooks ...HookCall) HookOption {
	return func(h *Hooks) {
		switch {
		case len(h.err) == 0:
			h.err = hooks
		default:
			h.err = append(h.err, hooks...)
		}
	}
}

func (h *Hooks) Before(ctx context.Context, input *HookInput) (context.Context, error) {
	var err error

	for i := range h.before {
		ctx, err = h.before[i].Call(ctx, input)
		if err != nil {
			return ctx, fmt.Errorf("before hook #%d call failed: %w", i, err)
		}
	}

	return ctx, nil
}

func (h *Hooks) After(ctx context.Context, input *HookInput) (context.Context, error) {
	var err error

	for i := range h.after {
		ctx, err = h.after[i].Call(ctx, input)
		if err != nil {
			return ctx, fmt.Errorf("after hook #%d call failed: %w", i, err)
		}
	}

	return ctx, nil
}

func (h *Hooks) Error(ctx context.Context, input *HookInput) (context.Context, error) {
	err := input.Error

	for i := range h.after {
		ctx, err = h.after[i].Call(ctx, input)
		if err != nil {
			return ctx, err //nolint:wrapcheck // need clear error
		}
	}

	return ctx, err //nolint:wrapcheck // need clear error
}
