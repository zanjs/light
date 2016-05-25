package enums

type Status int8

const (
	StatusNormal Status = 1
)

func (s Status) String() string {
	return "Normal"
}
