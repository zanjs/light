package persist

import (
	"testing"

	"github.com/gotips/log"

	m "github.com/arstd/gobatis/examples/domain"
)

var x DemoPersist
var id int

func TestDemoPersisterAdd(t *testing.T) {
	d := &m.Demo{
		Name:       "Name",
		ThirdField: true,
		Status:     m.StatusNormal,
		Content:    &m.Demo{},
	}

	err := x.Add(d)
	if err != nil {
		t.Fatalf("insert %#v error: %s", d, err)
	}

	id = d.Id
}

func TestDemoPersisterUpdate(t *testing.T) {
	d := &m.Demo{
		Id:         id,
		Name:       "Name_updated",
		ThirdField: true,
		Status:     m.StatusNormal,
		Content:    &m.Demo{},
	}

	err := x.Modify(d)
	if err != nil {
		t.Fatalf("update %#v error: %s", d, err)
	}
}

func TestDemoGet(t *testing.T) {
	d, err := x.Get(id)

	if err != nil {
		t.Fatalf("get(%+v) error: %s", d, err)
	}

	log.JSON(d)
}

func TestDemoPersisterList(t *testing.T) {
	d := &m.Demo{
		Id:         1,
		Name:       "Name",
		ThirdField: true,
		Status:     m.StatusNormal,
		Content:    &m.Demo{},
	}
	log.JSON(d.Content)
	tx, err := BeginTx()
	defer RollbackTx(tx)
	mts, err := x.List(tx, d, []m.Status{m.StatusNormal}, 1, 9999)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("list(%+v) error: %s", d, err)
	}

	log.JSON(mts)
}

func TestDemoPersisterCount(t *testing.T) {
	d := &m.Demo{
		Id:         1,
		Name:       "Name",
		ThirdField: true,
		Status:     m.StatusNormal,
		Content:    &m.Demo{},
	}
	log.JSON(d.Content)
	tx, err := BeginTx()
	defer RollbackTx(tx)
	count, err := x.Count(tx, d, []m.Status{m.StatusNormal})
	CommitTx(tx)

	if err != nil {
		t.Fatalf("list(%+v) error: %s", d, err)
	}

	log.Debug(count)
}

func TestDemoPersisterRemove(t *testing.T) {
	err := x.Remove(id)
	if err != nil {
		t.Fatalf("delete by Id=%d error: %s", id, err)
	}
}
