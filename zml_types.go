package zml

const (
	// RECT sets elemenet type to rectangle
	RECT int = 0
	// DESCISION sets elemenet type to decision node
	DESCISION = 1
	// CIRCLE sets elemenet type to circle
	CIRCLE = 2
)

type elemenet struct {
	Name string
	Type int
}

// Font font settings
type Font struct {
	Name  string
	Size  float64
	Color Color
}

type edge struct {
	from        elemenet
	to          elemenet
	directional bool
	Label       string
}

func (e *edge) From() elemenet {
	return e.from
}

func (e *edge) To() elemenet {
	return e.to
}

type elemenetCoord struct {
	X float64
	Y float64
}
