package app

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"strings"
)

func formatContainer(c types.Container) string {
	return fmt.Sprintf(
		"%v %v %v %v",
		c.ID,
		c.Image,
		strings.Join(c.Names, ":"),
		c.Status,
	)
}
