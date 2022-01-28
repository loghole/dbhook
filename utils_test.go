package dbhook

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func BenchmarkNamedValueToValue(b *testing.B) {
	named := make([]driver.NamedValue, 100)
	for i := range named {
		named[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   fmt.Sprintf("it's number: %d", i),
		}
	}

	b.Run("namedValueToValue", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			namedValueToValue(named)
		}
	})

	b.Run("argsToValue", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			argsToValue(named)
		}
	})
}

func Test_namedValueToValue(t *testing.T) {
	t.Parallel()

	named := make([]driver.NamedValue, 10)
	for i := range named {
		named[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   fmt.Sprintf("it's number: %d", i),
		}
	}

	type args struct {
		named []driver.NamedValue
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pass1",
			args: args{
				named: named,
			},
		},
		{
			name: "pass2",
			args: args{
				named: nil,
			},
		},
		{
			name: "has_error",
			args: args{
				named: []driver.NamedValue{
					{
						Ordinal: 1,
						Value:   1,
					},
					{
						Name:    "name",
						Ordinal: 2,
						Value:   2,
					},
				},
			},
			wantErr: true,
		},
	}

	snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory(testDataPath))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := namedValueToValue(tt.args.named)
			if (err != nil) != tt.wantErr {
				t.Errorf("namedValueToValue() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			snapshotter.SnapshotT(t, got)
		})
	}
}

func Test_argsToValue(t *testing.T) {
	t.Parallel()

	values := make([]driver.NamedValue, 10)
	for i := range values {
		values[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   fmt.Sprintf("it's number: %d", i),
		}
	}

	type args struct {
		args []driver.NamedValue
	}

	tests := []struct {
		name string
		args args
		want []driver.Value
	}{
		{
			name: "pass1",
			args: args{
				args: values,
			},
		},
		{
			name: "pass2",
			args: args{
				args: nil,
			},
		},
		{
			name: "with_name",
			args: args{
				args: []driver.NamedValue{
					{
						Ordinal: 1,
						Value:   1,
					},
					{
						Name:    "name",
						Ordinal: 2,
						Value:   2,
					},
				},
			},
		},
	}

	snapshotter := cupaloy.New(cupaloy.SnapshotSubdirectory(testDataPath))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := argsToValue(tt.args.args)

			snapshotter.SnapshotT(t, got)
		})
	}
}
