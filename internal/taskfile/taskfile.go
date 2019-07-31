package taskfile

// Taskfile represents a Taskfile.yml
type Taskfile struct {
	Version    string
	Expansions int
	Output     string
	Libraries  map[string]string
	Includes   map[string]string
	Vars       Vars
	Env        Vars
	Tasks      Tasks
}

// UnmarshalYAML implements yaml.Unmarshaler interface
func (tf *Taskfile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&tf.Tasks); err == nil {
		tf.Version = "1"
		return nil
	}

	var taskfile struct {
		Version    string
		Expansions int
		Output     string
		Libraries  map[string]string
		Includes   map[string]string
		Vars       Vars
		Env        Vars
		Tasks      Tasks
	}
	if err := unmarshal(&taskfile); err != nil {
		return err
	}
	tf.Version = taskfile.Version
	tf.Expansions = taskfile.Expansions
	tf.Output = taskfile.Output
	tf.Libraries = taskfile.Libraries
	tf.Includes = taskfile.Includes
	tf.Vars = taskfile.Vars
	tf.Env = taskfile.Env
	tf.Tasks = taskfile.Tasks
	if tf.Expansions <= 0 {
		tf.Expansions = 2
	}
	return nil
}
