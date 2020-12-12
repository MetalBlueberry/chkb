package layerFile

import (
	"github.com/MetalBlueberry/chkb/pkg/chkb"
	"fmt"
	"strings"

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
	case chkb.KbActionPushLayer, chkb.KbActionPopLayer:
		str := layerString(lf.mapper)
		_, err := fmt.Fprint(lf, str)
		return true, err
	}
	return false, nil
}

func layerString(kb *chkb.Keyboard) string {
	builder := strings.Builder{}
	for i := range kb.Mapper.Layers {
		for name := range kb.Config.Layers {
			if kb.Mapper.Layers[i] == kb.Config.Layers[name] {
				builder.WriteString(name)
				builder.WriteString(" > ")
			}
		}
	}
	builder.WriteString("\n")
	return builder.String()
}
