package graph

import (
	"math/rand"
	"sort"
	common "github.com/lliryc/memego-common"
)

type GraphGenerator struct{
	AdjMap map[string]map[string]float32
	Vertices[] Vertex
}

func nkeys(m map[string]float32) []string{
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func filterAdj(passed map[string]bool, vertices[]string) []string {
	filtered := make([]string, 0)
	for _,v := range vertices {
		_, ok := passed[v]
		if ok {
			continue
		}
		filtered = append(filtered, v)
	}
	return filtered
}

func (gen GraphGenerator) Create(policy common.Policy) common.Generation {
	var vertices[]Vertex = gen.Vertices
	
	vdict := make(map[string]Vertex)	
	for _,v := range vertices {
		vdict[v.Id] = v 
	}
	n := len(vertices)
	adjMap := gen.AdjMap
	size := policy.GetPopulationSize()
	var generation common.Generation = make(common.Generation, size)
	neighbors := make(map[string][]string, n)
	for i := 0; i < n; i++ {
		v := vertices[i]
		neighbors[v.Id] = nkeys(adjMap[v.Id])
	}
	// generate population
	for i := 0; i < size; i++ {
		newVertices := make([]Vertex, n)
		vix := rand.Intn(n)
		v := vertices[vix]		
		newVertices[0] = v
		adjacents := neighbors[v.Id]
		passed := make(map[string]bool)
		passed[v.Id] = true
		adjacents = filterAdj(passed, adjacents)
		for j := 1; j < n; j++ {			
			nn := len(adjacents)
			vId := adjacents[rand.Intn(nn)]
			newVertices[j] = vdict[vId]
			passed[vId] = true
			adjacents = neighbors[vId]
			adjacents = filterAdj(passed, adjacents)
		}
		var instance common.Instance = &Path{Vertices: newVertices, AdjMap: adjMap}
		generation[i] = instance
	}
	sort.Sort(generation)
	return generation		
}