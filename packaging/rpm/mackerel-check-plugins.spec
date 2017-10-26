%define _binaries_in_noarch_packages_terminate_build   0
%define _localbindir /usr/local/bin
%define __targetdir /usr/bin
%define __oldtargetdir /usr/local/bin

Name:      mackerel-check-plugins
Version:   %{_version}
Release:   1
License:   Commercial
Summary:   macekrel.io check plugins
URL:       https://mackerel.io
Group:     Hatena
Packager:  Hatena
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
mackerel.io check plugins

%prep

%build

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

for i in aws-sqs-queue-size cert-file disk elasticsearch file-age file-size http jmx-jolokia load log mailq masterha memcached mysql ntpoffset postgresql procs redis solr ssh tcp uptime; do \
    %{__install} -m0755 %{_sourcedir}/build/check-$i %{buildroot}%{__targetdir}/; \
done

# put symlinks to /usr/local/bin for backward compatibility of following plugins
%{__install} -d -m755 %{buildroot}%{__oldtargetdir}
for i in elasticsearch file-age file-size http jmx-jolokia load log mailq memcached mysql ntpoffset postgresql procs redis solr tcp uptime; \
do
    ln -s ../../bin/check-$i %{buildroot}%{__oldtargetdir}/check-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*
%{__oldtargetdir}/*

%changelog
* Thu Oct 26 2017 <mackerel-developers@hatena.ne.jp> - 0.14.1
- fix check-disk options: -x, -X, -p, -N (by plaster)

* Thu Oct 12 2017 <mackerel-developers@hatena.ne.jp> - 0.14.0
- [check-log]Show matched filenames when we use `--return` option (by syou6162)

* Wed Sep 27 2017 <mackerel-developers@hatena.ne.jp> - 0.13.0
- build with Go 1.9 (by astj)
- [check-log] Add check-first option and skip the log file on the first run on default (by edangelion)

* Wed Aug 23 2017 <mackerel-developers@hatena.ne.jp> - 0.12.0
- add check-disk to package (by astj)
- add check-disk (by edangelion)
- [check-postgresql] Add dbname to postgresql-setting (by edangelion)

* Wed Aug 02 2017 <mackerel-developers@hatena.ne.jp> - 0.11.1
- Remove check-ssh binary (by astj)

* Wed Jul 26 2017 <mackerel-developers@hatena.ne.jp> - 0.11.0
- [check-http] Add -i flag to specify source IP (by mattn)
- [check-http] Add -s option to map HTTP status (by mattn)

* Wed Jun 07 2017 <mackerel-developers@hatena.ne.jp> - 0.10.4
- v2 packages (rpm and deb) (by Songmu)
- [check-log]  When specified multiple pattern, perform search that satisfies all conditions (by a-know)

* Tue May 16 2017 <mackerel-developers@hatena.ne.jp> - 0.10.3
- [ntpoffset] support chronyd (by Songmu)
- [check-ssh] fix the problem that check-ssh cannot invoke SSH connection (by astj)

* Mon May 15 2017 <mackerel-developers@hatena.ne.jp> - 0.10.2
- [experimental] update release scripts (by Songmu)

* Thu Apr 27 2017 <mackerel-developers@hatena.ne.jp> - 0.10.1-1
- use wmi query instead of running wmic command (by mattn)
- Use golib/pluginutil.PluginWorkDir() (by astj)

* Thu Apr 06 2017 <mackerel-developers@hatena.ne.jp> - 0.10.0-1
- bump go to 1.8 (by astj)

* Mon Mar 27 2017 <mackerel-developers@hatena.ne.jp> - 0.9.7-1
- check lower-case driver letters (by mattn)

* Wed Mar 22 2017 <mackerel-developers@hatena.ne.jp> - 0.9.6-1
- Change directory structure convention of each plugin (by Songmu)
- run tests under ./check-XXX/lib (by astj)
- fix test for AppVayor (by daiksy)

* Thu Mar 09 2017 <mackerel-developers@hatena.ne.jp> - 0.9.5-1
- add appveyor.yml and fix failing tests on windows (by Songmu)
- [check-tcp] connect timeout (by Songmu)
- [check-tcp] [bugfix] fix decimal threshold value handling (by Songmu)

* Wed Feb 22 2017 <mackerel-developers@hatena.ne.jp> - 0.9.4-1
- ensure close temporary file in writeFileAtomically (by itchyny)
- Write the file in same directory (by mattn)

* Wed Feb 08 2017 <mackerel-developers@hatena.ne.jp> - 0.9.3-1
- fix matching for 'Audit Success' and 'Audit Failure' (by mattn)

* Wed Jan 25 2017 <mackerel-developers@hatena.ne.jp> - 0.9.2-1
- [check-windows-eventlog] add --source-exclude, --message-exclude and --event-id (by mattn, daiksy)
- [check-windows-eventlog] remove --event-id and add --event-id-pattern, --event-id-exclude (by mattn)

* Wed Jan 11 2017 <mackerel-developers@hatena.ne.jp> - 0.9.1-1
- [check-log] support glob in --file option (by astj)

* Wed Jan 04 2017 <mackerel-developers@hatena.ne.jp> - 0.9.0-1
- add check-windows-eventlog (by daiksy)
- [check-log]fix encoding option (by daiksy)

* Tue Nov 29 2016 <mackerel-developers@hatena.ne.jp> - 0.8.1-1
- Fix state in check-procs (by itchyny)
- Fix the links to the document (by itchyny)
- remove checking Ignore (by mattn)
- [check-log] fix state file path (fix wrong slice copy) (by Songmu)

* Thu Oct 27 2016 <mackerel-developers@hatena.ne.jp> - 0.8.0-1
- [check-log] improve Windows support (by daiksy)

* Tue Oct 18 2016 <mackerel-developers@hatena.ne.jp> - 0.7.0-1
- Add option for skip searching logfile content if logfile does not exist (by a-know)
- [check-log] write file atomically when saving read position into state file (by Songmu)

* Wed Sep 07 2016 <mackerel-developers@hatena.ne.jp> - 0.6.3-1
- fix check-mysql replication to detect IO thread 'Connecting' (by hiroakis)
- [file-age] Remove unnecessary newline (by b4b4r07)

* Thu Jun 23 2016 <mackerel-developers@hatena.ne.jp> - 0.6.2-1
- Fixed argument parser error: (by karupanerura)

* Fri May 13 2016 <mackerel-developers@hatena.ne.jp> - 0.6.1-1
- no panic check-masterha when not found the targets, and refactoring (by karupanerura)

* Tue May 10 2016 <mackerel-developers@hatena.ne.jp> - 0.6.0-1
- supported gearman ascii protocol for check-tcp (by karupanerura)
- added check-masterha command to check masterha status (by karupanerura)
- Plugin for checking AWS SQS queue size (by hiroakis)
- fix: rpm should not include dir (by stanaka)
- added ssh checker (by karupanerura)
- remove 'golang.org/x/tools/cmd/vet' (by y-kuno)
- [uptime/procs] `--warn-over/under` is deprecated. use `--warning-over/under` instead (by Songmu)
- add aws-sqs-queue-size, cert-file, masterha and ssh into package (by Songmu)
- bump up go version to 1.6.2 (by stanaka)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.5.2
- Revert "use /usr/bin/check-*" (by Songmu)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.5.1
- use /usr/bin/check-* (by naokibtn)
- use GOARCH=amd64 for now (by Songmu)

* Wed Mar 02 2016 <y.songmu@gmail.com> - 0.5.0
- add check-solr (by supercaracal)
-  add check-jmx-jolokia (by y-kuno)
- check-memcached (by naokibtn)
- Add scheme option to check-elasticsearch (by yano3)
- Check file size (by hiroakis)
- Add uptime sub command to check-mysql (by kazeburo)
- add check-cert-file (by kazeburo)
- [tcp] Add --no-check-certificate (by kazeburo)
- Fixed slurp. conn.read returns content with EOF (by kazeburo)
- Fix check-tcp IPv6 testcase on OSX(?) (by hanazuki)
- check-redis: Report critical if master_link_status is down (by hanazuki)
- check-redis: Fixed panic (by yoheimuta)
- check-procs: Fixed the counting logic with -p (by yoheimuta)
- add check-uptime (by Songmu)
- add file-size, jmx-jolokia, memcached, solr, uptime into package config (by Songmu)

* Thu Feb 04 2016 <y.songmu@gmail.com> - 0.4.0
- Fix duplicated help message (by hfm)
- add qmail queue check to check-mailq (by tnmt)
- Add check-elasticsearch (by naokibtn)
- Add check-redis (by naokibtn)
- check-procs: check PPID (by hfm)

* Thu Jan 07 2016 <y.songmu@gmail.com> - 0.3.1
- build with go1.5 (by Songmu)

* Wed Jan 06 2016 <y.songmu@gmail.com> - 0.3.0
- add check-postgresql (by supercaracal)
- [check-ntpoffset] work on ntp 4.2.2 (by naokibtn)
- check-file-age: show time in message (by naokibtn)
- add --no-check-certificate options to check_http (by cl-lab-k)
- add check-mailq, currently only available for postfix. (by tnmt)
- add check-mailq and check-postgresql into package (by Songmu)

* Tue Dec 08 2015 <y.songmu@gmail.com> - 0.2.2
- A plugin to check ntpoffset (by hiroakis)
- Check tcp unix domain socket (by tkuchiki)
- [check-tcp] add ipv6 support (by Songmu)

* Wed Nov 25 2015 <y.songmu@gmail.com> - 0.2.1
- Fix bugs of check-log (by Songmu)
- [check-log] add --no-state option (by Songmu)

* Fri Nov 20 2015 <y.songmu@gmail.com> - 0.2.0
- [check-procs] support `--critical-over=0` and `--warn-over=0` (by Songmu)
- add check-tcp (by Songmu)

* Thu Nov 12 2015 <y.songmu@gmail.com> - 0.1.1
- check-procs for windows (by mechairoi)
- [bug] [check-log] fix large file handling (by matsuzj)

* Thu Nov 05 2015 <y.songmu@gmail.com> - 0.1.0
- check-log (by Songmu)
- Add check-log in the packages (by Songmu)

* Mon Oct 26 2015 <daiksy@hatena.ne.jp> - 0.0.5
- Add mysql in packages

* Mon Oct 26 2015 <daiksy@hatena.ne.jp> - 0.0.4
- Refactor Mysql checks (by Songmu)
- Add Mysql checks (by hiroakis)

* Thu Oct 15 2015 <itchyny@hatena.ne.jp> - 0.0.3
- reduce binary size by using ldflags (by Songmu)
- Remove cgo dependency from check-load (by Songmu)
- Add check-load in the packages

* Thu Oct 08 2015 <itchyny@hatena.ne.jp> - 0.0.2
- fix MatchSelf behaviour (by Songmu)

* Wed Oct 07 2015 <itchyny@hatena.ne.jp> - 0.0.1
- Fix release tools

* Wed Oct 07 2015 <itchyny@hatena.ne.jp> - 0.0.0
- Initial release
