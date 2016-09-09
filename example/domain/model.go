package domain

import (
	"time"

	"github.com/arstd/yan/example/enum"
)

// Model 模型示例
type Model struct {
	Id    int
	Name  string
	Flag  bool
	Score float32

	Map   map[string]interface{}
	Time  time.Time
	Slice []string

	Status      enum.Status
	Pointer     *Model
	StructSlice []*Model
}
