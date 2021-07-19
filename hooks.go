package dbhook

import (
	"context"
	"database/sql/driver"
	"fmt"
)

type HookBefore interface {
	Before(ctx context.Context, input *HookInput) (context.Context, error)
}

type HookAfter interface {
	After(ctx context.Context, input *HookInput) (context.Context, error)
}

type HookError interface {
	Error(ctx context.Context, input *HookInput) (context.Context, error)
}

type Hook interface {
	HookBefore
	HookAfter
	HookError
}

type HookInput struct {
	Error  error
	Args   []driver.Value
	Query  string
	Caller CallerType
}

type HookOption func(*Hooks)

type Hooks struct {
	before []HookBefore
	after  []HookAfter
	err    []HookError
}

func NewHooks(opts ...HookOption) *Hooks {
	hooks := &Hooks{}

	for _, opt := range opts {
		opt(hooks)
	}

	return hooks
}

func WithHook(hooks ...Hook) HookOption {
	return func(h *Hooks) {
		for _, hook := range hooks {
			h.before = append(h.before, hook)
			h.after = append(h.after, hook)
			h.err = append(h.err, hook)
		}
	}
}

func WithHooksBefore(hooks ...HookBefore) HookOption {
	return func(h *Hooks) {
		switch {
		case len(h.before) == 0:
			h.before = hooks
		default:
			h.before = append(h.before, hooks...)
		}
	}
}

func WithHooksAfter(hooks ...HookAfter) HookOption {
	return func(h *Hooks) {
		switch {
		case len(h.after) == 0:
			h.after = hooks
		default:
			h.after = append(h.after, hooks...)
		}
	}
}

func WithHooksError(hooks ...HookError) HookOption {
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
		ctx, err = h.before[i].Before(ctx, input)
		if err != nil {
			return ctx, fmt.Errorf("before hook #%d call failed: %w", i, err)
		}
	}

	return ctx, nil
}

func (h *Hooks) After(ctx context.Context, input *HookInput) (context.Context, error) {
	var err error

	for i := range h.after {
		ctx, err = h.after[i].After(ctx, input)
		if err != nil {
			return ctx, fmt.Errorf("after hook #%d call failed: %w", i, err)
		}
	}

	return ctx, nil
}

func (h *Hooks) Error(ctx context.Context, input *HookInput) (context.Context, error) {
	for i := range h.err {
		ctx, input.Error = h.err[i].Error(ctx, input)
		if input.Error == nil {
			return ctx, nil
		}
	}

	return ctx, input.Error
}
