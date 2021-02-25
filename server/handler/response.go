package handler

type privateInfo struct {
	ID string `json:"id"`
}

type publicInfo struct {
	ID          string
	IsOwner     bool
	Active      bool
	Title       string
	Description string
	Result      string
	Count       int
}
