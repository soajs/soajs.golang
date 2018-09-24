package structs

type Databases map[string]Database

type Database struct {
    Name                    string                  `json:"name"`
    Prefix                  string                  `json:"prefix"`
    Cluster                 string                  `json:"cluster"`
    Servers                 ServersList             `json:"servers"`
    Credentials             ServerCredentials       `json:"credentials"`
    Streaming               map[string]interface{}  `json:"streaming"`
    RegistryLocation        RegistryLocation        `json:"registryLocation"`
    URLParam                interface{}             `json:"URLParam"`
    ExtraParam              interface{}             `json:"extraParam"`

    //NOTE session specific entries
    Store                   interface{}             `json:"store,omitempty"`
    Collection              string                  `json:"collection,omitempty"`
    Stringify               bool                    `json:"stringify,omitempty"`
    ExpireAfter             int                     `json:"expireAfter,omitempty"`
}

type ServersList struct {
    Host                    string                  `json:"host"`
    Port                    int                     `json:"port"`
}

type ServerCredentials struct {
    Username                string                  `json:"username"`
    Password                string                  `json:"password"`
}

type RegistryLocation struct {
    L1                      string                  `json:"l1"`
    L2                      string                  `json:"l2"`
    Env                     string                  `json:"env"`
}
