package domain

import "github.com/arstd/gobatis/examples/enums"

import fmt "fmt"
import "net/http"

var _ = fmt.Print
var _ = http.StatusOK

// Demo 示例结构体
type Demo struct {
	Id         int
	Name       string
	ThirdField bool
	Status     enums.Status
	Content    *Demo
	Map        map[string]string
}

/*
create table demos (
	id serial primary key,
	name text,
	third_field bool,
	status smallint,
	Content text
)
*/
