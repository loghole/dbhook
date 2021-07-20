package dbhook

import (
	"context"
	"database/sql/driver"
)

type Tx struct {
	Tx    driver.Tx
	hooks Hook
	ctx   context.Context
}

func (tx *Tx) Commit() error {
	var (
		err       error
		hookInput = &HookInput{
			Caller: CallerCommit,
			Query:  "",
			Args:   nil,
			Error:  nil,
		}
	)

	if tx.hooks != nil {
		tx.ctx, err = tx.hooks.Before(tx.ctx, hookInput)
		if err != nil {
			return err
		}
	}

	err = tx.Tx.Commit()
	if err != nil {
		if tx.hooks == nil {
			return err
		}

		hookInput.Error = err

		if tx.ctx, err = tx.hooks.Error(tx.ctx, hookInput); err != nil {
			return err
		}
	}

	if tx.hooks != nil {
		if tx.ctx, err = tx.hooks.After(tx.ctx, hookInput); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Tx) Rollback() error {
	var (
		err       error
		hookInput = &HookInput{
			Caller: CallerRollback,
			Query:  "",
			Args:   nil,
			Error:  nil,
		}
	)

	if tx.hooks != nil {
		tx.ctx, err = tx.hooks.Before(tx.ctx, hookInput)
		if err != nil {
			return err
		}
	}

	err = tx.Tx.Rollback()
	if err != nil {
		if tx.hooks == nil {
			return err
		}

		hookInput.Error = err

		if tx.ctx, err = tx.hooks.Error(tx.ctx, hookInput); err != nil {
			return err
		}
	}

	if tx.hooks != nil {
		if tx.ctx, err = tx.hooks.After(tx.ctx, hookInput); err != nil {
			return err
		}
	}

	return nil
}
