package components

type CollisionType string

const (
	Solid   CollisionType = "solid"
	Trigger CollisionType = "trigger"
)

type Collision struct {
	Width  int
	Height int
	Type   CollisionType
}
