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
	"fmt"
	"strings"

	"github.com/MetalBlueberry/chkb/pkg/chkb"

	"github.com/spf13/afero"
)

type LayerFile struct {
	mapper *chkb.Keyboard

	afero.File
}

func NewLayerFile(fs afero.Fs, kb *chkb.Keyboard, fileName string) (*LayerFile, error) {
	lf := &LayerFile{
		mapper: kb,
	}

	file, err := fs.Create(fileName)
	if err != nil {
		return nil, err
	}

	lf.File = file

	return lf, nil
}

func (lf *LayerFile) Deliver(event chkb.MapEvent) (handled bool, err error) {
	switch event.Action {
	case chkb.KbActionPushLayer, chkb.KbActionPopLayer, chkb.KbActionChangeLayer:
		str := layerString(lf.mapper)
		_, err := fmt.Fprintln(lf, str)
		return true, err
	}
	return false, nil
}

func layerString(kb *chkb.Keyboard) string {
	fragments := make([]string, len(kb.Mapper.Layers))
	for i := range kb.Mapper.Layers {
		for name := range kb.Config.Layers {
			if kb.Mapper.Layers[i] == kb.Config.Layers[name] {
				fragments[i] = name
			}
		}
	}
	return strings.Join(fragments, ">")
}
