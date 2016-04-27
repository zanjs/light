package persist

import (
	"github.com/gotips/seed/gobatis"
)

// MyTestTable xxx
var MyTestTable myTestTable

type myTestTable struct{}

func (*myTestTable) Add(mt *gobatis.MyTest) (err error) {
	_ = `insert into my_test(test_name) values(${mt.TestName}) returning test_id`

	return nil
}

func (*myTestTable) Update(mt *gobatis.MyTest) (err error) {
	_ = `update my_test set test_name=${mt.TestName} where test_id=${mt.TestID}`

	return nil
}

func (*myTestTable) Get(testID int) (mt *gobatis.MyTest, err error) {
	_ = `select * from my_test where test_id = $testID`

	return nil, nil
}

func (*myTestTable) List(testID int) (mts []gobatis.MyTest, err error) {
	_ = `select * from my_test where test_id < $testID`

	return mts, nil
}

func (*myTestTable) Delete(testID int) (err error) {
	_ = `delete from my_test where test_id=$testID`

	return nil
}
