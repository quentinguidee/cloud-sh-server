package storage

import (
	"mime/multipart"
	"path"
)

func DetectFileType(nodeName string) string {
	// TODO: Shebang detection

	extension := path.Ext(nodeName)

	switch extension {
	case ".babelrc":
		return "babel"
	case ".c":
		return "c"
	case ".cpp", ".cxx":
		return "cpp"
	case ".css":
		return "css"
	case ".gitignore", ".gitkeep":
		return "git"
	case ".go":
		return "go"
	case ".html", ".htm":
		return "html"
	case ".js":
		return "javascript"
	case ".json":
		return "json"
	case ".kt":
		return "kotlin"
	case ".lock":
		return "yarn"
	case ".md":
		return "markdown"
	case ".ml", ".mli":
		return "ocaml"
	case ".mp3":
		return "mp3"
	case ".wav":
		return "wav"
	case ".flac":
		return "flac"
	case ".ogg":
		return "ogg"
	case ".mp4":
		return "mp4"
	case ".avi":
		return "avi"
	case ".mkv":
		return "mkv"
	case ".mov":
		return "mov"
	case ".png":
		return "png"
	case ".jpg":
		return "jpg"
	case ".jpeg":
		return "jpeg"
	case ".bmp":
		return "bmp"
	case ".raw":
		return "raw"
	case ".php":
		return "php"
	case ".py":
		return "python"
	case ".rb":
		return "ruby"
	case ".sass", ".scss":
		return "sass"
	case ".sc":
		return "scala"
	case ".sh":
		return "shell"
	case ".ts":
		return "typescript"
	case ".tsx":
		return "react"
	default:
		return "file"
	}
}

func DetectFileMime(file *multipart.FileHeader) string {
	return file.Header.Get("Content-Type")
}
