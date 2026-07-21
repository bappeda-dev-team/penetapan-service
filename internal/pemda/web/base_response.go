package web

type Response[T any] struct {
	Data T `json:"data"`
}
