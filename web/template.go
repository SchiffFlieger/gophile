package web

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/SchiffFlieger/gophile/auth"
)

type IndexTemplateData struct {
	Error string
}

type GophileTemplateData struct {
	User        string
	CurrentPath string
	FileTable   []FileTemplateData
	ActionText  string
}

type FileTemplateData struct {
	IsDir    bool
	FileName string
	FileSize string
	FileDate string
}

type SortDirs []FileTemplateData

type TemplateManager struct {
	index     *template.Template
	impressum *template.Template
	gophile   *template.Template
}

var GlobalTemplateManager *TemplateManager

func (s SortDirs) Len() int      { return len(s) }
func (s SortDirs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortDirs) Less(i, j int) bool {
	if s[i].IsDir == s[j].IsDir {
		return s[i].FileName < s[j].FileName
	} else {
		return s[i].IsDir
	}
}

// Creates a new template manager. All template files are cached inside the template manager.
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		index:     template.Must(template.New("index").Parse(INDEX_HTML)),
		impressum: template.Must(template.New("impressum").Parse(IMPRESSUM_HTML)),
		gophile:   template.Must(template.New("gophile").Parse(GOPHILE_HTML)),
	}
}

// Shows the landing page.
func ShowIndex(w http.ResponseWriter, msg string) {
	GlobalTemplateManager.index.Execute(w, IndexTemplateData{Error: msg})
}

// Shows the impressum.
func ShowImpressum(w http.ResponseWriter) {
	GlobalTemplateManager.impressum.Execute(w, nil)
}

// Shows the default page with the files and folders.
func ShowGophile(w http.ResponseWriter, data GophileTemplateData) {
	GlobalTemplateManager.gophile.Execute(w, data)
}

// Creates the struct used for filling the files and folders template.
func CreateFileData(user, currentpath, actionText string) GophileTemplateData {
	var fileData []FileTemplateData
	dir, _ := os.Open(path.Join(auth.USERDIR, user, currentpath))
	defer dir.Close()
	files, _ := dir.Readdir(-1)

	for _, f := range files {
		fileSize := "-"
		if !f.IsDir() {
			fileSize = intToBytes(f.Size())
		}
		fdata := FileTemplateData{FileName: f.Name(),
			FileSize: fileSize,
			FileDate: f.ModTime().Format("02.01.2006 15:04:05"),
			IsDir:    f.IsDir()}
		fileData = append(fileData, fdata)
	}
	sort.Sort(SortDirs(fileData))
	if len(files) == 0 && actionText == "" {
		actionText = "Directory is empty."
	}

	return GophileTemplateData{User: user, CurrentPath: currentpath, FileTable: fileData, ActionText: actionText}
}

// Converts an integer to a human readable file size string.
func intToBytes(size int64) string {
	sizeLength := size
	units := [...]string{"B", "KB", "MB", "GB"}
	order := 0
	for sizeLength >= 1024 && order < len(units) {
		order++
		sizeLength /= 1024
	}
	return strconv.Itoa(int(sizeLength)) + " " + units[order]
}
