%define __targetdir /usr/bin

Name:      mackerel-check-plugins
Version:   %{_version}
Release:   1%{?dist}
License:   ASL 2.0
Summary:   macekrel.io check plugins
URL:       https://mackerel.io
Group:     Applications/System
Packager:  Hatena
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
mackerel.io check plugins

%prep

%build

%install
%{__rm} -rf %{buildroot}

%{__mkdir} -p %{buildroot}%{__targetdir}

%{__install} -m0755 %{_sourcedir}/build/mackerel-check %{buildroot}%{__targetdir}/

for i in aws-cloudwatch-logs aws-sqs-queue-size cert-file disk elasticsearch file-age file-size http jmx-jolokia ldap load log mailq masterha memcached mysql ntpoffset ping postgresql procs redis smtp solr ssh ssl-cert tcp uptime; do \
    ln -s ./mackerel-check %{buildroot}%{__targetdir}/check-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

%changelog
* Tue Jul 06 2021 <mackerel-developers@hatena.ne.jp> - 0.39.5
- Bump github.com/shirou/gopsutil/v3 from 3.21.5 to 3.21.6 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.38.45 to 1.38.68 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.2 to 1.7.3 (by dependabot[bot])

* Wed Jun 23 2021 <mackerel-developers@hatena.ne.jp> - 0.39.4
- [ci]rewrite check-memcached tests. used docker. (by yseto)
- refactor check-log tests. (by yseto)
- [ci] run tests on the workflow (by lufia)
- [check-disk] upgrade gopsutil to v3 and fix treatment for fstype=none (by susisu)

