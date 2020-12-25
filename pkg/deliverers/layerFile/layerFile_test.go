/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package layerFile

import (
	"io/ioutil"
	"testing"

	"github.com/MetalBlueberry/chkb/pkg/chkb"

	"github.com/spf13/afero"
)

func TestLayerFile_Deliver(t *testing.T) {
	book := chkb.Config{
		Layers: chkb.LayerBook{
			"l0": &chkb.Layer{},
			"l1": &chkb.Layer{},
			"l2": &chkb.Layer{},
		},
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
					f, err := NewLayerFile(mfs, chkb.NewKeyboard(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0
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
					f, err := NewLayerFile(mfs, chkb.NewKeyboard(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					f.mapper.PushLayer("l2")
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0>l2
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
					f, err := NewLayerFile(mfs, chkb.NewKeyboard(book, "l0"), "file")
					if err != nil {
						panic(err)
					}
					f.mapper.PushLayer("l2")
					return f
				},
			},
			wantErr:     false,
			wantHandled: true,
			wantFileContent: `l0>l2
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
					f, err := NewLayerFile(mfs, chkb.NewKeyboard(book, "l0"), "file")
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
