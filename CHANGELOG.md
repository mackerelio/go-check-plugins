# Changelog

## 0.14.1 (2017-10-26)

* fix check-disk options: -x, -X, -p, -N #200 (plaster)


## 0.14.0 (2017-10-12)

* [check-log]Show matched filenames when we use `--return` option #197 (syou6162)


## 0.13.0 (2017-09-27)

* build with Go 1.9 #195 (astj)
* [check-log] Add check-first option and skip the log file on the first run on default #190 (edangelion)


## 0.12.0 (2017-08-23)

* add check-disk to package #192 (astj)
* add check-disk #178 (edangelion)
* [check-postgresql] Add dbname to postgresql-setting #189 (edangelion)


## 0.11.1 (2017-08-02)

* Remove check-ssh binary #186 (astj)


## 0.11.0 (2017-07-26)

* [check-http] Add -i flag to specify source IP #184 (mattn)
* [check-http] Add -s option to map HTTP status #183 (mattn)


## 0.10.4 (2017-06-07)

* v2 packages (rpm and deb) #175 (Songmu)
* [check-log]  When specified multiple pattern, perform search that satisfies all conditions #174 (a-know)


## 0.10.3 (2017-05-16)

* [ntpoffset] support chronyd #166 (Songmu)
* [check-ssh] fix the problem that check-ssh cannot invoke SSH connection #171 (astj)


## 0.10.2 (2017-05-15)

* [experimental] update release scripts #168 (Songmu)


## 0.10.1 (2017-04-27)

* use wmi query instead of running wmic command #157 (mattn)
* Use golib/pluginutil.PluginWorkDir() #163 (astj)


## 0.10.0 (2017-04-06)

* bump go to 1.8 #159 (astj)


## 0.9.7 (2017-03-27)

* check lower-case driver letters #156 (mattn)


## 0.9.6 (2017-03-22)

* Change directory structure convention of each plugin #144 (Songmu)
* run tests under ./check-XXX/lib #152 (astj)
* fix test for AppVayor #154 (daiksy)


## 0.9.5 (2017-03-09)

* add appveyor.yml and fix failing tests on windows #145 (Songmu)
* [check-tcp] connect timeout #146 (Songmu)
* [check-tcp] [bugfix] fix decimal threshold value handling #147 (Songmu)


## 0.9.4 (2017-02-22)

* ensure close temporary file in writeFileAtomically #141 (itchyny)
* Write the file in same directory #142 (mattn)


## 0.9.3 (2017-02-08)

* fix matching for 'Audit Success' and 'Audit Failure' #139 (mattn)


## 0.9.2 (2017-01-25)

* [check-windows-eventlog] add --source-exclude, --message-exclude and --event-id #136 (mattn, daiksy)
* [check-windows-eventlog] remove --event-id and add --event-id-pattern, --event-id-exclude #137 (mattn)


## 0.9.1 (2017-01-11)

* [check-log] support glob in --file option #133 (astj)


## 0.9.0 (2017-01-04)

* add check-windows-eventlog #129 (daiksy)
* [check-log]fix encoding option #131 (daiksy)


## 0.8.1 (2016-11-29)

* Fix state in check-procs #124 (itchyny)
* Fix the links to the document #125 (itchyny)
* remove checking Ignore #126 (mattn)
* [check-log] fix state file path (fix wrong slice copy) #127 (Songmu)


## 0.8.0 (2016-10-27)

* [check-log] improve Windows support #122 (daiksy)


## 0.7.0 (2016-10-18)

* Add option for skip searching logfile content if logfile does not exist #115 (a-know)
* [check-log] write file atomically when saving read position into state file #119 (Songmu)


## 0.6.3 (2016-09-07)

* fix check-mysql replication to detect IO thread 'Connecting' #116 (hiroakis)
* [file-age] Remove unnecessary newline #117 (b4b4r07)


## 0.6.2 (2016-06-23)

* Fixed argument parser error: #113 (karupanerura)


## 0.6.1 (2016-05-13)

* no panic check-masterha when not found the targets, and refactoring #108 (karupanerura)


## 0.6.0 (2016-05-10)

* supported gearman ascii protocol for check-tcp #89 (karupanerura)
* added check-masterha command to check masterha status #90 (karupanerura)
* Plugin for checking AWS SQS queue size #92 (hiroakis)
* fix: rpm should not include dir #98 (stanaka)
* added ssh checker #101 (karupanerura)
* remove 'golang.org/x/tools/cmd/vet' #102 (y-kuno)
* [uptime/procs] `--warn-over/under` is deprecated. use `--warning-over/under` instead #104 (Songmu)
* add aws-sqs-queue-size, cert-file, masterha and ssh into package #105 (Songmu)
* bump up go version to 1.6.2 #106 (stanaka)


## 0.5.2 (2016-03-25)

* Revert "use /usr/bin/check-*" #95 (Songmu)


## 0.5.1 (2016-03-25)

* use /usr/bin/check-* #91 (naokibtn)
* use GOARCH=amd64 for now #93 (Songmu)


## 0.5.0 (2016-03-02)

* add check-solr #46 (supercaracal)
*  add check-jmx-jolokia #68 (y-kuno)
* check-memcached #69 (naokibtn)
* Add scheme option to check-elasticsearch #70 (yano3)
* Check file size #72 (hiroakis)
* Add uptime sub command to check-mysql #73 (kazeburo)
* add check-cert-file #74 (kazeburo)
* [tcp] Add --no-check-certificate #75 (kazeburo)
* Fixed slurp. conn.read returns content with EOF #76 (kazeburo)
* Fix check-tcp IPv6 testcase on OSX(?) #77 (hanazuki)
* check-redis: Report critical if master_link_status is down #79 (hanazuki)
* check-redis: Fixed panic #80 (yoheimuta)
* check-procs: Fixed the counting logic with -p #81 (yoheimuta)
* add check-uptime #82 (Songmu)
* add file-size, jmx-jolokia, memcached, solr, uptime into package config #84 (Songmu)


## 0.4.0 (2016-02-04)

* Fix duplicated help message #58 (hfm)
* add qmail queue check to check-mailq #59 (tnmt)
* Add check-elasticsearch #62 (naokibtn)
* Add check-redis #63 (naokibtn)
* check-procs: check PPID #64 (hfm)


## 0.3.1 (2016-01-07)

* build with go1.5 #43 (Songmu)


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
