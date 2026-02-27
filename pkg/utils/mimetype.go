package utils

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// FileMimeType представляет результат поиска MIME-типа
type FileMimeType struct {
	Extension string
	MimeType  string
	Found     bool
}

func MimeType(data []byte, filename string) *FileMimeType {
	if data == nil && filename == "" {
		return &FileMimeType{
			Found: false,
		}
	}

	result := getFileMimeType(filename)

	if result.Extension == "" {
		index := strings.LastIndex(filename, ".")
		if index >= 0 {
			result.Extension = filename[index:]
		}
	}

	if result.Found {
		return &result
	}

	var contentType string
	mtype := mimetype.Detect(data)
	if mtype != nil {
		return &FileMimeType{
			Extension: mtype.Extension(),
			MimeType:  mtype.String(),
		}
	}

	fromHttp := http.DetectContentType(data)
	if contentType != fromHttp {
		result = getExtensionByMimeType(contentType)
	}

	return &result
}

func getMimeTypesMap() map[string]string {
	return map[string]string{
		// application/*
		".atom": "application/atom+xml",
		".json": "application/json",
		".js":   "application/javascript",
		".bin":  "application/octet-stream",
		".exe":  "application/octet-stream",
		".ogx":  "application/ogg",
		".pdf":  "application/pdf",
		".ps":   "application/postscript",
		".xml":  "application/xml",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".xls":  "application/vnd.ms-excel",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".yaml": "application/x-yaml",
		".yml":  "application/x-yaml",
		".apk":  "application/vnd.android.package-archive",
		".tar":  "application/x-tar",
		".rar":  "application/x-rar-compressed",
		".zip":  "application/zip",
		".gz":   "application/gzip",
		".7z":   "application/x-7z-compressed",
		".bz2":  "application/x-bzip2",

		// audio/*
		".aac":  "audio/aac",
		".mp3":  "audio/mpeg",
		".oga":  "audio/ogg",
		".ogg":  "audio/ogg",
		".wav":  "audio/x-wav",
		".weba": "audio/webm",
		".flac": "audio/flac",
		".m4a":  "audio/x-m4a",
		".opus": "audio/opus",
		".amr":  "audio/amr",
		".aiff": "audio/x-aiff",
		".ape":  "audio/x-ape",

		// image/*
		".gif":  "image/gif",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".svg":  "image/svg+xml",
		".tiff": "image/tiff",
		".tif":  "image/tiff",
		".ico":  "image/vnd.microsoft.icon",
		".webp": "image/webp",
		".heif": "image/heif",
		".heic": "image/heic",
		".avif": "image/avif",
		".bmp":  "image/bmp",

		// model/*
		".igs":  "model/iges",
		".iges": "model/iges",
		".msh":  "model/mesh",
		".mesh": "model/mesh",
		".wrl":  "model/vrml",
		".vrml": "model/vrml",
		".x3db": "model/x3d+binary",
		".x3dv": "model/x3d+vrml",
		".x3d":  "model/x3d+xml",
		".obj":  "model/obj",
		".stl":  "model/stl",
		".u3d":  "model/u3d",
		".glb":  "model/gltf-binary",
		".dae":  "model/vnd.collada+xml",

		// text/*
		".css":      "text/css",
		".csv":      "text/csv",
		".html":     "text/html",
		".htm":      "text/html",
		".txt":      "text/plain",
		".md":       "text/markdown",
		".markdown": "text/markdown",
		".php":      "text/php",
		".rtf":      "text/rtf",
		".vcard":    "text/vcard",
		".vcf":      "text/vcard",
		".vtt":      "text/vtt",
		".go":       "text/x-go",

		// video/*
		".mpeg": "video/mpeg",
		".mpg":  "video/mpeg",
		".mp4":  "video/mp4",
		".ogv":  "video/ogg",
		".mov":  "video/quicktime",
		".webm": "video/webm",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".avi":  "video/x-msvideo",
		".3gp":  "video/3gpp",
		".3g2":  "video/3gpp2",
		".mkv":  "video/x-matroska",
		".m4v":  "video/x-m4v",

		// x/*
		".dvi":   "application/x-dvi",
		".latex": "application/x-latex",
		".ttf":   "application/x-font-ttf",
		".swf":   "application/x-shockwave-flash",
		".sit":   "application/x-stuffit",
	}
}

func getFileMimeType(filename string) FileMimeType {
	// Extension
	ext := strings.ToLower(filepath.Ext(filename))

	mimeTypes := getMimeTypesMap()

	// MimeType
	if mimeType, found := mimeTypes[ext]; found {
		return FileMimeType{
			Extension: ext,
			MimeType:  mimeType,
			Found:     true,
		}
	}

	// Empty
	return FileMimeType{
		Extension: ext,
		MimeType:  "",
		Found:     false,
	}
}

func getExtensionByMimeType(mimeType string) FileMimeType {
	mimeTypeLower := strings.ToLower(strings.TrimSpace(mimeType))

	mimeTypes := getMimeTypesMap()

	for ext, mt := range mimeTypes {
		if strings.ToLower(mt) == mimeTypeLower {
			return FileMimeType{
				Extension: ext,
				MimeType:  mimeType,
				Found:     true,
			}
		}
	}

	return FileMimeType{
		Extension: "",
		MimeType:  mimeType,
		Found:     false,
	}
}
