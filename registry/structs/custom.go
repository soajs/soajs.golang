package structs

type CustomRegistries map[string]CustomRegistry

type CustomRegistry struct {
    ID             string          `json:"_id"`
	Name           string          `json:"name"`
	Locked         bool            `json:"locked"`
	Plugged        bool            `json:"plugged"`
	Shared         bool            `json:"shared"`
	Value          interface{}     `json:"value"`
	Created        string          `json:"created"`
	Author         string          `json:"author"`
}
