package persist

import (
	"database/sql"

	"github.com/arstd/gobatis/examples/domain"
	"github.com/arstd/gobatis/examples/enums"
)

//go:generate gobatis

// DemoPersister 示例接口
type DemoPersister interface {

	// insert into demos(name, third_field, status, content)
	// values(${d.Name}, ${d.ThirdField}, ${d.Status}, ${d.Content})
	// returning id
	Add(d *domain.Demo) error

	// update demos
	// set name=${d.Name}, third_field=${d.ThirdField},
	//   status=${d.Status}, content=${d.Content}
	// where id=${d.Id}
	Modify(d *domain.Demo) (int64, error)

	// delete from demos where id=${id}
	Remove(id int) (int64, error)

	// select id, name, third_field, status, content
	// from demos where id=${id}
	Get(id int) (*domain.Demo, error)

	// select count(id)
	// from demos
	// where name=${d.Name}
	//   [?{d.ThirdField != false} and third_field=${d.ThirdField} ]
	//   [?{d.Content != nil} and content=${d.Content} ]
	//   [?{len(statuses) != 0} and status in (${statuses}) ]
	Count(tx *sql.Tx, d *domain.Demo, statuses []enums.Status) (int64, error)

	// select id, name, third_field, status, content
	// from demos
	// where name=${d.Name}
	//   [?{d.ThirdField != false} and third_field=${d.ThirdField} ]
	//   [?{d.Content != nil} and content=${d.Content} ]
	//   [?{len(statuses) != 0} and status in (${statuses}) ]
	List(tx *sql.Tx, d *domain.Demo, statuses []enums.Status, page, size int) ([]*domain.Demo, error)
}
