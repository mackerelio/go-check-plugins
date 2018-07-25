%define __targetdir /usr/bin

Name:      mackerel-check-plugins
Version:   %{_version}
Release:   1%{?dist}
License:   ASL 2.0
Summary:   macekrel.io check plugins
URL:       https://mackerel.io
Group:     Applications/System
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

%{__install} -m0755 %{_sourcedir}/build/mackerel-check %{buildroot}%{__targetdir}/

for i in aws-sqs-queue-size cert-file disk elasticsearch file-age file-size http jmx-jolokia ldap load log mailq masterha memcached mysql ntpoffset postgresql procs redis solr ssh ssl-cert tcp uptime; do \
    ln -s ./mackerel-check %{buildroot}%{__targetdir}/check-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

%changelog
* Wed Jul 25 2018 <mackerel-developers@hatena.ne.jp> - 0.21.2
- modify message check-windows-eventlog (by daiksy)

* Wed Jun 20 2018 <mackerel-developers@hatena.ne.jp> - 0.21.1
- [check-http] Set Host header via req.Host (by nonylene)

* Thu Jun 07 2018 <mackerel-developers@hatena.ne.jp> - 0.21.0
- add check-ssl-cert (by Songmu)
- Add feature use ntp server (by netmarkjp)
- [check-mysql] add unix socket option (by sugy)

* Wed May 16 2018 <mackerel-developers@hatena.ne.jp> - 0.20.0
- [check-log] Add option to suppress pattern display (by k-hal)
- [check-windows-eventlog] fix README - Some of the listed EVENTTYPEs can not be detected as alerts (by a-know)
- [check-procs][check-cert-file] Fix duplicated output of usage with --help argument (by itchyny)

* Wed Mar 28 2018 <mackerel-developers@hatena.ne.jp> - 0.19.0
- add check-ldap (by taku-k)
- [check-http] add regexp pattern option (by taku-k)

* Thu Mar 15 2018 <mackerel-developers@hatena.ne.jp> - 0.18.0
- [check-http] add host header option (by taku-k)

* Thu Mar 01 2018 <mackerel-developers@hatena.ne.jp> - 0.17.1
- [check-log]improve skip bytes when file size is zero (by hayajo)

* Thu Feb 08 2018 <mackerel-developers@hatena.ne.jp> - 0.17.0
- [check-procs] fix `getProcs` error handling (by mechairoi)
- [ntpoffset] Refine NTP daemon detection and add tests (by Songmu)
- [check-tcp] add -W option (by dozen)
- [check-procs] Error handling is required (by mechairoi)
- Avoid race condition in checkhttp.Run() (by astj)
- [check-http] Expose checkhttp.Run due for using check-http as a library (by aereal)

* Tue Jan 23 2018 <mackerel-developers@hatena.ne.jp> - 0.16.0
- Setting password via environment variable (by hayajo)
- update rpm-v2 task for building Amazon Linux 2 package (by hayajo)

* Wed Jan 10 2018 <mackerel-developers@hatena.ne.jp> - 0.15.0
- [check-mysql] add readonly subcommand (by ichirin2501)
- [uptime] use go-osstat/uptime instead of golib/uptime for getting more accurate uptime (by Songmu)

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
