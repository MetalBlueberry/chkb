package layerFile

import (
	"MetalBlueberry/cheap-keyboard/pkg/chkb"
	"fmt"
	"strings"

	"github.com/spf13/afero"
)

type LayerFile struct {
	mapper *chkb.Mapper

	afero.File
}

func NewLayerFile(fs afero.Fs, kb *chkb.Mapper, fileName string) (*LayerFile, error) {
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

func layerString(kb *chkb.Mapper) string {
	builder := strings.Builder{}
	for i := range kb.Layers {
		for name := range kb.LayerBook {
			if kb.Layers[i] == kb.LayerBook[name] {
				builder.WriteString(name)
				builder.WriteString(" > ")
			}
		}
	}
	builder.WriteString("\n")
	return builder.String()
}
