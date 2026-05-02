package domain

// Container describes a running Docker container available for log streaming.
type Container struct {
	ID    string
	Name  string
	Image string
}

func (c Container) DisplayName() string {
	if c.Name != "" {
		return c.Name
	}
	if c.ID != "" {
		return c.ID
	}
	return "<unknown>"
}
