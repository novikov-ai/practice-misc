package groups

type Database struct {
	InMemory bool   `toml:"in_memory"`
	Driver   string `toml:"driver"`
	Source   string `toml:"source"`
}
