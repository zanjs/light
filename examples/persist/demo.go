package persist

import "github.com/arstd/persist/examples/domain"

//go:generate go run ../../main.go

// DemoPersister xxx
type DemoPersister interface {

	// insert into demos(demo_name) values(${d.DemoName}) returning id
	Add(d *domain.Demo) (err error)

	// update demos
	// set demo_name=${d.DemoName}
	// where id=${d.Id}
	Update(d *domain.Demo) (err error)

	// select id, demo_name, demo_status
	// from demos where id = ${id}
	Get(id string) (d *domain.Demo, err error)

	// select id, demo_name, demo_status
	// from demos
	// where id < $id[ and demo_name!=${d.DemoName}][ and demo_status=${d.DemoStatus}]
	List(d *domain.Demo) (ds []*domain.Demo, err error)

	// delete from demos where id=${id}
	Delete(id string) (err error)
}
