package database

import (
	"database/internal/database/commands"
	"database/internal/database/storage"
	"database/internal/shared/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestDatabase_Handle(t *testing.T) {
	t.Parallel()

	type fields struct {
		logger   *logger.MockLogger
		engine   *MockEngine
		analyzer *MockAnalyzer
		parser   *MockParser
	}

	type args struct {
		input string
	}

	tests := []struct {
		name    string
		args    args
		prepare func(*fields)
		want    string
	}{
		{
			name: "Should succeed execute get operation",
			args: args{input: "GET key"},
			prepare: func(f *fields) {
				f.logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

				f.parser.EXPECT().
					Parse("GET key").
					Return([]string{"GET", "key"}, nil).
					Times(1)

				f.analyzer.EXPECT().
					Analyze([]string{"GET", "key"}).
					Return(commands.Command{
						Operation: "GET",
						Arguments: []string{"key"},
					}, nil).
					Times(1)

				f.engine.EXPECT().
					Execute(commands.Command{
						Operation: "GET",
						Arguments: []string{"key"},
					}, true).Return(storage.Result{Out: "value"}, nil)
			},
			want: "value",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctr := gomock.NewController(t)
			defer t.Cleanup(func() {
				ctr.Finish()
			})

			f := fields{
				logger:   logger.NewMockLogger(ctr),
				engine:   NewMockEngine(ctr),
				analyzer: NewMockAnalyzer(ctr),
				parser:   NewMockParser(ctr),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			db := &Database{
				logger:   f.logger,
				engine:   f.engine,
				analyzer: f.analyzer,
				parser:   f.parser,
			}

			got := db.Handle(tt.args.input)

			is := require.New(t)

			is.Equal(tt.want, got)
		})
	}
}
