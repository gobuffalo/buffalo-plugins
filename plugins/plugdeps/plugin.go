package plugdeps

import (
	"encoding/json"
)

// Plugin represents a Go plugin for Buffalo applications
type Plugin struct {
	Binary   string    `toml:"binary" json:"binary"`
	GoGet    string    `toml:"go_get,omitempty" json:"go_get,omitempty"`
	Local    string    `toml:"local,omitempty" json:"local,omitempty"`
	Commands []Command `toml:"command,omitempty" json:"commands,omitempty"`
}

// String implementation of fmt.Stringer
func (p Plugin) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}
