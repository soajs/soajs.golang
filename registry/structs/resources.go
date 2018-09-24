package structs

type Resources map[string]Resource

type Resource struct {
  Id            string          `json:"_id"`
  Name          string          `json:"name"`
  Type          string          `json:"type"`
  Category      string          `json:"category"`
  Created       string          `json:"created"`
  Author        string          `json:"author"`
  Locked        bool            `json:"locked"`
  Plugged       bool            `json:"plugged"`
  Shared        bool            `json:"shared"`
  Config        interface{}     `json:"config"`
}
