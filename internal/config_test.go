package internal

import (
	"flag"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	t.Log("TestWriteConfig")

	type args struct {
		context *cli.Context
	}

	tempFile := "/tmp/go-aws-sso/generated-config.yaml"
	defer func(file string) {
		t.Log("wat")
		if r := recover(); r != nil {
			t.Log("Recovered in f", r)
		}
		dir := path.Dir(file)
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}(tempFile)
	t.Log("TestWriteConfig", tempFile)

	flagSet := flag.NewFlagSet("path", flag.ContinueOnError)
	flagSet.String("path", tempFile, "")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Should create a default config file",
			args:    args{context: cli.NewContext(nil, flagSet, nil)},
			wantErr: false,
		},
	}
	t.Log("TestWriteConfig")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("TestWriteConfig")

			wantAppConfig := AppConfig{
				StartUrl: "https://my-login.awsapps.com/start#/",
				Region:   "eu-central-1",
			}

			t.Log("TestWriteConfig")
			got := writeConfig(tempFile, wantAppConfig)
			if got != nil {
				t.Errorf("Not expected: %q", got)
			}

			t.Log("TestWriteConfig")
			configFile, err := os.Open(tempFile)
			fail(err, t)

			bytes, err := ioutil.ReadFile(configFile.Name())
			fail(err, t)

			t.Log("TestWriteConfig")
			gotAppConfig := AppConfig{}
			err = yaml.Unmarshal(bytes, &gotAppConfig)
			fail(err, t)

			t.Log("TestWriteConfig")
			if !reflect.DeepEqual(gotAppConfig, wantAppConfig) {
				t.Errorf("got: %q, want: %q", gotAppConfig, wantAppConfig)
			}
		})
	}
}

func fail(err error, t *testing.T) {
	t.Log("fail", err)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
}
