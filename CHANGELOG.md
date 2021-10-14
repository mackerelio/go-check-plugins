# Changelog

## 0.41.1 (2021-10-14)

* Bump github.com/aws/aws-sdk-go from 1.39.4 to 1.40.59 #529 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.21.6 to 3.21.9 #527 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.3 to 1.7.4 #521 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.2 to 1.10.3 #520 (dependabot[bot])
* Bump golang.org/x/text from 0.3.6 to 0.3.7 #519 (dependabot[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.3.0 to 3.4.1 #518 (dependabot[bot])


## 0.41.0 (2021-10-06)

* update golib, checkers #525 (yseto)
* [check-log] add search-in-directory option #524 (yseto)
* [check-redis] migrate redis client library to redigo #516 (pyto86pri)


## 0.40.1 (2021-09-29)

* check-mysql: Closes `checkReplication` rows #515 (mechairoi)


## 0.40.0 (2021-08-24)

* [check-mysql] add --tls, --tls-root-cert and --tls-skip-verify options #511 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.38.68 to 1.39.4 #507 (dependabot[bot])


## 0.39.5 (2021-07-06)

* Bump github.com/shirou/gopsutil/v3 from 3.21.5 to 3.21.6 #505 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.38.45 to 1.38.68 #501 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.2 to 1.7.3 #494 (dependabot[bot])


## 0.39.4 (2021-06-23)

* [ci]rewrite check-memcached tests. used docker. #498 (yseto)
* refactor check-log tests. #497 (yseto)
* [ci] run tests on the workflow #495 (lufia)
* [check-disk] upgrade gopsutil to v3 and fix treatment for fstype=none #492 (susisu)


## 0.39.3 (2021-06-03)

* Bump github.com/aws/aws-sdk-go from 1.38.40 to 1.38.45 #489 (dependabot[bot])
* Bump github.com/go-sql-driver/mysql from 1.5.0 to 1.6.0 #473 (dependabot[bot])
* Bump github.com/jmoiron/sqlx from 1.3.1 to 1.3.4 #486 (dependabot[bot])
* Bump github.com/jessevdk/go-flags from 1.4.0 to 1.5.0 #470 (dependabot[bot])
* Bump github.com/lib/pq from 1.10.0 to 1.10.2 #488 (dependabot[bot])
* Bump golang.org/x/text from 0.3.5 to 0.3.6 #474 (dependabot[bot])
* Bump github.com/aws/aws-sdk-go from 1.37.30 to 1.38.40 #487 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 #484 (dependabot[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.2.4 to 3.3.0 #476 (dependabot[bot])
* upgrade go1.14 -> 1.16 #480 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.37.1 to 1.37.30 #468 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.7.0 to 1.7.2 #466 (dependabot[bot])
* Bump github.com/lib/pq from 1.9.0 to 1.10.0 #469 (dependabot[bot])
* [check-mysql] Use go-mysql-driver and sqlx instead of mymysql #464 (nonylene)
* [ci] replace token #461 (yseto)
* [ci] replace mackerel-github-release #459 (yseto)


## 0.39.2 (2021-02-03)

* Bump github.com/aws/aws-sdk-go from 1.36.31 to 1.37.1 #454 (dependabot[bot])
* Closes #455 check-load fix null pointer issue when critical #456 (hurrycaine)
* Bump github.com/aws/aws-sdk-go from 1.36.28 to 1.36.31 #452 (dependabot[bot])


## 0.39.1 (2021-01-21)

* Bump github.com/aws/aws-sdk-go from 1.36.19 to 1.36.28 #447 (dependabot[bot])
* Bump github.com/fsouza/go-dockerclient from 1.6.6 to 1.7.0 #445 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 #448 (dependabot[bot])
* Bump github.com/go-ole/go-ole from 1.2.4 to 1.2.5 #449 (dependabot[bot])
* Bump golang.org/x/text from 0.3.4 to 0.3.5 #443 (dependabot[bot])


## 0.39.0 (2021-01-14)

* Bump github.com/aws/aws-sdk-go from 1.35.35 to 1.36.19 #442 (dependabot[bot])
* Bump github.com/lib/pq from 1.8.0 to 1.9.0 #433 (dependabot[bot])
* [check-disk] Closes #440 added sort the chec-disk output #441 (hurrycaine)
* Bump github.com/fsouza/go-dockerclient from 1.6.5 to 1.6.6 #434 (dependabot[bot])


## 0.38.0 (2020-12-09)

* Bump github.com/shirou/gopsutil from 2.20.8+incompatible to 2.20.9+incompatible #416 (dependabot-preview[bot])
* migrate CIs to GitHub Actions #432 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.34.32 to 1.35.35 #431 (dependabot[bot])
* Bump golang.org/x/text from 0.3.3 to 0.3.4 #422 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.2.3 to 3.2.4 #420 (dependabot-preview[bot])
* Update Dependabot config file #428 (dependabot-preview[bot])
* [check-postgresql] Add sslrootcert option #425 (nonylene)


## 0.37.1 (2020-10-01)

