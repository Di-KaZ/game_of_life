package grid

type Cell struct {
	Alive bool `json:"alive"`
	Turns int  `json:"turns"`
}
