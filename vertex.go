package graph

type Vertex struct {
	Id string
}

func (p *Vertex) String() string { return p.Id }

func (p *Vertex) Less(other *Vertex) bool {
	return p.Id < other.Id
}
