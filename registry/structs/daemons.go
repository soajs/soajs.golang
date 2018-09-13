package structs

type Daemons map[string]Daemon

type Daemon struct {
    Group           string                      `json:"group"`
    Port            int                         `json:"port"`
    Versions        map[string]DaemonVersion    `json:"versions"`
}

type DaemonVersion struct {
    Jobs            map[string]DaemonJob        `json:"jobs"`
}

type DaemonJob map[string]interface
