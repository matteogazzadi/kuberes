package domain

type K8sStats struct {
	PodName   string
	Namespace string
	Cpu       Resource
	Memory    Resource
}

type Resource struct {
	Request int64
	Limit   int64
}
