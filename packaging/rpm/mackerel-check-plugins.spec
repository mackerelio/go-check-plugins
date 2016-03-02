%define _binaries_in_noarch_packages_terminate_build   0
%define _localbindir /usr/local/bin
%define __targetdir /usr/local/bin

Name:      mackerel-check-plugins
Version:   0.5.0
Release:   1
License:   Commercial
Summary:   macekrel.io check plugins
URL:       https://mackerel.io
Group:     Hatena
Packager:  Hatena
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
mackerel.io check plugins

%prep

%build

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

for i in elasticsearch file-age http load log mailq ntpoffset postgresql procs redis tcp mysql;do \
    %{__install} -m0755 %{_sourcedir}/build/check-$i %{buildroot}%{__targetdir}/; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}

%changelog
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
