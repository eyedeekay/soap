package unciv

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"log"
)

var invalidFilenameChars = []string{"..", "*"}

func ReplaceInvalidChars(filename string) (string, error) {
	if len(filename) > 253 {
		return "", fmt.Errorf("Too-long filename submitted, potential attack")
	}
	spf := ""
	for _, i := range invalidFilenameChars {
		spf = strings.ReplaceAll(filename, i, "")
	}
	return filepath.Clean(spf), nil
}

var tx map[string]*sync.Mutex

func writeFile(path, filename string, bytes []byte) error {
	if tx == nil {
		tx = make(map[string]*sync.Mutex)
	}
	txe, ok := tx[filename]
	if ok {
		if !txe.TryLock() {
			return fmt.Errorf("Lock already held, file is open for writing")
		}
	} else {
		tx[filename] = &sync.Mutex{}
		return writeFile(path, filename, bytes)
	}
	fullFileName := filepath.Join(path, filename)
	err := ioutil.WriteFile(fullFileName, bytes, 0644)
	defer txe.Unlock()
	return err
}

func deleteFile(path, filename string) error {
	if tx == nil {
		tx = make(map[string]*sync.Mutex)
	}
	txe, ok := tx[filename]
	if ok {
		if !txe.TryLock() {
			return fmt.Errorf("Lock already held, file is open for writing")
		}
	} else {
		tx[filename] = &sync.Mutex{}
		return deleteFile(path, filename)
	}
	fullFileName := filepath.Join(path, filename)
	err := os.RemoveAll(fullFileName)
	defer txe.Unlock()
	return err
}

func readFile(path, filename string) ([]byte, error) {
	filepath := filepath.Join(path, filename)
	return ioutil.ReadFile(filepath)
}

func stat(path, filename string) (os.FileInfo, error) {
	filename = filepath.Join(path, filename)
	return os.Stat(filename)
}

type UncivServerInterface interface {
	Alive(w http.ResponseWriter, r *http.Request) error
	Put(w http.ResponseWriter, r *http.Request) error
	Get(w http.ResponseWriter, r *http.Request) error
	Delete(w http.ResponseWriter, r *http.Request) error
	Files(w http.ResponseWriter, r *http.Request) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type UncivServer struct {
	Directory string
}

func (u *UncivServer) directory() string {
	if u.Directory == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		u.Directory = filepath.Join(dir, "UncivServer")
	}
	os.MkdirAll(u.Directory, 0755)
	return u.Directory
}

func (u *UncivServer) writeFile(filename string, bytes []byte) error {
	return writeFile(u.directory(), filename, bytes)
}

func (u *UncivServer) deleteFile(filename string) error {
	return deleteFile(u.directory(), filename)
}

func (u *UncivServer) readFile(filename string) ([]byte, error) {
	return readFile(u.directory(), filename)
}

func (u *UncivServer) stat(filename string) (os.FileInfo, error) {
	return stat(u.directory(), filename)
}

func (u *UncivServer) Alive(w http.ResponseWriter, r *http.Request) (err error) {
	w.Write([]byte("true"))
	return err
}

/*  put("/files/{fileName}") {
    val fileName = call.parameters["fileName"] ?: throw Exception("No fileName!")
    log.info("Receiving file: ${fileName}")
    val file = File(fileFolderName, fileName)
    withContext(Dispatchers.IO) {
        file.outputStream().use {
            call.request.receiveChannel().toInputStream().copyTo(it)
        }
    }
    call.respond(HttpStatusCode.OK)
}*/
func (u *UncivServer) Put(w http.ResponseWriter, r *http.Request) (err error) {
	//get the "fileName" parameter
	filename, err := ReplaceInvalidChars(fileName(r))
	if filename == "" || filename == "." {
		w.WriteHeader(http.StatusUnauthorized)
		err = fmt.Errorf("no filename passed to server")
	}
	if err != nil {
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	err = u.writeFile(filename, bodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	return
}

/*  get("/files/{fileName}") {
	val fileName = call.parameters["fileName"] ?: throw Exception("No fileName!")
	log.info("File requested: $fileName")
	val file = File(fileFolderName, fileName)
	if (!file.exists()) {
		log.info("File $fileName not found")
		call.respond(HttpStatusCode.NotFound, "File does not exist")
		return@get
	}
	val fileText = withContext(Dispatchers.IO) { file.readText() }
	call.respondText(fileText)
}*/
func (u *UncivServer) Get(w http.ResponseWriter, r *http.Request) (err error) {
	//get the "fileName" parameter
	filename, err := ReplaceInvalidChars(fileName(r))
	if filename == "" || filename == "." {
		w.WriteHeader(http.StatusNotFound)
		err = fmt.Errorf("invalid filename passed to server")
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err = u.stat(filename); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bytes, err := u.readFile(filename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(bytes)
	return
}

/*  delete("/files/{fileName}") {
    val fileName = call.parameters["fileName"] ?: throw Exception("No fileName!")
    log.info("Deleting file: $fileName")
    val file = File(fileFolderName, fileName)
    if (!file.exists()) {
        call.respond(HttpStatusCode.NotFound, "File does not exist")
        return@delete
    }
    file.delete()
}*/
func (u *UncivServer) Delete(w http.ResponseWriter, r *http.Request) (err error) {
	//get the "fileName" parameter
	filename, err := ReplaceInvalidChars(fileName(r))
	if filename == "" || filename == "." {
		err = fmt.Errorf("invalid filename passed to server")
	}
	if err != nil {
		return
	}
	err = u.deleteFile(filename)
	return
}

func (u *UncivServer) Files(w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "GET":
		err = u.Get(w, r)
		if err != nil {
			log.Println("Get error", err)
		}
		return
	case "PUT":
		err = u.Put(w, r)
		if err != nil {
			log.Println("Put error", err)
		}
		return
	case "DELETE":
		err = u.Delete(w, r)
		if err != nil {
			log.Println("Delete error", err)
		}
		return
	}
	return fmt.Errorf("Method error: method not found")
}

func pathIsFiles(r *http.Request) bool {
	path := strings.TrimLeft(r.URL.Path, "/")
	if strings.HasPrefix(path, "files") {
		return true
	}
	return false
}

func (u *UncivServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formVals := r.ParseForm()
	log.Println(formVals)
	log.Println(r.URL.String())

	if pathIsFiles(r) {
		err := u.Files(w, r)
		if err != nil {
			log.Println("Server Error:", err)
		}
	} else {
		err := u.Alive(w, r)
		if err != nil {
			log.Println("Server Error:", err)
		}
	}
}

func fileName(r *http.Request) string {
	return strings.TrimLeft(r.URL.Path, "/files")
}
