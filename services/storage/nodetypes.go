package storage

import (
	"path"
)

func DetectFileType(nodeName string) string {
	// TODO: Shebang detection

	extension := path.Ext(nodeName)

	switch extension {
	case ".c":
		return "c"
	case ".cpp", ".cxx":
		return "cpp"
	case ".py":
		return "python"
	default:
		return "file"
	}
}