* Bump github.com/aws/aws-sdk-go from 1.34.22 to 1.34.32 #414 (dependabot-preview[bot])


## 0.37.0 (2020-09-15)

* add arm64 architecture packages #410 (lufia)
* Bump github.com/shirou/gopsutil from 2.20.7+incompatible to 2.20.8+incompatible #408 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.33.17 to 1.34.22 #411 (dependabot-preview[bot])
* [check-log] stabilize time-dependent tests #406 (astj)
* [check-http]adding more options to check-http #403 (fgouteroux)
* [check-load]fix check-load percpu output #405 (fgouteroux)
* Bump github.com/shirou/gopsutil from 2.20.6+incompatible to 2.20.7+incompatible #400 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.33.12 to 1.33.17 #399 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.7.1 to 1.8.0 #398 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.10 to 3.2.3 #390 (dependabot-preview[bot])


## 0.36.0 (2020-07-29)

* [check-aws-sqs-queue-size] Replace GoAMZ with aws-sdk-go #396 (astj)
* Bump github.com/aws/aws-sdk-go from 1.31.11 to 1.33.12 #393 (dependabot-preview[bot])
* Bump github.com/shirou/gopsutil from 2.20.4+incompatible to 2.20.6+incompatible #386 (dependabot-preview[bot])
* Bump github.com/mattn/go-zglob from 0.0.1 to 0.0.3 #394 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.7.0 to 1.7.1 #395 (dependabot-preview[bot])


## 0.35.0 (2020-07-20)

* Bump github.com/lib/pq from 1.6.0 to 1.7.0 #379 (dependabot-preview[bot])
* Bump golang.org/x/text from 0.3.2 to 0.3.3 #383 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 #373 (dependabot-preview[bot])
* [check-ssl-cert] check intermediate- and root-certificates #377 (lufia)
* aws-sdk-go 1.31.11 #372 (astj)
* Bump github.com/aws/aws-sdk-go from 1.30.26 to 1.31.7 #369 (dependabot-preview[bot])
* Bump github.com/lib/pq from 1.5.2 to 1.6.0 #370 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 #371 (dependabot-preview[bot])
* Go 1.14 #352 (ne-sachirou)
* [check-postgresql] add test.sh #365 (lufia)
* Bump github.com/aws/aws-sdk-go from 1.30.9 to 1.30.26 #362 (dependabot-preview[bot])


## 0.34.1 (2020-05-14)

* Bump github.com/lib/pq from 1.3.0 to 1.5.2 #359 (dependabot-preview[bot])
* ignore diffs of go.mod and go.sum #363 (lufia)
* Bump github.com/shirou/gopsutil from 2.20.3+incompatible to 2.20.4+incompatible #358 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.8 to 3.1.10 #360 (dependabot-preview[bot])
* Bump github.com/fsouza/go-dockerclient from 1.6.3 to 1.6.5 #355 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.29.30 to 1.30.9 #351 (dependabot-preview[bot])
* Bump github.com/shirou/gopsutil from 2.20.2+incompatible to 2.20.3+incompatible #347 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.7 to 3.1.8 #346 (dependabot-preview[bot])
* Add documents for testing #348 (lufia)


## 0.34.0 (2020-04-03)

