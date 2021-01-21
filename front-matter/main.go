package main

import (
   "bytes"
   "github.com/89z/x"
   "github.com/89z/x/toml"
   "io/ioutil"
   "os"
   "path"
)

func main() {
   if len(os.Args) != 2 {
      println(`front-matter D:\Git`)
      os.Exit(1)
   }
   root := os.Args[1]
   content := path.Join(root, "autumn", "content")
   e := os.Chdir(content)
   check(e)
   dir, e := ioutil.ReadDir(".")
   check(e)
   for _, entry := range dir {
      index_s := path.Join(entry.Name(), "_index.md")
      index_y, e := ioutil.ReadFile(index_s)
      check(e)
      data := bytes.SplitN(index_y, toml_sep, 3)[1]
      front, e := toml.LoadBytes(data)
      check(e)
      if front["_build"] != nil {
         continue
      }
      example := front.A("example")
      exFile := path.Join(root, example.S(0))
      if ! x.IsFile(exFile) {
         println(index_s)
         continue
      }
      example_y, e := ioutil.ReadFile(exFile)
      check(e)
      substr_s := example.S(1)
      substr_y := []byte(substr_s)
      if ! bytes.Contains(example_y, substr_y) {
         println(index_s)
      }
   }
}