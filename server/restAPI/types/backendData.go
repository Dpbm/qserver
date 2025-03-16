package types

type BackendData struct {
	Name    string `json:"backend_name"`
	ID      string `json:"id"`
	Pointer uint64 `json:"pointer"`
	Plugin  string `json:"plugin"`
}
