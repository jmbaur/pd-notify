package notifications

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gen2brain/beeep"
)

func Notify(title, description, iconPath string) {
	var err error
	if _, insideTmux := os.LookupEnv("TMUX"); insideTmux {
		msg := fmt.Sprintf("%s\n%s\n%s", title, strings.Repeat("=", len(title)), description)
		cmd := exec.Command("tmux", "display-popup", "printf", msg)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	} else {
		err = beeep.Alert(title, description, iconPath)
	}

	if err != nil {
		fmt.Printf("%s\n%s\n%s", title, strings.Repeat("=", len(title)), description)
	}
}
