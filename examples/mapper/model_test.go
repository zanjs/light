package mapper

import (
	"testing"

	"github.com/wothing/log"

	m "github.com/arstd/gobatis/examples/domain"
	e "github.com/arstd/gobatis/examples/enums"
)

// var x ModelImplExample
var x ModelMapperImpl
var id int

func TestModelMapperInsert(t *testing.T) {
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
	err = x.Insert(tx, m)
	if err != nil {
		t.Fatalf("add error: %s", err)
	}

	CommitTx(tx)
	id = m.Id
	log.Infof("id=%d", m.Id)
}

func TestModelMapperUpdate(t *testing.T) {
	m := &m.Model{
		Id:          id,
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
		BuildinRune:    '个',
		BuildinString:  "update",
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
	a, err := x.Update(tx, m)
	if err != nil {
		t.Fatalf("add error: %s", err)
	}

	CommitTx(tx)
	log.Infof("affected=%d", a)
}

func TestModelMapperGet(t *testing.T) {
	tx, err := BeginTx()
	defer RollbackTx(tx)
	m, err := x.Get(tx, id)
	if err != nil {
		t.Fatalf("add error: %s", err)
	}

	CommitTx(tx)
	log.JSONIndent(m)
}

func TestModelMapperCount(t *testing.T) {
	m := &m.Model{
		BuildinBool: true,
		EnumStatus:  e.StatusNormal,
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	count, err := x.Count(tx, m, []e.Status{e.StatusNormal, e.StatusDeleted})
	if err != nil {
		t.Fatalf("count(%+v) error: %s", m, err)
	}

	CommitTx(tx)
	log.JSON(count)
}

func TestModelMapperSelect(t *testing.T) {
	m := &m.Model{
		BuildinBool: true,
		EnumStatus:  e.StatusNormal,
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	ms, err := x.List(tx, m, []e.Status{e.StatusNormal, e.StatusDeleted}, 0, 20)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	CommitTx(tx)
	log.JSONIndent(ms)
}

func TestModelMapperDelete(t *testing.T) {
	tx, err := BeginTx()
	defer RollbackTx(tx)
	a, err := x.Delete(tx, id)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("add error: %s", err)
	}

	log.JSON(a)
}
