package layerFile

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
)

func TestLayerFile_Deliver(t *testing.T) {
	book := chkb.Config{
		"l0": &chkb.Layer{},
		"l1": &chkb.Layer{},
		"l2": &chkb.Layer{},
	}
	mfs := afero.NewMemMapFs()

	type fields struct {
		File func() *LayerFile
	}
	type args struct {
		event chkb.MapEvent
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantHandled     bool
		wantErr         bool
		wantFileContent string
	}{
		{
			name: "simple",
			args: args{
				event: chkb.MapEvent{
					Action:    chkb.KbActionPushLayer,
					LayerName: "test",
				},
			},
			fields: fields{
				File: func() *LayerFile {
					f, err := NewLayerFile(mfs, chkb.NewMapper(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0 > 
`,
		},
		{
			name: "two layers",
			args: args{
				event: chkb.MapEvent{
					Action:    chkb.KbActionPushLayer,
					LayerName: "test",
				},
			},
			fields: fields{
				File: func() *LayerFile {
					f, err := NewLayerFile(mfs, chkb.NewMapper(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					f.mapper.PushLayer("l2")
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0 > l2 > 
`,
		},
		{
			name: "handle pop",
			args: args{
				event: chkb.MapEvent{
					Action:    chkb.KbActionPopLayer,
					LayerName: "test",
				},
			},
			fields: fields{
				File: func() *LayerFile {
					f, err := NewLayerFile(mfs, chkb.NewMapper(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					f.mapper.PushLayer("l2")
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0 > l2 > 
`,
		},
		{
			name: "ignore event",
			args: args{
				event: chkb.MapEvent{
					Action:    chkb.KbActionDown,
					LayerName: "test",
				},
			},
			fields: fields{
				File: func() *LayerFile {
					f, err := NewLayerFile(mfs, chkb.NewMapper(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					return f
				},
			},
			wantErr:         false,
			wantHandled:     false,
			wantFileContent: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHandled, err := tt.fields.File().Deliver(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("LayerFile.Deliver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotHandled != tt.wantHandled {
				t.Errorf("LayerFile.Deliver() = %v, want %v", gotHandled, tt.wantHandled)
				return
			}
			f, err := mfs.Open("file")
			if err != nil {
				t.Errorf("Cannot open file %v", err)
				return
			}
			defer f.Close()
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Errorf("Cannot read file %v", err)
				return
			}
			if string(content) != tt.wantFileContent {
				t.Errorf("LayerFile.Deliver() out content %v, want %v", string(content), tt.wantFileContent)
				return
			}
		})
	}
}
