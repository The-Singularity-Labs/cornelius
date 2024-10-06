package sync

import (
	"time"

	"gopkg.in/yaml.v3"
)

type Duration time.Duration

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}

	tmp, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = Duration(tmp)
	return nil
}
