package persist

import (
	"testing"

	"github.com/wothing/log"

	m "github.com/arstd/gobatis/examples/domain"
	e "github.com/arstd/gobatis/examples/enums"
)

var x ModelImplExample

func TestModelMapperAdd(t *testing.T) {
	m := &m.Model{
		BuildinBool: true,
		BuildinByte: 'x',
		// BuildinComplex128 complex128
		// BuildinComplex64  complex64
		// BuildinError   error
		BuildinFloat32: float32(1.23),
		BuildinFloat64: float64(1.2345678),
		BuildinInt:     2,
		BuildinInt16:   2,
		BuildinInt32:   2,
		BuildinInt64:   2,
		BuildinInt8:    2,
		BuildinRune:    '中',
		BuildinString:  "text",
		BuildinUint:    2,
		BuildinUint16:  2,
		BuildinUint32:  2,
		BuildinUint64:  2,
		BuildinUint8:   2,
		// BuildinUintptr: 2,
		BuildinMap: map[string]interface{}{"a": 1},
		EnumStatus: e.StatusNormal,
		PtrModel:   &m.Model{BuildinString: "text"},
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	err = x.Add(tx, m)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("add error: %s", err)
	}

	log.JSON(m)
}

func TestModelPersisterList(t *testing.T) {
	m := &m.Model{
		BuildinBool: true,
		EnumStatus:  e.StatusNormal,
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	mts, err := x.List(tx, m, []e.Status{e.StatusNormal, e.StatusDeleted}, 1, 9999)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	log.JSON(mts)
}
