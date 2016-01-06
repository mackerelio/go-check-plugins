# Changelog

## 0.3.0 (2016-01-06)

* add check-postgresql #47 (supercaracal)
* [check-ntpoffset] work on ntp 4.2.2 #50 (naokibtn)
* check-file-age: show time in message #51 (naokibtn)
* add --no-check-certificate options to check_http #52 (cl-lab-k)
* add check-mailq, currently only available for postfix. #54 (tnmt)
* add check-mailq and check-postgresql into package #55 (Songmu)


## 0.2.2 (2015-12-08)

* A plugin to check ntpoffset #38 (hiroakis)
* Check tcp unix domain socket #39 (tkuchiki)
* [check-tcp] add ipv6 support #42 (Songmu)


## 0.2.1 (2015-11-25)

* Fix bugs of check-log #35 (Songmu)
* [check-log] add --no-state option #36 (Songmu)


## 0.2.0 (2015-11-20)

* [check-procs] support `--critical-over=0` and `--warn-over=0` #31 (Songmu)
* add check-tcp #32 (Songmu)


## 0.1.1 (2015-11-12)

* check-procs for windows #24 (mechairoi)
* [bug] [check-log] fix large file handling #27 (matsuzj)


## 0.1.0 (2015-11-05)

* check-log #21 (Songmu)
* Add check-log in the packages #25 (Songmu)


## 0.0.5 (2015-10-26)

* Add mysql in packages

## 0.0.4 (2015-10-26)

* Refactor Mysql checks #20 (Songmu)
* Add Mysql checkes #19 (hiroakis)

## 0.0.3 (2015-10-15)

* reduce binary size by using ldflags #15 (Songmu)
* Remove cgo dependency from check-load #16 (Songmu)
* Add check-load in the packages


## 0.0.2 (2015-10-08)

* fix MatchSelf behaviour #9 (Songmu)


## 0.0.1 (2015-10-07)

* Fix release tools

## 0.0.0 (2015-10-07)

* Initial release
