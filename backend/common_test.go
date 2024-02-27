package backend

import (
	"fmt"
	"testing"

	"github.com/jedib0t/go-pretty/v6/progress"
)

func TestProgress(t *testing.T) {
	fmt.Printf("%s\n", GetMessage(88, &progress.UnitsBytes))
	fmt.Printf("%s\n", ProgressBar(54, 100, "progress"))
	fmt.Printf("%s\n", ScaleBar(0, 10, 100))
}
