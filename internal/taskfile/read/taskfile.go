package read

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-task/task/internal/libraries"

	"github.com/go-task/task/internal/taskfile"
	"github.com/mitchellh/go-homedir"

	"gopkg.in/yaml.v2"
)

var (
	// ErrIncludedTaskfilesCantHaveIncludes is returned when a included Taskfile contains includes
	ErrIncludedTaskfilesCantHaveIncludes = errors.New("task: Included Taskfiles can't have includes. Please, move the include to the main Taskfile")
)

// Taskfile reads a Taskfile for a given directory
func Taskfile(dir string, entrypoint string, u *libraries.Updater, up bool) (*taskfile.Taskfile, error) {
	path := filepath.Join(dir, entrypoint)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf(`task: No Taskfile found on "%s". Use "task --init" to create a new one`, path)
	}

	t, err := readTaskfile(path)
	if err != nil {
		return nil, err
	}

	if up {
		e := u.Update(t)
		if e != nil {
			log.Fatalln(e)
		}
	}

	for namespace, path := range t.Includes {

		if strings.Contains(path, ":") {

			home, err := homedir.Dir()
			if err != nil {
				return nil, err
			}

			libMap := strings.Split(path, ":")
			libPath := filepath.Join(home, ".task", "libraries", libMap[0])

			if _, err := os.Stat(libPath); os.IsNotExist(err) {
				e := u.Update(t)
				if e != nil {
					log.Fatalln(e)
				}
			}

			path = filepath.Join(libPath, libMap[1])

		} else {
			path = filepath.Join(dir, path)
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			path = filepath.Join(path, "Taskfile.yml")
		}
		includedTaskfile, err := readTaskfile(path)
		if err != nil {
			return nil, err
		}
		if len(includedTaskfile.Includes) > 0 {
			return nil, ErrIncludedTaskfilesCantHaveIncludes
		}
		if err = taskfile.Merge(t, includedTaskfile, namespace); err != nil {
			return nil, err
		}
	}

	path = filepath.Join(dir, fmt.Sprintf("Taskfile_%s.yml", runtime.GOOS))
	if _, err = os.Stat(path); err == nil {
		osTaskfile, err := readTaskfile(path)
		if err != nil {
			return nil, err
		}
		if err = taskfile.Merge(t, osTaskfile); err != nil {
			return nil, err
		}
	}

	for name, task := range t.Tasks {
		task.Task = name
	}

	return t, nil
}

func readTaskfile(file string) (*taskfile.Taskfile, error) {
	fmt.Sprintf("%s", file)

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var t taskfile.Taskfile
	return &t, yaml.NewDecoder(f).Decode(&t)
}
