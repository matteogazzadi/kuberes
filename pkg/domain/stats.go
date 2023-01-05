package domain

type PodStats struct {
	Namespace string
	Cpu       Resource
	Memory    Resource
}

type Resource struct {
	Request int64
	Limit   int64
}
