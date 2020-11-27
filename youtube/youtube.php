<?php
declare(strict_types = 1);
extension_loaded('openssl') or die('openssl');
require_once 'cove/helper.php';

class YouTubeInfo {
   function __construct(string $watch_s) {
      # part 1
      $query_s = parse_url($watch_s, PHP_URL_QUERY);
      parse_str($query_s, $query_m);
      # part 2
      $this->id = $query_m['v'];
      # part 3
      $info_s = 'https://www.youtube.com/get_video_info?video_id=' . $this->id;
      echo $info_s, "\n";
      # part 4
      $get_s = file_get_contents($info_s);
      parse_str($get_s, $get_m);
      # part 5
      $resp_s = $get_m['player_response'];
      # part 6
      $resp_o = json_decode($resp_s);
      if (! property_exists($resp_o, 'microformat')) {
         return;
      }
      foreach ($resp_o->microformat->playerMicroformatRenderer as $k => $v) {
         $this->$k = $v;
      }
   }
}

function format_number(float $n): string {
   $n2 = (int)(log10($n) / 3);
   return sprintf('%.3f', $n / 1e3 ** $n2) . ['', ' K', ' M', ' B'][$n2];
}

class YouTubeViews extends YouTubeInfo {
   function color(): string {
      if (! property_exists($this, 'viewCount')) {
         return 'undefined property: viewCount';
      }
      $views_n = (int)($this->viewCount);
      $then_n = strtotime($this->publishDate);
      $now_n = time();
      $diff_n = ($now_n - $then_n) / 60 / 60 / 24 / 365;
      $rate_n = $views_n / $diff_n;
      $rate_s = format_number($rate_n);
      if ($rate_n > 8_000_000) {
         return 'RED ' . color_red($rate_s);
      }
      return 'GREEN ' . color_green($rate_s);
   }
}
