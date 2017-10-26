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

for i in aws-sqs-queue-size cert-file disk elasticsearch file-age file-size http jmx-jolokia load log mailq masterha memcached mysql ntpoffset postgresql procs redis solr ssh tcp uptime; do \
    ln -s ./mackerel-check %{buildroot}%{__targetdir}/check-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

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
