package registry

type Service struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    Address string `json:"address"`
    Port    int    `json:"port"`
}