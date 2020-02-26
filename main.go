package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// MyFile ...
type MyFile struct {
	FileName string
	Length   int64
	Source   io.Writer
	Files    io.Reader
}

// m := &MyFile{ReadWriter: file, FileName: line}

// CompFile ...
type CompFile struct {
	Name string
	Data chan []byte
}

func main() {
	ch := make(chan interface{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(out chan<- interface{}) {
		defer wg.Done()
		ch <- "hello world!"
		log.Println("from func")
	}(ch)

	go func() {
		defer wg.Done()
		value := <-ch
		log.Println(value)
	}()

	err := Compress("simple.txt")
	if err != nil {
		log.Println(err)
	}

	// adding single file to tar archive
	err = ArchiveTar("simple2.txt")
	if err != nil {
		log.Println(err)
	}

	// archive simple file
	err = ArchiveZip("test")
	if err != nil {
		log.Println(err)
	}

	// archive a archive
	err = ArchiveZip("simple2.tar")
	if err != nil {
		log.Println(err)
	}

	// compress archive
	err = Compress("simple2.tar")
	if err != nil {
		log.Println(err)
	}

	// data := make([]string, 0)
	// data = append(data, "simple.txt", "simple3.txt", "simple2.txt")
	// err = CreateArchTar(data)
	// if err != nil {
	// 	log.Println(err)
	// }

	UnPack("asimple43.tar.gz")

	wg.Wait()
}

// UnPack ...
func UnPack(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error :", err)
	}
	defer file.Close()
	UnTar(file)
	stat, _ := file.Stat()
	if stat.IsDir() {
		log.Println("is dir")
	} else {
		log.Println("it is a file")
	}

	log.Println("congratulations unpacked")
	return err
}

// UnTar ...
func UnTar(r io.Reader) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err != nil {
			log.Println(err)
			return
		}
		switch header.Typeflag {
		case tar.TypeDir:
			log.Println("dir")
		case tar.TypeReg:
			log.Println("regular")
		}
	}
}

// createArchDir ...
func createArchDir(path string) error {
	if path == "dir" {
		os.Mkdir("name_of_dir", 0600)
		// walk to the 
	} else {
		os.Create("filename")
		// than call a tar function
	}
	return nil
}

// Compress ...
func Compress(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	out, err := os.Create(filename + ".gz")
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	gzout := gzip.NewWriter(out)
	_, _ = gzout.Write(in)
	gzout.Close()
	return err
}

// ArchiveTar ...
func ArchiveTar(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}
	flName := []byte(filename)
	out, err := os.Create(string(flName[:len(flName)-4]) + ".tar")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer out.Close()

	tarout := tar.NewWriter(out)
	header := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(in)),
	}

	if err := tarout.WriteHeader(header); err != nil {
		log.Fatal(err)
	}

	if _, err := tarout.Write(in); err != nil {
		log.Fatal(err)
	}

	if err := tarout.Close(); err != nil {
		log.Fatal(err)
	}

	return err
}

// CreateArchTar ...
func CreateArchTar(files []string) error {
	file, err := os.Create("a" + "simple43" + ".tar.gz")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for i := range files {
		if err := addFile(tw, files[i]); err != nil {
			log.Fatalln(err)
		}
	}
	return err
}

func addFile(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		header := new(tar.Header)
		header.Name = path
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

// CreateArchZip ...
func CreateArchZip(filename, fl2 string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	in2, err := ioutil.ReadFile(fl2)
	if err != nil {
		log.Println(err)
	}

	out, err := os.Create("a" + filename + ".zip")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer out.Close()

	zipout := zip.NewWriter(out)
	z, _ := zipout.Create(filename)
	z.Write(in)
	z.Write(in2)
	zipout.Close()
	return err
}

// ArchiveZip ...
func ArchiveZip(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	out, err := os.Create(filename + ".zip")
	if err != nil {
		log.Println("Error: ", err)
	}
	defer out.Close()

	tarout := zip.NewWriter(out)
	z, _ := tarout.Create(filename)

	if _, err := z.Write(in); err != nil {
		log.Fatal(err)
	}

	if err := tarout.Close(); err != nil {
		log.Fatal(err)
	}

	return err
}
