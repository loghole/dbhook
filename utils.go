package dbhook

import (
	"database/sql/driver"
)

func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))

	for n, param := range named {
		if param.Name != "" {
			return nil, ErrNamedParametersNotSupport
		}

		dargs[n] = param.Value
	}

	return dargs, nil
}

func argsToValue(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))

	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}

	return values
}
