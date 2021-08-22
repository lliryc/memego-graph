package graph

type TspPolicy struct {
	VerticesNum int
	g int
	r int
}

func r(g int, v int) int{
	return int(0.2*float32(g) + 0.05 * float32(v) + 10)
}

func (policy TspPolicy) R() int {
	if policy.r == 0 {
		policy.r = r(policy.g, policy.VerticesNum)
	}
	return policy.r 
}

func (policy TspPolicy) GetPopulationSize() int {
	return 11 * policy.r
}

func (policy TspPolicy) GetSolutionN() int {
	return policy.VerticesNum
}

func (policy TspPolicy) GetCrossoverN() int {
	return 8 * policy.r
}

func (policy TspPolicy) GetReproductionN() int {
	return policy.r
}

func (policy TspPolicy) GetMutationN() int {
	return 2 * policy.r
}

func (policy TspPolicy) SetGeneration(g int) {
	policy.g = g
	policy.r = r(policy.g, policy.VerticesNum)
}

