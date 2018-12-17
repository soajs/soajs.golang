package structs

type Services map[string]Service

type Service struct {
	Group                 string `json:"group"`
	Port                  int    `json:"port"`
	RequestTimeout        int    `json:"requestTimeout"`
	RequestTimeoutRenewal int    `json:"requestTimeoutRenewal"`
	MaxPoolSize           int    `json:"maxPoolSize"`
	Authorization         bool   `json:"authorization"`
	Version               int    `json:"version"`
	ExtKeyRequired        bool   `json:"extKeyRequired"`
}