* Bump github.com/aws/aws-sdk-go from 1.29.24 to 1.29.30 #342 (dependabot-preview[bot])
* Bump github.com/beevik/ntp from 0.2.0 to 0.3.0 #340 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.29.14 to 1.29.24 #339 (dependabot-preview[bot])
* [check-redis] Add password support for check-redis #338 (n-rodriguez)
* Bump github.com/aws/aws-sdk-go from 1.28.7 to 1.29.14 #336 (dependabot-preview[bot])
* Bump github.com/shirou/gopsutil from 2.20.1+incompatible to 2.20.2+incompatible #335 (dependabot-preview[bot])
* Bump github.com/fsouza/go-dockerclient from 1.6.0 to 1.6.3 #332 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 #331 (dependabot-preview[bot])
* Bump github.com/go-ldap/ldap/v3 from 3.1.5 to 3.1.7 #329 (dependabot-preview[bot])
* [check-aws-cloudwatch-logs] Added option to specify AWS API retry count #334 (masahide)
* [check-aws-cloudwatch-logs] fix removing of quote (") implicitly from few options #330 (lufia)
* Bump github.com/shirou/gopsutil from 2.19.12+incompatible to 2.20.1+incompatible #323 (dependabot-preview[bot])
* rename: github.com/motemen/gobump -> github.com/x-motemen/gobump #322 (lufia)


## 0.33.1 (2020-01-22)

* Bump github.com/aws/aws-sdk-go from 1.27.0 to 1.28.7 #319 (dependabot-preview[bot])
* [check-aws-cloudwatch-logs] Use "errors" instead of "github.com/pkg/errors" #318 (astj)
* Bump github.com/shirou/gopsutil from 2.19.11+incompatible to 2.19.12+incompatible #312 (dependabot-preview[bot])
* Bump github.com/aws/aws-sdk-go from 1.26.7 to 1.27.0 #313 (dependabot-preview[bot])
* [check-log] When specified multiple exclude pattern, perform search that excludes all conditions. #294 (tukaelu)
* Bump github.com/aws/aws-sdk-go from 1.26.5 to 1.26.7 #310 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.3.0 to 1.4.0 #309 (dependabot-preview[bot])
* add .dependabot/config.yml #307 (lufia)
* refactor Makefile and update dependencies #306 (lufia)


## 0.33.0 (2019-11-21)

* [check-log] Make building the error lines efficiently #304 (ygurumi)
* [check-log] Ignore broken/unexpected json on reading state #302 (astj)


## 0.32.1 (2019-10-24)

* Build with Go 1.12.12


## 0.32.0 (2019-10-02)

* [doc]add repository policy #300 (lufia)
* add --user to check-http #297 (lausser)
* Update modules #298 (ne-sachirou)
* [check-procs] If more than one pattern is specified, find processes that meet any of the conditions. #293 (tukaelu)


## 0.31.1 (2019-07-22)

* Build with Go 1.12 #292 (astj)


## 0.31.0 (2019-06-11)

* Support go modules #289 (astj)


## 0.30.0 (2019-05-08)

* [check-aws-cloudwatch-logs] make handling a state file more safe #286 (lufia)
* [check-log] make handling a state file more safe #285 (lufia)
* [check-ldap] Update go-ldap to v3 #284 (nonylene)


## 0.29.0 (2019-03-27)

* Add check-ping #280 (a-know)
* [ntservice] Enable to specify service name to exclude from match #279 (a-know)
* Add NTP stratum check to check-ntpoffset #276 #278 (EijiSugiura, susisu)
* [check-http] add --proxy option #277 (astj)


## 0.28.0 (2019-02-13)

* Improve READMEs #274 (a-know)
* added support for netbsd on check-load #273 (paulbsd)
* [check-cert-file] improve README #272 (a-know)
*  [check-log][check-windows-eventlog] Improve atomic write #270 (itchyny)
* [check-log]stop reading logs on SIGTERM #268 (lufia)


## 0.27.0 (2019-01-10)

* [check-disk] Do not check inodes% along with disk% #266 (susisu)
* [check-disk] skip checks if there are no inodes #265 (susisu)


## 0.26.0 (2018-11-26)

* [check-http] Support --connect-to option #263 (astj)


## 0.25.0 (2018-11-12)

* [check-redis] add replication subcommand #261 (yuuki)


## 0.24.0 (2018-10-17)

* add User-Agent header to http check plugins #257 (astj)


## 0.23.0 (2018-09-27)

* add aws-cloudwatch-logs #255 (syou6162)
* Add CloudWatch Logs plugin #248 (itchyny)


## 0.22.1 (2018-09-13)

* [check-log] Trace an old file after logrotation with the inode number  #250 (yuuki)
* [check-log] Jsonize status file #252 (yuuki)
* Use Go 1.11 #253 (astj)
* [check-http] Add --max-redirects option #249 (nonylene)


## 0.22.0 (2018-08-30)

* Add check-smtp #243 (shiimaxx)


## 0.21.2 (2018-07-25)

* modify message check-windows-eventlog #241 (daiksy)


## 0.21.1 (2018-06-20)

* [check-http] Set Host header via req.Host #239 (nonylene)


## 0.21.0 (2018-06-07)

* add check-ssl-cert #34 (Songmu)
* Add feature use ntp server #237 (netmarkjp)
* [check-mysql] add unix socket option #236 (sugy)


## 0.20.0 (2018-05-16)

* [check-log] Add option to suppress pattern display #234 (k-hal)
* [check-windows-eventlog] fix README - Some of the listed EVENTTYPEs can not be detected as alerts #233 (a-know)
* [check-procs][check-cert-file] Fix duplicated output of usage with --help argument #231 (itchyny)


## 0.19.0 (2018-03-28)

* add check-ldap #227 (taku-k)
* [check-http] add regexp pattern option #225 (taku-k)


## 0.18.0 (2018-03-15)

* [check-http] add host header option #224 (taku-k)


## 0.17.1 (2018-03-01)

* [check-log]improve skip bytes when file size is zero #222 (hayajo)


## 0.17.0 (2018-02-08)

* [check-procs] fix `getProcs` error handling #216 (mechairoi)
* [ntpoffset] Refine NTP daemon detection and add tests #219 (Songmu)
* [check-tcp] add -W option #212 (dozen)
* [check-procs] Error handling is required #204 (mechairoi)
* Avoid race condition in checkhttp.Run() #214 (astj)
* [check-http] Expose checkhttp.Run due for using check-http as a library #210 (aereal)


## 0.16.0 (2018-01-23)

* Setting password via environment variable #209 (hayajo)
* update rpm-v2 task for building Amazon Linux 2 package #208 (hayajo)


## 0.15.0 (2018-01-10)

* [check-mysql] add readonly subcommand #205 (ichirin2501)
* [uptime] use go-osstat/uptime instead of golib/uptime for getting more accurate uptime #206 (Songmu)


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
