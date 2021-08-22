package graph

import (
	common "github.com/lliryc/memego-common"
	"container/heap"
	"math"
	"math/rand"	
)

func ArgMin(array []Vertex) int {
	min := array[0].Id
	argMin := 0
	for i, value := range array {
		if min > value.Id {
			min = value.Id
			argMin = i
		}
	}
	return argMin
}

type Path struct {
	Vertices []Vertex
	FitVal   float32
	AdjMap   map[string]map[string]float32
}

func minFirst(vertices []Vertex) []Vertex {
	argMin := ArgMin(vertices)
	startSeq := vertices[argMin:]
	endSeq := vertices[:argMin]
	resSeq := append(startSeq, endSeq...)
	return resSeq
}

func filter(forFilter []Vertex, src []Vertex) []Vertex {
	res := []Vertex{}
	for _, w := range src {
		skip := false
		for _, v := range forFilter {
			if v == w {
				skip = true
				break
			}
		}
		if !skip {
			res = append(res, w)
		}
	}
	return res
}

func (path Path) CrossOver(other common.Instance) common.Instance {
	otherPath := other.(Path)
	n := len(path.Vertices)
	vertices1 := minFirst(path.Vertices)
	vertices2 := minFirst(otherPath.Vertices)
	crispPos := rand.Intn(n - 2)
	endPos := rand.Intn(n - crispPos - 1)
	seq := vertices1[crispPos:endPos]
	preSeq := vertices2[:crispPos]
	postSeq := vertices2[endPos:]
	restSeq := append(postSeq, preSeq...)
	restSeq = filter(seq, restSeq)
	seq = append(seq, restSeq...)
	return Path{Vertices: seq}
}

func (path Path) Mutation() common.Instance {
	const minLConst float64 = 0.05
	const maxLConst float64 = 0.3
	var n int = len(path.Vertices)
	minSeqLen := minLConst * float64(n)
	maxSeqLen := maxLConst * float64(n)
	seqLen := int(math.Max(minSeqLen+rand.Float64()*(maxSeqLen-minSeqLen), 1))
	crispPos := rand.Intn(n - seqLen - 1)
	preSeq := path.Vertices[0:crispPos]
	cutSeq := path.Vertices[crispPos : crispPos+seqLen]
	postSeq := path.Vertices[crispPos+seqLen:]
	tmpSeq := append(preSeq, postSeq...)
	tmpLen := len(tmpSeq)
	insertPos := rand.Intn(tmpLen)
	preSeq = tmpSeq[0:insertPos]
	postSeq = tmpSeq[insertPos:]
	vertices := append(preSeq, cutSeq...)
	vertices = append(vertices, postSeq...)
	return &Path{Vertices: vertices}
}

func (path Path) Reproduce() common.Instance {
	newPath := path
	return &newPath
}

type LocalSearch func (float32) ([]Vertex, float32)

func (path Path) Improve() common.Instance {	
	n := len(path.Vertices)
	vertices := path.Vertices
	lowerBound := path.eval(math.MaxFloat32)
	newVertices := make([]Vertex, n)
	copy(newVertices, vertices)
	newPath := Path{Vertices: newVertices, AdjMap: path.AdjMap, FitVal: path.FitVal}	
	var localSearch []LocalSearch = []LocalSearch{newPath.insert, newPath.directOpt2, newPath.opt2, newPath.swap2}	
	l := len(localSearch)
	i := 0
	for {		
		if len(localSearch) <= 0 {
			break
		}
		vertices, newLowerBound := localSearch[i](lowerBound)
		if newLowerBound >= lowerBound {
			localSearch = append(localSearch[:i], localSearch[i+1:]...)
			continue
		} else {
			lowerBound = newLowerBound
		}
		newPath = Path{
			Vertices: vertices,
			FitVal: lowerBound,
			AdjMap: path.AdjMap,
		}
		i = (i + 1) % l
	}
	return &newPath
}

func (path Path) swap2(lowerBound float32) ([]Vertex, float32) {
	vertices := path.Vertices
	n := len(vertices)
	// TODO: implement index permutation instead of vertices clone
	candidate := make([]Vertex, len(vertices))
	copy(candidate, vertices)
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			if i == 0 && j == n-1 {
				continue
			}
			vertices[i], vertices[j] = vertices[j], vertices[i]
			res := path.eval(lowerBound)
			if res < lowerBound {
				lowerBound = res
				copy(candidate, vertices)
			}
			vertices[i], vertices[j] = vertices[j], vertices[i]
		}
	}
	return candidate, lowerBound
}