* Thu Jun 03 2021 <mackerel-developers@hatena.ne.jp> - 0.39.3
- Bump github.com/aws/aws-sdk-go from 1.38.40 to 1.38.45 (by dependabot[bot])
- Bump github.com/go-sql-driver/mysql from 1.5.0 to 1.6.0 (by dependabot[bot])
- Bump github.com/jmoiron/sqlx from 1.3.1 to 1.3.4 (by dependabot[bot])
- Bump github.com/jessevdk/go-flags from 1.4.0 to 1.5.0 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.0 to 1.10.2 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.5 to 0.3.6 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.37.30 to 1.38.40 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.2.4 to 3.3.0 (by dependabot[bot])
- upgrade go1.14 -> 1.16 (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.37.1 to 1.37.30 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.0 to 1.7.2 (by dependabot[bot])
- Bump github.com/lib/pq from 1.9.0 to 1.10.0 (by dependabot[bot])
- [check-mysql] Use go-mysql-driver and sqlx instead of mymysql (by nonylene)
- [ci] replace token (by yseto)
- [ci] replace mackerel-github-release (by yseto)

* Wed Feb 03 2021 <mackerel-developers@hatena.ne.jp> - 0.39.2
- Bump github.com/aws/aws-sdk-go from 1.36.31 to 1.37.1 (by dependabot[bot])
- Closes #455 check-load fix null pointer issue when critical (by hurrycaine)
- Bump github.com/aws/aws-sdk-go from 1.36.28 to 1.36.31 (by dependabot[bot])

* Thu Jan 21 2021 <mackerel-developers@hatena.ne.jp> - 0.39.1
- Bump github.com/aws/aws-sdk-go from 1.36.19 to 1.36.28 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.6.6 to 1.7.0 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (by dependabot[bot])
- Bump github.com/go-ole/go-ole from 1.2.4 to 1.2.5 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.4 to 0.3.5 (by dependabot[bot])

* Thu Jan 14 2021 <mackerel-developers@hatena.ne.jp> - 0.39.0
- Bump github.com/aws/aws-sdk-go from 1.35.35 to 1.36.19 (by dependabot[bot])
- Bump github.com/lib/pq from 1.8.0 to 1.9.0 (by dependabot[bot])
- [check-disk] Closes #440 added sort the chec-disk output (by hurrycaine)
- Bump github.com/fsouza/go-dockerclient from 1.6.5 to 1.6.6 (by dependabot[bot])

* Wed Dec 09 2020 <mackerel-developers@hatena.ne.jp> - 0.38.0
- Bump github.com/shirou/gopsutil from 2.20.8+incompatible to 2.20.9+incompatible (by dependabot-preview[bot])
- migrate CIs to GitHub Actions (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.34.32 to 1.35.35 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.3 to 0.3.4 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.2.3 to 3.2.4 (by dependabot-preview[bot])
- Update Dependabot config file (by dependabot-preview[bot])
- [check-postgresql] Add sslrootcert option (by nonylene)

* Thu Oct 01 2020 <mackerel-developers@hatena.ne.jp> - 0.37.1
- Bump github.com/aws/aws-sdk-go from 1.34.22 to 1.34.32 (by dependabot-preview[bot])

* Tue Sep 15 2020 <mackerel-developers@hatena.ne.jp> - 0.37.0
- add arm64 architecture packages (by lufia)
- Bump github.com/shirou/gopsutil from 2.20.7+incompatible to 2.20.8+incompatible (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.33.17 to 1.34.22 (by dependabot-preview[bot])
- [check-log] stabilize time-dependent tests (by astj)
- [check-http]adding more options to check-http (by fgouteroux)
- [check-load]fix check-load percpu output (by fgouteroux)
- Bump github.com/shirou/gopsutil from 2.20.6+incompatible to 2.20.7+incompatible (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.33.12 to 1.33.17 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.7.1 to 1.8.0 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.10 to 3.2.3 (by dependabot-preview[bot])

* Wed Jul 29 2020 <mackerel-developers@hatena.ne.jp> - 0.36.0
- [check-aws-sqs-queue-size] Replace GoAMZ with aws-sdk-go (by astj)
- Bump github.com/aws/aws-sdk-go from 1.31.11 to 1.33.12 (by dependabot-preview[bot])
- Bump github.com/shirou/gopsutil from 2.20.4+incompatible to 2.20.6+incompatible (by dependabot-preview[bot])
- Bump github.com/mattn/go-zglob from 0.0.1 to 0.0.3 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.7.0 to 1.7.1 (by dependabot-preview[bot])

* Mon Jul 20 2020 <mackerel-developers@hatena.ne.jp> - 0.35.0
- Bump github.com/lib/pq from 1.6.0 to 1.7.0 (by dependabot-preview[bot])
- Bump golang.org/x/text from 0.3.2 to 0.3.3 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 (by dependabot-preview[bot])
- [check-ssl-cert] check intermediate- and root-certificates (by lufia)
- aws-sdk-go 1.31.11 (by astj)
- Bump github.com/aws/aws-sdk-go from 1.30.26 to 1.31.7 (by dependabot-preview[bot])
- Bump github.com/lib/pq from 1.5.2 to 1.6.0 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 (by dependabot-preview[bot])
- Go 1.14 (by ne-sachirou)
- [check-postgresql] add test.sh (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.30.9 to 1.30.26 (by dependabot-preview[bot])

* Thu May 14 2020 <mackerel-developers@hatena.ne.jp> - 0.34.1
- Bump github.com/lib/pq from 1.3.0 to 1.5.2 (by dependabot-preview[bot])
- ignore diffs of go.mod and go.sum (by lufia)
- Bump github.com/shirou/gopsutil from 2.20.3+incompatible to 2.20.4+incompatible (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.8 to 3.1.10 (by dependabot-preview[bot])
- Bump github.com/fsouza/go-dockerclient from 1.6.3 to 1.6.5 (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.29.30 to 1.30.9 (by dependabot-preview[bot])
- Bump github.com/shirou/gopsutil from 2.20.2+incompatible to 2.20.3+incompatible (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.7 to 3.1.8 (by dependabot-preview[bot])
- Add documents for testing (by lufia)

* Fri Apr 03 2020 <mackerel-developers@hatena.ne.jp> - 0.34.0
- Bump github.com/aws/aws-sdk-go from 1.29.24 to 1.29.30 (by dependabot-preview[bot])
- Bump github.com/beevik/ntp from 0.2.0 to 0.3.0 (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.29.14 to 1.29.24 (by dependabot-preview[bot])
- [check-redis] Add password support for check-redis (by n-rodriguez)
- Bump github.com/aws/aws-sdk-go from 1.28.7 to 1.29.14 (by dependabot-preview[bot])
- Bump github.com/shirou/gopsutil from 2.20.1+incompatible to 2.20.2+incompatible (by dependabot-preview[bot])
- Bump github.com/fsouza/go-dockerclient from 1.6.0 to 1.6.3 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 (by dependabot-preview[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.1.5 to 3.1.7 (by dependabot-preview[bot])
- [check-aws-cloudwatch-logs] Added option to specify AWS API retry count (by masahide)
- [check-aws-cloudwatch-logs] fix removing of quote (") implicitly from few options (by lufia)
- Bump github.com/shirou/gopsutil from 2.19.12+incompatible to 2.20.1+incompatible (by dependabot-preview[bot])
- rename: github.com/motemen/gobump -> github.com/x-motemen/gobump (by lufia)

* Wed Jan 22 2020 <mackerel-developers@hatena.ne.jp> - 0.33.1
- Bump github.com/aws/aws-sdk-go from 1.27.0 to 1.28.7 (by dependabot-preview[bot])
- [check-aws-cloudwatch-logs] Use "errors" instead of "github.com/pkg/errors" (by astj)
- Bump github.com/shirou/gopsutil from 2.19.11+incompatible to 2.19.12+incompatible (by dependabot-preview[bot])
- Bump github.com/aws/aws-sdk-go from 1.26.7 to 1.27.0 (by dependabot-preview[bot])
- [check-log] When specified multiple exclude pattern, perform search that excludes all conditions. (by tukaelu)
- Bump github.com/aws/aws-sdk-go from 1.26.5 to 1.26.7 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.3.0 to 1.4.0 (by dependabot-preview[bot])
- add .dependabot/config.yml (by lufia)
- refactor Makefile and update dependencies (by lufia)

* Thu Nov 21 2019 <mackerel-developers@hatena.ne.jp> - 0.33.0
- [check-log] Make building the error lines efficiently (by ygurumi)
- [check-log] Ignore broken/unexpected json on reading state (by astj)

* Thu Oct 24 2019 <mackerel-developers@hatena.ne.jp> - 0.32.1
- Build with Go 1.12.12

* Wed Oct 02 2019 <mackerel-developers@hatena.ne.jp> - 0.32.0
- [doc]add repository policy (by lufia)
- add --user to check-http (by lausser)
- Update modules (by ne-sachirou)
- [check-procs] If more than one pattern is specified, find processes that meet any of the conditions. (by tukaelu)

* Mon Jul 22 2019 <mackerel-developers@hatena.ne.jp> - 0.31.1
- Build with Go 1.12 (by astj)

* Tue Jun 11 2019 <mackerel-developers@hatena.ne.jp> - 0.31.0
- Support go modules (by astj)

* Wed May 08 2019 <mackerel-developers@hatena.ne.jp> - 0.30.0
- [check-aws-cloudwatch-logs] make handling a state file more safe (by lufia)
- [check-log] make handling a state file more safe (by lufia)
- [check-ldap] Update go-ldap to v3 (by nonylene)

* Wed Mar 27 2019 <mackerel-developers@hatena.ne.jp> - 0.29.0
- Add check-ping (by a-know)
- [ntservice] Enable to specify service name to exclude from match (by a-know)
- Add NTP stratum check to check-ntpoffset (by EijiSugiura, susisu)
- [check-http] add --proxy option (by astj)

* Wed Feb 13 2019 <mackerel-developers@hatena.ne.jp> - 0.28.0
- Improve READMEs (by a-know)
- added support for netbsd on check-load (by paulbsd)
- [check-cert-file] improve README (by a-know)
-  [check-log][check-windows-eventlog] Improve atomic write (by itchyny)
- [check-log]stop reading logs on SIGTERM (by lufia)

* Thu Jan 10 2019 <mackerel-developers@hatena.ne.jp> - 0.27.0
- [check-disk] Do not check inodes% along with disk% (by susisu)
- [check-disk] skip checks if there are no inodes (by susisu)

* Mon Nov 26 2018 <mackerel-developers@hatena.ne.jp> - 0.26.0
- [check-http] Support --connect-to option (by astj)

* Mon Nov 12 2018 <mackerel-developers@hatena.ne.jp> - 0.25.0
- [check-redis] add replication subcommand (by yuuki)

* Wed Oct 17 2018 <mackerel-developers@hatena.ne.jp> - 0.24.0
- add User-Agent header to http check plugins (by astj)

* Thu Sep 27 2018 <mackerel-developers@hatena.ne.jp> - 0.23.0
- add aws-cloudwatch-logs (by syou6162)
- Add CloudWatch Logs plugin (by itchyny)

* Thu Sep 13 2018 <mackerel-developers@hatena.ne.jp> - 0.22.1
- [check-log] Trace an old file after logrotation with the inode number  (by yuuki)
- [check-log] Jsonize status file (by yuuki)
- Use Go 1.11 (by astj)
- [check-http] Add --max-redirects option (by nonylene)

* Thu Aug 30 2018 <mackerel-developers@hatena.ne.jp> - 0.22.0
- Add check-smtp (by shiimaxx)

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
