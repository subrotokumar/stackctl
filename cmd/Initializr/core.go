package initializr

import (
	"fmt"

	"github.com/subrotokumar/stackctl/v4/cmd/core"
	"github.com/subrotokumar/stackctl/v4/cmd/ui/inputtext"
	"github.com/subrotokumar/stackctl/v4/cmd/ui/selector"
)

type Step func() (back bool)

func RunSteps(steps []Step) {
	for i := 0; i < len(steps); {
		if steps[i]() {
			if i > 0 {
				i--
			}
		} else {
			i++
		}
	}
}

func AskSelect(title string, options []string, dst *string) bool {
	v, back := selector.New(title, options).WithDefault(*dst).RunWithBack()
	if back {
		return true
	}
	*dst = v
	return false
}

func AskText(title string, dst *string) bool {
	v, back := inputtext.New(title, *dst).RunWithBack()
	if back {
		return true
	}
	*dst = v
	return false
}

func PrintKV(pairs [][2]string) {
	for _, kv := range pairs {
		fmt.Printf("  %s: %s\n", core.LogoStyle.Render(kv[0]), kv[1])
	}
}
