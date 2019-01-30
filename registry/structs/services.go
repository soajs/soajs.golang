package structs

type Services map[string]Service

type Service struct {
	Group                 string `json:"group"`
	Port                  int    `json:"port"`
	RequestTimeout        int    `json:"requestTimeout"`
	RequestTimeoutRenewal int    `json:"requestTimeoutRenewal"`
	MaxPoolSize           int    `json:"maxPoolSize"`
	Version               int    `json:"version"`
	Authorization         bool   `json:"authorization"`
	ExtKeyRequired        bool   `json:"extKeyRequired"`
}
