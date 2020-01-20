package mapstructure

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type Test struct {
	ID            string    `json:"ID" yaml:"ID"`
	Name          string    `json:"Name" yaml:"Name"`
	LogoImageURL  string    `json:"LogoImageURL" yaml:"LogoImageURL"`
	MapSearchWord string    `json:"MapSearchWord" yaml:"MapSearchWord"`
	CreatedAt     time.Time `json:"CreatedAt" yaml:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt" yaml:"UpdatedAt"`
}

func Test_CreateStructBySkpFields(t *testing.T) {
	type args struct {
		st         interface{}
		skipFields []string
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "creating anonymous struct test without CreatedAt and UpdatedAt",
			args: args{
				st:         Test{},
				skipFields: []string{"CreatedAt", "UpdatedAt"},
			},
			want: struct {
				ID            string `json:"ID" yaml:"ID"`
				Name          string `json:"Name" yaml:"Name"`
				LogoImageURL  string `json:"LogoImageURL" yaml:"LogoImageURL"`
				MapSearchWord string `json:"MapSearchWord" yaml:"MapSearchWord"`
			}{
				ID:            "",
				Name:          "",
				LogoImageURL:  "",
				MapSearchWord: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStructBySkpFields(tt.args.st, tt.args.skipFields)
			if err != nil {
				t.Errorf("%s", err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("CreateStructBySkpFields differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func Test_MapToStruct(t *testing.T) {
	type args struct {
		targetMap    map[string]string
		targetStruct interface{}
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "convert map to struct test",
			args: args{
				targetMap: map[string]string{
					"ID":            "xxxxxx",
					"Name":          "test",
					"LogoImageURL":  "http://logo.hoge",
					"MapSearchWord": "word",
				},
				targetStruct: Test{},
			},
			want: struct {
				ID            string `json:"ID" yaml:"ID"`
				Name          string `json:"Name" yaml:"Name"`
				LogoImageURL  string `json:"LogoImageURL" yaml:"LogoImageURL"`
				MapSearchWord string `json:"MapSearchWord" yaml:"MapSearchWord"`
			}{
				ID:            "xxxxxx",
				Name:          "test",
				LogoImageURL:  "http://logo.hoge",
				MapSearchWord: "word",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapToStruct(tt.args.targetMap, tt.args.targetStruct)
			if err != nil {
				t.Errorf("%s", err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("MapToStruct differs: (-got +want)\n%s", diff)
			}
		})
	}
}