func (path Path) insert(lowerBound float32) ([]Vertex, float32) {	
	vertices := path.Vertices
	adjMap := path.AdjMap
	n := len(vertices)
	candidate := make([]Vertex, len(vertices))
	vertices2 := make([]Vertex, len(vertices))
	for i := 0; i < n; i++ {
		for j := i + 2; j < n; j++ {
			if i == 0 && j == n-1 {
				continue
			}
			copy(vertices2, vertices)
			v := vertices2[i]
			copy(vertices2[i:], vertices2[i+1:])
			copy(vertices2[j+1:], vertices2[j:])
			vertices2[j] = v
			newPath := Path{Vertices: vertices2, AdjMap: adjMap}
			res := newPath.eval(lowerBound)
			if res < lowerBound {
				copy(candidate, vertices2)
				lowerBound = res
			}
		}
	}
	return candidate, lowerBound
}

func reverse(vertices []Vertex, start int, end int) []Vertex {
	for i, j := start, end; i < j; i, j = i+1, j-1 {
		vertices[i], vertices[j] = vertices[j], vertices[i]
	}
	return vertices
}

func (path Path) opt2(lowerBound float32) ([]Vertex, float32) {
	vertices := path.Vertices
	adjMap := path.AdjMap
	n := len(vertices)	
	vertices2 := make([]Vertex, len(vertices))
	copy(vertices2, vertices)
	for i := 0; i < n-1; i++ {
		vi1 := vertices2[i].Id
		vi2 := vertices2[i+1].Id
		vali, _ := adjMap[vi1][vi2]
		for j := i + 2; j < n-1; j++ {
			vj1 := vertices2[j].Id
			vj2 := vertices2[j+1].Id
			valj, _ := adjMap[vj1][vj2]
			valij1, ok := adjMap[vi1][vj1]
			if !ok {
				continue
			}
			valij2, ok := adjMap[vi2][vj2]
			if !ok {
				continue
			}
			if vali+valj > valij1+valij2 {
				reverse(vertices2[i+1:j+1], i+1, j)
				newPath := Path{Vertices: vertices2, AdjMap: adjMap}
				lowerBound = newPath.eval(lowerBound)
				break
			}
		}
	}
	return vertices2, lowerBound
}

func (path Path) directOpt2(lowerBound float32) ([]Vertex, float32) {
	vertices := path.Vertices
	adjMap := path.AdjMap
	const longestEdgeNumber = 2
	var n int = len(vertices)
	var edges EdgeHeap = make([]Edge, 0)
	heap.Init(&edges)
	heap.Push(&edges, Edge{V: n - 1, W: 0, Value: adjMap[vertices[n-1].Id][vertices[0].Id]})
	for i := 0; i < n-1; i++ {
		v := vertices[i].Id
		w := vertices[i+1].Id
		heap.Push(&edges, Edge{V: i, W: i + 1, Value: adjMap[v][w]})
	}
	vertices2 := make([]Vertex, len(vertices))
	for i := 0; i < longestEdgeNumber; i++ {
		longEdge := heap.Pop(&edges).(Edge)
		i1 := longEdge.V
		i2 := longEdge.W
		vi1 := vertices2[i1].Id
		vi2 := vertices2[i2].Id
		vali := longEdge.Value
		copy(vertices2, vertices)
		for j := 0; j < n-1; j++ {
			j1 := j
			j2 := j + 1
			vj1 := vertices2[j1].Id
			vj2 := vertices2[j2].Id
			if i1 == j1 || i1 == j2 || i2 == j1 {
				continue
			}
			valj, _ := adjMap[vj1][vj2]
			valij1, ok := adjMap[vi1][vj1]
			if !ok {
				continue
			}
			valij2, ok := adjMap[vi2][vj2]
			if !ok {
				continue
			}
			if vali+valj > valij1+valij2 {
				reverse(vertices2[i+1:j+1], i+1, j)
				newPath := Path{Vertices: vertices2, AdjMap: adjMap}
				lowerBound = newPath.eval(lowerBound)
				break
			}
		}
	}
	return vertices2, lowerBound
}

func (path Path) eval(lowerBound float32) float32 {
	vertices := path.Vertices
	adjMap := path.AdjMap
	var n int = len(vertices)
	var sum float32 = adjMap[vertices[0].Id][vertices[n-1].Id]
	for i := 0; i < len(vertices)-1; i++ {
		v := vertices[i].Id
		wDict, ok := adjMap[v]
		if !ok {
			return math.MaxFloat32
		}
		for j := i + 1; j < len(vertices); j++ {
			w := vertices[j].Id
			val, ok := wDict[w]
			if !ok {
				return math.MaxFloat32
			}
			sum += val
			if sum > lowerBound {
				return math.MaxFloat32
			}
		}
	}

	return sum
}

func (path Path) ComputeFitness() float32 {
	path.FitVal = path.eval(math.MaxFloat32)
	return path.FitVal
}

func (path Path) Fitness() float32 {
	return path.FitVal
}

func (path Path) Less(other common.Instance) bool{
	otherPath := other.(Path)
	return path.FitVal > otherPath.FitVal
}
