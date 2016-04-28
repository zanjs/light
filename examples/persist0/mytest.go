package persist0

import (
	// import sql driver
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wothing/log"
)

// MyTestTable xxx
var MyTestTable myTestTable

type myTestTable struct{}

// Add xxx
func (*myTestTable) Add(mt *MyTest) (err error) {
	q := `insert into my_test(test_name) values($1) returning test_id`
	err = db.QueryRow(q, mt.TestName).Scan(&mt.TestID)
	if err != nil {
		log.Errorf("insert(%s, %s) error: %s", q, mt.TestName, err)
		return err
	}

	return nil
}

// Update xxx
func (*myTestTable) Update(mt *MyTest) (err error) {
	q := `update my_test set test_name=$1 where test_id=$2`
	res, err := db.Exec(q, mt.TestName, mt.TestID)
	if err != nil {
		log.Errorf("update(%s, %s, %d) error: %s", q, mt.TestName, mt.TestID, err)
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		log.Errorf("update(%s, %s, %d) error: %s", q, mt.TestName, mt.TestID, err)
		return err
	} else if a != 1 {
		log.Errorf("update(%s, %s, %d) expected affected 1 row, but actual affected %d rows",
			q, mt.TestName, mt.TestID, a)
		return fmt.Errorf("expected affected 1 row, but actual affected %d rows", a)
	}

	return nil
}

// List xxx
func (*myTestTable) List(testID int) (mts []MyTest, err error) {
	q := `select * from my_test where test_id < $1`
	rows, err := db.Query(q, testID)
	if err != nil {
		log.Errorf("query(%s, %d) error: %s", q, testID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MyTest
		err = rows.Scan(&ms.TestID, &ms.TestName)
		if err != nil {
			log.Errorf("scan rows for query(%s, %d) error: %s", q, testID, err)
			return nil, err
		}
		mts = append(mts, ms)
	}
	if err = rows.Err(); err != nil {
		log.Errorf("scan rows for query(%s, %d) last error: %s", q, testID, err)
		return nil, err
	}

	return mts, nil
}

// Delete xxx
func (*myTestTable) Delete(testID int) (err error) {
	q := `delete from my_test where test_id=$1`
	res, err := db.Exec(q, testID)
	if err != nil {
		log.Errorf("delete(%s, %d) error: %s", q, testID, err)
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		log.Errorf("delete(%s, %d) error: %s", q, testID, err)
		return err
	} else if a != 1 {
		log.Errorf("delete(%s, %d) expected affected 1 row, but actual affected %d rows",
			q, testID, a)
		return fmt.Errorf("expected affected 1 row, but actual affected %d rows", a)
	}

	return nil
}
