package dbhook

import "database/sql/driver"

func Wrap(drv driver.Driver, hks Hook) driver.Driver {
	return &Driver{Driver: drv, hooks: hks}
}
