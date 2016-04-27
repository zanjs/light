package persist

import (
	"github.com/arstd/persist/examples/domain"

	// import sql driver
	_ "github.com/lib/pq"
	"github.com/wothing/log"
)

type DemoPersisterData struct{}

func (*DemoPersisterData) Add(d *domain.Demo) (err error) {
	q := "insert into demos(demo_name) values($1) returning id"

	err = db.QueryRow(q, d.DemoName).Scan(&d.Id)
	if err != nil {
		log.Errorf("insert(%s, %s) error: %s", q, d.DemoName, err)
		return err
	}
	return nil
}

func (*DemoPersisterData) Update(d *domain.Demo) (err error) {
	q := "update demos set demo_name=$1 where id=$2"

	res, err := db.Exec(q, d.DemoName, d.Id)
	if err != nil {
		log.Errorf("update(%s, %s, %s) error: %s", q, d.DemoName, d.Id, err)
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		log.Errorf("update(%s, %s, %s) error: %s", q, d.DemoName, d.Id, err)
		return err
	} else if a != 1 {
		log.Errorf("update(%s, %s, %s) expected affected 1 row, but actual affected %d rows",
			q, d.DemoName, d.Id, a)
		return err
	}
	return nil
}

func (*DemoPersisterData) Get(id string) (d *domain.Demo, err error) {
	q := "select id, demo_name, demo_status from demos where id = $1"

	x := domain.Demo{}
	err = db.QueryRow(q, id).
		Scan(&x.Id, &x.DemoName, &x.DemoStatus)
	if err != nil {
		log.Errorf("query(%s, %s) error: %s", q, id, err)
		return nil, err
	}
	return &x, nil
}

func (*DemoPersisterData) List(d *domain.Demo) (ds []*domain.Demo, err error) {
	q := "select id, demo_name, demo_status from demos where id < $id[ and demo_name!=$1][ and demo_status=$2]"

	rows, err := db.Query(q, d.DemoName, d.DemoStatus)
	if err != nil {
		log.Errorf("query(%s, %s, %s) error: %s", q, d.DemoName, d.DemoStatus, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var x domain.Demo
		err = rows.Scan(&x.Id, &x.DemoName, &x.DemoStatus)
		if err != nil {
			log.Errorf("scan rows for query(%s, %s, %s) error: %s", q, d.DemoName, d.DemoStatus, err)
			return nil, err
		}
		ds = append(ds, &x)
	}
	if err = rows.Err(); err != nil {
		log.Errorf("scan rows for query(%s, %s, %s) last error: %s", q, d.DemoName, d.DemoStatus, err)
		return nil, err
	}
	return ds, nil
}

func (*DemoPersisterData) Delete(id string) (err error) {
	q := "delete from demos where id=$1"

	res, err := db.Exec(q, id)
	if err != nil {
		log.Errorf("delete(%s, %s) error: %s", q, id, err)
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		log.Errorf("delete(%s, %s) error: %s", q, id, err)
		return err
	} else if a != 1 {
		log.Errorf("delete(%s, %s) expected affected 1 row, but actual affected %d rows",
			q, id, a)
		return err
	}
	return nil
}
