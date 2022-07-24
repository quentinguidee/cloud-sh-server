package storage

import (
	"mime/multipart"
	"path"
)

func DetectFileType(nodeName string) string {
	// TODO: Shebang detection

	extension := path.Ext(nodeName)

	switch extension {
	case ".avi":
		return "avi"
	case ".babelrc":
		return "babel"
	case ".bmp":
		return "bmp"
	case ".c":
		return "c"
	case ".cpp", ".cxx":
		return "cpp"
	case ".css":
		return "css"
	case ".word", ".odt", ".doc", ".docx":
		return "document"
	case ".flac":
		return "flac"
	case ".gitignore", ".gitkeep":
		return "git"
	case ".go":
		return "go"
	case ".html", ".htm":
		return "html"
	case ".js":
		return "javascript"
	case ".jpeg":
		return "jpeg"
	case ".jpg":
		return "jpg"
	case ".json":
		return "json"
	case ".kt":
		return "kotlin"
	case ".md":
		return "markdown"
	case ".mkv":
		return "mkv"
	case ".mov":
		return "mov"
	case ".mp3":
		return "mp3"
	case ".mp4":
		return "mp4"
	case ".ml", ".mli":
		return "ocaml"
	case ".ogg":
		return "ogg"
	case ".pdf":
		return "pdf"
	case ".php":
		return "php"
	case ".png":
		return "png"
	case ".ppt", ".ppdb", ".odp":
		return "presentation"
	case ".py":
		return "python"
	case ".raw":
		return "raw"
	case ".tsx":
		return "react"
	case ".rb":
		return "ruby"
	case ".sass", ".scss":
		return "sass"
	case ".sc":
		return "scala"
	case ".sh":
		return "shell"
	case ".xls", ".xlsx", ".ods":
		return "spreadsheet"
	case ".ts":
		return "typescript"
	case ".wav":
		return "wav"
	case ".lock":
		return "yarn"
	default:
		return "file"
	}
}

func DetectFileMime(file *multipart.FileHeader) string {
	return file.Header.Get("Content-Type")
}
