<?php
declare(strict_types = 1);
error_reporting(E_ALL);

require_once 'cove/helper.php';
require_once 'sienna/musicbrainz.php';
require_once 'sienna/youtube.php';

function yt_result(string $s_query): string {
   $m_query['search_query'] = $s_query;
   $s_res = 'https://www.youtube.com/results?' . http_build_query($m_query);
   echo $s_res, "\n";
   $s_get = file_get_contents($s_res);
   preg_match('!/watch\?v=[^"]*!', $s_get, $a_mat);
   return $a_mat[0];
}

if ($argc != 2) {
   echo <<<eof
usage:
musicbrainz-views.php <URL>

examples:
https://musicbrainz.org/release-group/d03bb6b1-d7b4-38ea-974e-847cbb31dca4
https://musicbrainz.org/release/7a629d52-6a61-3ea1-a0a0-dd50bdef63b4
eof;
   exit(1);
}

$s_url = $argv[1];
$s_mbid = basename($s_url);

if (str_contains($s_url, 'release-group')) {
   # RELEASE GROUP
   $a_releases = mb_decode_group($s_mbid);
   $n_re = 0;
   foreach ($a_releases as $n_idx => $o_cur) {
      $n_re = mb_reduce_group($n_re, $o_cur, $n_idx, $a_releases);
   }
   $o_re = $a_releases[$n_re];
} else {
   # RELEASE
   $o_re = mb_decode_release($s_mbid);
}

foreach ($o_re->{'artist-credit'} as $o_artist) {
   $a_out[] = $o_artist->name;
}

$s_artists = implode(' ', $a_out);

foreach ($o_re->media as $o_media) {
   foreach ($o_media->tracks as $o_track) {
      $s_url = yt_result($s_artists . ' ' . $o_track->title);
      $o = new YouTubeViews($s_url);
      echo $o->color(), "\n";
      usleep(500_000);
   }
}