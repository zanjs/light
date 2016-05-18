package domain

import "github.com/arstd/gobatis/examples/enums"

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
