package main

import (
		"log"
		"os"
		"bufio"
		"path/filepath"
		"io/ioutil"
		"flag"
		"sync"
		"sync/atomic"
)

var fileChan = make([]*os.FileInfo,100)

var toatal int32 = 0

var bom1 int32 = 0

var bom2 int32 = 0

var noBom int32 = 0

var fixBom int32 = 0

var waitGroup = sync.WaitGroup{}
func main(){


		removeBom := flag.Bool("rb", false, "is remove the bom? true will remove the bom.false will only print is the file with bom ")
		dst := flag.String("dst", "", "which folder you nead to scan")

		flag.Parse()
		log.Printf("this action will removeBom:%t",*removeBom)
		if *dst == "" {
				*dst = getCurrentDirectory()
		}
		ReadFile(*dst,*removeBom)
		waitGroup.Wait()
		log.Println("-------------------------")
		log.Printf("total file:%d\n",toatal)
		log.Printf("2bom file:%d\n",bom2)
		log.Printf("1bom file:%d\n",bom1)
		log.Printf("nobom file:%d\n",noBom)
		log.Printf("fixed bom file:%d\n",fixBom)
		log.Println("-------------------------")
		log.Println("scan bom finsh")
}

func getCurrentDirectory() string {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
				log.Fatal(err)
		}
		return dir
}

func ReadFile(path string,removeBom bool)  {
		filepath.Walk(path,
				func(path string, f os.FileInfo, err error) error {
						if f == nil {
								return err
						}
						if !f.IsDir() {
								waitGroup.Add(1)
								go IsBom(path,removeBom)
						}
						return nil
				})
}


func IsBom(name string,removeBom bool)  {
		f, err := os.Open(name)
		defer f.Close()
		checkError(err)
		br := bufio.NewReader(f)
		r1, _, err := br.ReadRune()
		checkError(err)

		r2, _, err := br.ReadRune()
		checkError(err)

		atomic.AddInt32(&toatal,1)
		if r1 == '\uFEFF' && r2 == '\uFEFF'{
				log.Println("2bom file "+ name)
				atomic.AddInt32(&bom2,1)
				if removeBom {
						RemoveBom(name)
				}
				waitGroup.Done()
		}else if r1 == '\uFEFF' && r2 != '\uFEFF'{
				log.Println("1bom file "+ name)
				atomic.AddInt32(&bom1,1)
				if removeBom {
						RemoveBom(name)
				}
				waitGroup.Done()
		} else{
				atomic.AddInt32(&noBom,1)
				//log.Println("file "+ name + ":without bom")
				waitGroup.Done()
		}
}

func RemoveBom(name string)   {
		b, err  := ioutil.ReadFile(name)
		checkError(err)
		raw := b[3:]
		err = ioutil.WriteFile(name,raw,os.ModePerm)
		if err == nil {
				atomic.AddInt32(&fixBom,1)
				log.Println(":remove bom "+name)
		}else{
				log.Println(err)
		}
}


func checkError(err error) {
		if err != nil {
				//log.Println(err)
		}
}
