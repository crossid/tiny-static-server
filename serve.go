package main

import (
  "flag"
  "log"
  "net/http"
  "os"
)

type (
  onlyFilesFS struct {
    fs    http.FileSystem
    index string
  }

  neuteredReaddirFile struct {
    http.File
  }
)

// Dir returns a http.Filesystem that can be used by http.FileServer().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(index, root string, listDirectory bool) http.FileSystem {
  fs := http.Dir(root)
  if listDirectory {
    return fs
  }
  return &onlyFilesFS{fs, index}
}

// Conforms to http.Filesystem
// Tries to open the given file name, otherwise fallback to index.html
// This is required by SPA apps that perform internal routing
func (fs onlyFilesFS) Open(name string) (http.File, error) {
  f, err := fs.fs.Open(name)
  if err != nil {
    f, err = fs.fs.Open(fs.index)
  }

  if err != nil {
    return nil, err
  }
  return neuteredReaddirFile{f}, nil
}

// Overrides the http.File default implementation
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
  // this disables directory listing
  return nil, nil
}

// Conforms http.FileSystem interface
type localFileSystem struct {
  http.FileSystem
  root string
}

func LocalFile(index, root string, listDir bool) *localFileSystem {
  return &localFileSystem{
    FileSystem: Dir(index, root, listDir),
    root:       root,
  }
}

func main() {
  dir := flag.String("dir", "/serve", "location of the dir to serve")
  bind := flag.String("bind", ":8005", "bind address")
  flag.Parse()

  // admin
  adminFS := LocalFile("admin.html", *dir, false)
  adminHandler := http.FileServer(*adminFS)
  adminHandler = http.StripPrefix("/admin/", adminHandler)
  http.Handle("/admin/", adminHandler)

  // portal
  portalFS := LocalFile("portal.html", *dir, false)
  portalHandler := http.FileServer(*portalFS)
  portalHandler = http.StripPrefix("/portal/", portalHandler)
  http.Handle("/portal/", portalHandler)

  // common
  commonFS := LocalFile("", *dir, false)
  commonHandler := http.FileServer(*commonFS)
  http.Handle("/common/", commonHandler)

  log.Printf("Serving from dir %s\n", *dir)
  panic(http.ListenAndServe(*bind, nil))
}
