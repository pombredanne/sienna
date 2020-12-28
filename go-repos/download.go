package main

import (
   "encoding/json"
   "golang.org/x/build/repos"
   "net/http"
   "os"
   "sienna/assert"
   "strings"
)

var bad_repo = map[string]bool{
   "golang.org/x/build": true,
   "golang.org/x/crypto": true,
   "golang.org/x/oauth2": true,
   "golang.org/x/tools": true,
}

func Copy(source, dest string) (int64, error) {
   println(source)
   get_o, e := http.Get(source)
   if e != nil {
      return 0, e
   }
   create_o, e := os.Create(dest)
   if e != nil {
      return 0, e
   }
   return create_o.ReadFrom(get_o.Body)
}

func Download() error {
   os.Mkdir("x", os.ModeDir)
   os.Chdir("x")
   for repo_s, repo_o := range repos.ByImportPath {
      if ! repo_o.ShowOnDashboard() {
         continue
      }
      if bad_repo[repo_s] {
         continue
      }
      url_s := "https://api.godoc.org/search?q=" + repo_s + "/"
      println(url_s)
      get_o, e := http.Get(url_s)
      if e != nil {
         return e
      }
      get_m := assert.Map{}
      e = json.NewDecoder(get_o.Body).Decode(&get_m)
      if e != nil {
         return e
      }
      result_a := get_m.A("results")
      for n := range result_a {
         path_s := result_a.M(n).S("path")
         if ! strings.HasPrefix(path_s, "golang.org/x/") {
            continue
         }
         path_a := strings.Split(path_s, "/")
         if len(path_a) > 4 {
            continue
         }
         dest := strings.ReplaceAll(path_s, "/", "-")
         _, e = Copy("https://pkg.go.dev/" + path_s, dest + ".html")
         if e != nil {
            return e
         }
      }
   }
   return nil
}
