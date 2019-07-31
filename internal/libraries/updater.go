package libraries

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-task/task/internal/execext"

	"github.com/go-task/task/internal/logger"
	"github.com/go-task/task/internal/taskfile"
	"github.com/mitchellh/go-homedir"
)

// Updater clones updates libraries in home directory
type Updater struct {
	Logger logger.Logger
}

func (u *Updater) Update(t *taskfile.Taskfile) error {
	_, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	libDir := filepath.Join(home, ".task", "libraries")
	if _, err := os.Stat(libDir); os.IsNotExist(err) {
		err := os.MkdirAll(libDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var stdout bytes.Buffer
	for namespace, url := range t.Libraries {
		libPath := filepath.Join(libDir, namespace)

		u.Logger.VerboseOutf("Checking library %s at %s", namespace, libPath)

		_, err := os.Stat(libPath)
		if os.IsNotExist(err) {
			// Clone the repo
			u.Logger.Outf("Cloning library %s from %s", namespace, url)
			opts := &execext.RunCommandOptions{
				Command: strings.Join([]string{"git clone", url, namespace}, " "),
				Dir:     libDir,
				Stdout:  &stdout,
				Stderr:  u.Logger.Stderr,
			}
			execext.RunCommand(context.Background(), opts)
		} else {
			// Check for updates and update the repos if they are outdated
			u.Logger.Outf("Updating library %s from %s", namespace, url)
			opts := &execext.RunCommandOptions{
				Command: "git pull",
				Dir:     filepath.Join(libDir, namespace),
				Stdout:  &stdout,
				Stderr:  u.Logger.Stderr,
			}
			execext.RunCommand(context.Background(), opts)
		}

	}

	return nil
}
