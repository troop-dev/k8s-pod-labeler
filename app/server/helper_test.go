package server

import (
	"reflect"
	"testing"
)

func Test_parseB64Map(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
		err  error
	}{
		{
			name: "accepts complex values",
			args: args{in: "eyJjbHVzdGVyLWF1dG9zY2FsZXIua3ViZXJuZXRlcy5pby9zYWZlLXRvLWV2aWN0IjoidHJ1ZSIsImNvbmZpZy5saW5rZXJkLmlvL3NraXAtb3V0Ym91bmQtcG9ydHMiOiI0MzE3In0="},
			want: map[string]string{
				"cluster-autoscaler.kubernetes.io/safe-to-evict": "true",
				"config.linkerd.io/skip-outbound-ports":          "4317",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseB64Map(tt.args.in)

			if err != tt.err {
				t.Errorf("parseB64Map() error = %v, wantErr %v", err, tt.err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseB64Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
