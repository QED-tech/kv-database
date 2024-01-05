package compute

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	t.Parallel()

	type args struct {
		query string
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Should return tokens with multiple words query",
			args: args{
				query: `SET query "hello world value"`,
			},
			want:    []string{"SET", "query", "hello world value"},
			wantErr: false,
		},
		{
			name: "Should return tokens with multiple words query",
			args: args{
				query: `GET 'multiple words key'`,
			},
			want:    []string{"GET", "multiple words key"},
			wantErr: false,
		},
		{
			name: "Should return tokens for GET+query+value",
			args: args{
				query: `SET query value`,
			},
			want:    []string{"SET", "query", "value"},
			wantErr: false,
		},
		{
			name: "Should return tokens for DEL+query",
			args: args{
				query: `DEL query`,
			},
			want:    []string{"DEL", "query"},
			wantErr: false,
		},
		{
			name: "Should return tokens for DEL and spaces",
			args: args{
				query: `DEL   key   `,
			},
			want:    []string{"DEL", "key"},
			wantErr: false,
		},
		{
			name: "Should return tokens for DEL and spaces",
			args: args{
				query: `DEL   `,
			},
			want:    []string{"DEL"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &Parser{}
			got, err := p.Parse(tt.args.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			is := require.New(t)

			is.Equal(tt.want, got)
		})
	}
}
