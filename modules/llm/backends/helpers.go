package backends

import (
	"fmt"

	"github.com/YourTechBud/inferix/modules/llm/types"
)

func (b *Backends) getBackend(name string) (types.Backend, error) {
	backend, ok := b.backends[name]
	if !ok {
		return nil, fmt.Errorf("backend %s not found", name)
	}
	return backend, nil
}
