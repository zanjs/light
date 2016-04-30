package persist

import m "github.com/arstd/persist/examples/domain"

//go:generate go run ../../main.go

// DemoPersister xxx
type DemoPersister interface {

	// insert into demos(demo_name) values(${d.DemoName}) returning id
	Add(d *m.Demo) (err error)

	// update demos
	// set demo_name=${d.DemoName}
	// where id=${d.Id}
	Update(d *m.Demo) (err error)

	// select id, demo_name, demo_status
	// from demos where id=${id}
	Get(id string) (d *m.Demo, err error)

	// select id, demo_name, demo_status
	// from demos
	// where id<${d.Id} and demo_name!=${d.DemoName} and demo_status='1'
	List(d *m.Demo) (ds []*m.Demo, err error)

	// delete from demos where id=${id}
	Delete(id string) (err error)
}
