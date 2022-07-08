package storage

import (
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
	case ".mp3", ".wav", ".flac", ".ogg":
		return "audio"
	case ".mp4", ".avi", ".mkv", ".mov":
		return "video"
	case ".png", "jpg", "jpeg", "bmp", "raw":
		return "image"
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
