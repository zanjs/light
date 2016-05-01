package domain

// Demo xxx
type Demo struct {
	Id         string
	DemoName   string
	DemoStatus string
	DemoStruct *Demo
}

//
// create table demos (
// id uuid primary key default gen_random_uuid(),
// demo_name text,
// demo_status text,
// demo_struct jsonb
// );
