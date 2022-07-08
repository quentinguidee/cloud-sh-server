package storage

import (
	"mime/multipart"
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

func DetectFileMime(file *multipart.FileHeader) string {
	return file.Header.Get("Content-Type")
}
