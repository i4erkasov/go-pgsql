package go_pgsql

import "github.com/i4erkasov/go-pgsql/pgx"

// Deprecated: in mockery v3 this approach will be removed
// rename the file .mockery.yaml.dist to .mockery.yaml and remove go:generate
// there's a problem with aliased type: for aliased types mocks won't be generated
// type Foo = baz.Baz - this mock won't be generated
// see https://github.com/vektra/mockery/issues with replace-type params
//go:generate mockery --config=.mockery.yaml

type TxMock interface {
	pgx.Tx
}
