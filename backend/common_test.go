package backend

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"testing"
)

func TestProgress(t *testing.T) {
	fmt.Printf("%s\n", GetMessage(88, &progress.UnitsBytes))
}