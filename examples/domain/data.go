package domain

// Demo 示例结构体
type Demo struct {
	Id         int
	Name       string
	ThirdField bool
	Status     Status
	Content    *Demo
	Map        map[string]Status
}

type Status int8

const (
	StatusNormal Status = 1
)

/*
create table demos (
	id serial primary key,
	name text,
	third_field bool,
	status smallint,
	Content text
)
*/
