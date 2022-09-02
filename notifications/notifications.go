package notifications

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gen2brain/beeep"
)

func Notify(title, description, iconPath string) error {
	if _, insideTmux := os.LookupEnv("TMUX"); insideTmux {
		msg := fmt.Sprintf("%s\n%s", title, description)
		cmd := exec.Command("tmux", "display-popup", "printf", msg)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		err := beeep.Alert(title, description, iconPath)
		if err != nil {
			return err
		}
	}

	return nil
}
