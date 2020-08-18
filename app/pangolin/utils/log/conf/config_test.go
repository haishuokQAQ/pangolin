package conf

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/AlekSi/pointer"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

func TestUnmarshalYamlConfig(t *testing.T) {

	text := `core: zap
level: debug
formatter: console
outputs:
  - type: stdout
  - type: file
    file: /tmp/app_name.log
  - type: rotate_file
    rotate_file:
      file_name: /tmp/app_name_rotated.log
      max_size: 100
      max_age: 10
      max_backups: 3
      localtime: true
      compress: false
`

	actual := &Config{}

	err := yaml.Unmarshal([]byte(text), actual)
	if err != nil {
		t.Fatal("fail to unmarshal text")
	}

	expected := &Config{
		Core:      ZapCore,
		Level:     LevelDebug,
		Formatter: ConsoleFormater,
		Outputs: []Output{
			{
				Type: "stdout",
			},
			{
				Type: "file",
				File: pointer.ToString("/tmp/app_name.log"),
			},
			{
				Type: "rotate_file",
				RotateFile: &RotateFile{
					FileName:   "/tmp/app_name_rotated.log",
					MaxSize:    100,
					MaxAge:     10,
					MaxBackups: 3,
					LocalTime:  true,
					Compress:   false,
				},
			},
		},
	}

	spew.Printf("==actual==\n%v\n==expected==\n%v\n", actual, expected)

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
		return
	}
}

func TestUnmarshalJSONConfig(t *testing.T) {

	text := `{
    "core": "zap",
    "level": "debug",
    "formatter": "console",
    "outputs": [
        {
            "type": "stdout"
        },
        {
            "type": "file",
            "file": "/tmp/app_name.log"
        },
        {
            "type": "rotate_file",
            "rotate_file": {
                "file_name": "/tmp/app_name_rotated.log",
                "max_size": 100,
                "max_age": 10,
                "max_backups": 3,
                "localtime": true,
                "compress": false
            }
        }
    ]
}`

	actual := &Config{}

	err := json.Unmarshal([]byte(text), actual)
	if err != nil {
		t.Fatalf("fail to unmarshal text: %v", err)
	}

	expected := &Config{
		Core:      ZapCore,
		Level:     LevelDebug,
		Formatter: ConsoleFormater,
		Outputs: []Output{
			{
				Type: "stdout",
			},
			{
				Type: "file",
				File: pointer.ToString("/tmp/app_name.log"),
			},
			{
				Type: "rotate_file",
				RotateFile: &RotateFile{
					FileName:   "/tmp/app_name_rotated.log",
					MaxSize:    100,
					MaxAge:     10,
					MaxBackups: 3,
					LocalTime:  true,
					Compress:   false,
				},
			},
		},
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
		return
	}
}
