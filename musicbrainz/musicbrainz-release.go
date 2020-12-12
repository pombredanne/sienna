package main

import (
   "log"
   "sienna/musicbrainz"
   "os"
   "path"
   "strings"
)

func main() {
   if len(os.Args) != 2 {
      println(`Usage:
musicbrainz-release <URL>

URL:
https://musicbrainz.org/release/7cc21f46-16b4-4479-844c-e779572ca834
https://musicbrainz.org/release-group/67898886-90bd-3c37-a407-432e3680e872`)
      os.Exit(1)
   }
   url_s := os.Args[1]
   mbid_s := path.Base(url_s)
   dec_o := decode.NewDecode(mbid_s)
   if strings.Contains(url_s, "release-group") {
      rel_a, e := dec_o.Group()
      if e != nil {
         log.Fatal(e)
      }
      rel_n := 0
      foreach ($rel_a as $idx_n => $cur_o) {
         $rel_n = MusicBrainzReduce($rel_n, $cur_o, $idx_n, $rel_a);
      }
      $rel_o = $rel_a[$rel_n];
   } else {
      $rel_o = $dec_o->release();
   }
   $min_n = 179.5 * 1000;
   $max_n = 15 * 60 * 1000;
   $album_m[$rel_o->title]['@date'] = $rel_o->date;
   foreach ($rel_o->media as $media_o) {
      foreach ($media_o->tracks as $track_o) {
         $len_n = $track_o->length;
         if ($len_n < $min_n) {
            $note_s = 'short';
         } else if ($len_n > $max_n) {
            $note_s = 'long';
         } else {
            $note_s = '';
         }
         $track_m = &$album_m[$rel_o->title];
         if (key_exists($track_o->title, $track_m)) {
            $track_m[$track_o->number . '. ' . $track_o->title] = $note_s;
         } else {
            $track_m[$track_o->title] = $note_s;
         }
      }
   }
   $opt_n = JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE;
   $rate_s = json_encode($album_m, $opt_n);
   echo $rate_s, "\n";
}