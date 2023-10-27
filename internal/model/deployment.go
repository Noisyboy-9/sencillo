package model

type Deployment struct {
	id             int
	name           string
	coresRequired  float64
	memoryRequired float64
}

func NewDeployment(id int, name string, cores float64, memory float64) *Deployment {
	return &Deployment{
		id:             id,
		name:           name,
		coresRequired:  cores,
		memoryRequired: memory,
	}
}

func (dep *Deployment) Id() int {
	return dep.id
}

func (dep *Deployment) Name() string {
	return dep.name
}

func (dep *Deployment) CoresRequired() float64 {
	return dep.coresRequired
}

func (dep *Deployment) MemoryRequired() float64 {
	return dep.memoryRequired
}
