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

for i in aws-cloudwatch-logs aws-sqs-queue-size cert-file disk dns elasticsearch file-age file-size http jmx-jolokia ldap load log mailq masterha memcached mysql ntpoffset ping postgresql procs redis smtp solr ssh ssl-cert tcp uptime; do \
    ln -s ./mackerel-check %{buildroot}%{__targetdir}/check-$i; \
done

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-, root, root, 0755)
%{__targetdir}/*

%changelog
* Wed Dec 10 2025 <mackerel-developers@hatena.ne.jp> - 0.51.0
- Add check-redis --username flag (by fujiwara)
- Bump the aws-aws-sdk-go-v2 group with 5 updates (by dependabot[bot])
- Bump golang.org/x/sys from 0.38.0 to 0.39.0 in the golang-x group (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/setup-go-matrix.yml from 1.7.0 to 1.8.0 (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/go-lint.yml from 1.7.0 to 1.8.0 (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/go-test.yml from 1.7.0 to 1.8.0 (by dependabot[bot])
- Bump actions/checkout from 6.0.0 to 6.0.1 (by dependabot[bot])
- Bump the aws-aws-sdk-go-v2 group with 5 updates (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/go-lint.yml from 1.6.0 to 1.7.0 (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/setup-go-matrix.yml from 1.6.0 to 1.7.0 (by dependabot[bot])
- Bump mackerelio/workflows/.github/workflows/go-test.yml from 1.6.0 to 1.7.0 (by dependabot[bot])
- Bump actions/checkout from 5.0.0 to 6.0.0 (by dependabot[bot])
- Bump actions/setup-go from 6.0.0 to 6.1.0 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.12.2 to 1.12.3 (by dependabot[bot])
- Bump golang.org/x/crypto from 0.42.0 to 0.45.0 (by dependabot[bot])
- fix buildtag (by yseto)
- update CI (by yseto)
- Bump github.com/gomodule/redigo from 1.9.2 to 1.9.3 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.4.8 to 3.4.12 (by dependabot[bot])
- Bump github.com/beevik/ntp from 1.4.3 to 1.5.0 (by dependabot[bot])
- Bump the aws-aws-sdk-go-v2 group with 5 updates (by dependabot[bot])
- Bump github.com/go-sql-driver/mysql from 1.8.1 to 1.9.3 (by dependabot[bot])

* Mon Oct 20 2025 <mackerel-developers@hatena.ne.jp> - 0.50.1
- added dependabot cooldown (by yseto)
- Bump github.com/mackerelio/go-osstat from 0.2.5 to 0.2.6 in the mackerelio group (by dependabot[bot])
- Bump actions/setup-go from 5 to 6 (by dependabot[bot])
- Bump actions/checkout from 4 to 5 (by dependabot[bot])
- implement status-as option to check-ping (by kga)
- Bump actions/download-artifact from 4 to 5 (by dependabot[bot])
- Bump github.com/miekg/dns from 1.1.50 to 1.1.68 (by dependabot[bot])
- Bump mackerelio/workflows from 1.4.0 to 1.5.0 (by dependabot[bot])
- Bump the golang-x group across 1 directory with 3 updates (by dependabot[bot])
- Bump the testlibs group across 1 directory with 2 updates (by dependabot[bot])
- Bump github.com/mattn/go-zglob from 0.0.4 to 0.0.6 (by dependabot[bot])
- Bump github.com/jessevdk/go-flags from 1.5.0 to 1.6.1 (by dependabot[bot])
- Bump github.com/beevik/ntp from 1.3.1 to 1.4.3 (by dependabot[bot])

* Fri Sep 19 2025 <mackerel-developers@hatena.ne.jp> - 0.50.0
- Fix error handling in check-disk to avoid unnecessary failures when using --path option (by mechairoi)
- Bump github.com/fsouza/go-dockerclient from 1.11.0 to 1.12.2 (by dependabot[bot])
- Bump golang.org/x/net from 0.36.0 to 0.38.0 (by dependabot[bot])

* Fri May 16 2025 <mackerel-developers@hatena.ne.jp> - 0.49.0
- use Go 1.24 (by yseto)
- Release version 0.49.0 (by mackerelbot)
- Remove rewrite some files on every releases (by yseto)
- introduce status-as option to check-ntservice (by kmuto)
- replace to aws-sdk-go-v2 (by yseto)

* Mon Mar 31 2025 <mackerel-developers@hatena.ne.jp> - 0.48.0
- [check-windows-eventlog] add status-as option (by masarasi)
- [check-windows-eventlog] add target event type (by masarasi)
- replace to newer runner-images (by yseto)
- Bump golang.org/x/net from 0.25.0 to 0.36.0 (by dependabot[bot])
- Bump mackerelio/workflows from 1.3.0 to 1.4.0 (by dependabot[bot])
- Bump mackerelio/workflows from 1.2.0 to 1.3.0 (by dependabot[bot])
- use mackerelio/workflows@v1.2.0 (by yseto)

* Wed Jun 12 2024 <mackerel-developers@hatena.ne.jp> - 0.47.0
- return CRITICAL instead of UNKNOWN when check-redis reachable is failed (by kmuto)
- Bump the golang-x group with 3 updates (by dependabot[bot])
- use go 1.22.x on build phase (by lufia)
- update dependencies (by lufia)
- [check-mailq] fix pattern (by lufia)
- Bump github.com/docker/docker from 25.0.4+incompatible to 25.0.5+incompatible (by dependabot[bot])

* Tue Apr 23 2024 <mackerel-developers@hatena.ne.jp> - 0.46.3
- Revert "Bump github.com/miekg/dns from 1.1.50 to 1.1.59" (by ne-sachirou)
- Bump github.com/go-ldap/ldap/v3 from 3.4.4 to 3.4.8 (by dependabot[bot])
- Bump github.com/miekg/dns from 1.1.50 to 1.1.59 (by dependabot[bot])
- Bump golang.org/x/net from 0.17.0 to 0.23.0 (by dependabot[bot])
- Fix: check-log panic with invalid memory address or nil pointer dereference (by ne-sachirou)
- Bump the golang-x group with 2 updates (by dependabot[bot])
- Bump github.com/docker/docker from 23.0.0+incompatible to 24.0.9+incompatible (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.9.4 to 1.11.0 (by dependabot[bot])
- Bump mackerelio/workflows from 1.0.2 to 1.1.0 (by dependabot[bot])
- Bump the testlibs group with 1 update (by dependabot[bot])
- Bump github.com/opencontainers/runc from 1.1.2 to 1.1.12 (by dependabot[bot])
- Bump actions/cache from 3 to 4 (by dependabot[bot])
- Bump github.com/beevik/ntp from 0.3.0 to 1.3.1 (by dependabot[bot])
- Bump github.com/containerd/containerd from 1.6.18 to 1.6.26 (by dependabot[bot])
- Bump actions/upload-artifact from 3 to 4 (by dependabot[bot])
- Bump actions/download-artifact from 3 to 4 (by dependabot[bot])
- Bump actions/setup-go from 4 to 5 (by dependabot[bot])
- Bump github.com/go-ole/go-ole from 1.2.6 to 1.3.0 (by dependabot[bot])

* Tue Feb 27 2024 <mackerel-developers@hatena.ne.jp> - 0.46.2
- Reduce check-log errors when a file in the log directory has been removed at the moment of running check-log (by ne-sachirou)
- Bump the golang-x group with 3 updates (by dependabot[bot])
- Fix path (by yohfee)

* Wed Nov 15 2023 <mackerel-developers@hatena.ne.jp> - 0.46.1
- CGO_ENABLED=0 when build for packaging (by Arthur1)

* Mon Nov 13 2023 <mackerel-developers@hatena.ne.jp> - 0.46.0
- Bump github.com/aws/aws-sdk-go from 1.47.3 to 1.47.9 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.23.8 to 3.23.10 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.271 to 1.47.3 (by dependabot[bot])
- bump go version on build (by yseto)
- fix the line handling of UTF-16le (by kmuto)
- use mackerelio/workflows and upgrade maximum Go version to 1.21 (by lufia)
- check-ntservice: exclude should be -x, not -E (by kmuto)
- Bump github.com/lib/pq from 1.10.7 to 1.10.9 (by dependabot[bot])
- Bump github.com/go-sql-driver/mysql from 1.7.0 to 1.7.1 (by dependabot[bot])
- Bump actions/checkout from 3 to 4 (by dependabot[bot])
- support overwrite status (by yseto)

* Fri Sep 22 2023 <mackerel-developers@hatena.ne.jp> - 0.45.0
- Bump golang.org/x/crypto from 0.6.0 to 0.13.0 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.23.1 to 3.23.8 (by dependabot[bot])
- [check-tcp] Supports option to monitor that ports are closed. (by tukaelu)
- Improve ntservice (by tukaelu)
- Remove old rpm packaging (by yseto)
- Bump github.com/aws/aws-sdk-go from 1.44.199 to 1.44.271 (by dependabot[bot])
- Bump actions/setup-go from 3 to 4 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.3 to 0.2.4 (by dependabot[bot])

* Thu Jul 13 2023 <mackerel-developers@hatena.ne.jp> - 0.44.1
- added build tests. (by yseto)

* Mon Feb 27 2023 <mackerel-developers@hatena.ne.jp> - 0.44.0
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (by dependabot[bot])
- fix gosimple, ineffassign (by wafuwafu13)
- Bump github.com/containerd/containerd from 1.6.14 to 1.6.18 (by dependabot[bot])
- check-dns: add `expected-string` option (by wafuwafu13)
- Bump github.com/aws/aws-sdk-go from 1.44.189 to 1.44.199 (by dependabot[bot])
- Bump golang.org/x/crypto from 0.5.0 to 0.6.0 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.9.3 to 1.9.4 (by dependabot[bot])
- Bump golang.org/x/text from 0.6.0 to 0.7.0 (by dependabot[bot])
- Bump golang.org/x/sys from 0.4.0 to 0.5.0 (by dependabot[bot])
- added dns plugin on package (by yseto)
- Remove `circle.yml` (by wafuwafu13)
- Bump github.com/shirou/gopsutil/v3 from 3.22.12 to 3.23.1 (by dependabot[bot])
- Add check-dns plugin (by wafuwafu13)

* Wed Feb 1 2023 <mackerel-developers@hatena.ne.jp> - 0.43.0
- fix generate docs (by yseto)
- Bump actions/checkout from 2 to 3 (by dependabot[bot])
- Bump actions/setup-go from 2 to 3 (by dependabot[bot])
- Bump actions/cache from 2 to 3 (by dependabot[bot])
- Bump actions/upload-artifact from 2 to 3 (by dependabot[bot])
- Bump actions/download-artifact from 2 to 3 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.157 to 1.44.189 (by dependabot[bot])
- Enables Dependabot version updates for GitHub Actions (by Arthur1)
- Stop build for apt v1 (by Arthur1)
- Bump github.com/fsouza/go-dockerclient from 1.9.0 to 1.9.3 (by dependabot[bot])
- [check-http] add test.sh (by lufia)
- check-ssl-cert: add `ca-file`, `cert-file`, `key-file`, `no-check-certificate` options (by wafuwafu13)
- Bump golang.org/x/text from 0.5.0 to 0.6.0 (by dependabot[bot])
- Bump golang.org/x/crypto from 0.4.0 to 0.5.0 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.22.11 to 3.22.12 (by dependabot[bot])
- Bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0 (by dependabot[bot])

* Wed Jan 18 2023 <mackerel-developers@hatena.ne.jp> - 0.42.4
- check-cert-file: add test (by wafuwafu13)
- test: use `T.TempDir` to create temporary test directory (by Juneezee)
- combine lint, lint-windows, fix test on windows. (by yseto)
- added compile option, fix packaging format (by yseto)
- Update dependencies (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.44.116 to 1.44.157 (by dependabot[bot])

* Thu Oct 20 2022 <mackerel-developers@hatena.ne.jp> - 0.42.3
- Bump golang.org/x/text from 0.3.7 to 0.4.0 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.44.56 to 1.44.116 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.8.3 to 1.9.0 (by dependabot[bot])
- use Go 1.19 on build (by yseto)
- Bump github.com/shirou/gopsutil/v3 from 3.22.2 to 3.22.9 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.6 to 1.10.7 (by dependabot[bot])
- Bump github.com/mackerelio/checkers from 0.0.3 to 0.0.4 (by dependabot[bot])
- [uptime] rewite to testable and add test (by wafuwafu13)
- go.mod from 1.16 to 1.18 (by yseto)
- added test check-file-age (by yseto)
- added test check-file-size (by yseto)
- Bump github.com/mackerelio/go-osstat from 0.2.2 to 0.2.3 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.8.1 to 1.8.3 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.4.3 to 3.4.4 (by dependabot[bot])
- Improve tests for check-mysql (by susisu)

* Wed Jul 27 2022 <mackerel-developers@hatena.ne.jp> - 0.42.2
- Ignores fuse.portal partitions (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.44.37 to 1.44.56 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.8 to 1.8.9 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.10 to 1.8.1 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.4 to 1.10.6 (by dependabot[bot])
- Bump github.com/jmoiron/sqlx from 1.3.4 to 1.3.5 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.4.2 to 3.4.3 (by dependabot[bot])

* Wed Jun 22 2022 <mackerel-developers@hatena.ne.jp> - 0.42.1
- Bump github.com/aws/aws-sdk-go from 1.43.26 to 1.44.37 (by dependabot[bot])

* Wed Mar 30 2022 <mackerel-developers@hatena.ne.jp> - 0.42.0
- [check-aws-cloudwatch-logs] stop gracefully on timeout signal (by pyto86pri)
- [check-aws-cloudwatch-logs] stop gracefully on timeout signal (by pyto86pri)
- Bump github.com/mackerelio/checkers from 0.0.2 to 0.0.3 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.43.12 to 1.43.26 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.1 to 0.2.2 (by dependabot[bot])
- [check-aws-cloudwatch-logs] use FilterLogEventsPages API (by pyto86pri)
- Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.9 to 1.7.10 (by dependabot[bot])

* Tue Mar 15 2022 <mackerel-developers@hatena.ne.jp> - 0.41.7
- Bump github.com/aws/aws-sdk-go from 1.43.7 to 1.43.12 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.22.1 to 3.22.2 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.8 to 1.7.9 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.42.52 to 1.43.7 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.4.1 to 3.4.2 (by dependabot[bot])

* Wed Feb 16 2022 <mackerel-developers@hatena.ne.jp> - 0.41.6
- Bump github.com/fsouza/go-dockerclient from 1.7.7 to 1.7.8 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.42.44 to 1.42.52 (by dependabot[bot])
- upgrade Go: 1.16 -> 1.17 (by lufia)
- Bump github.com/shirou/gopsutil/v3 from 3.21.12 to 3.22.1 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.4 to 1.7.7 (by dependabot[bot])

* Wed Feb 2 2022 <mackerel-developers@hatena.ne.jp> - 0.41.5
- Bump github.com/aws/aws-sdk-go from 1.42.35 to 1.42.44 (by dependabot[bot])
- Bump github.com/aws/aws-sdk-go from 1.42.9 to 1.42.35 (by dependabot[bot])
- Bump github.com/gomodule/redigo from 1.8.6 to 1.8.8 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.21.10 to 3.21.12 (by dependabot[bot])

* Wed Jan 12 2022 <mackerel-developers@hatena.ne.jp> - 0.41.4
- Bump github.com/gomodule/redigo from 1.8.5 to 1.8.6 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.3 to 1.10.4 (by dependabot[bot])

* Wed Dec 1 2021 <mackerel-developers@hatena.ne.jp> - 0.41.3
- Bump github.com/aws/aws-sdk-go from 1.40.59 to 1.42.9 (by dependabot[bot])
- Bump github.com/mackerelio/go-osstat from 0.2.0 to 0.2.1 (by dependabot[bot])

* Thu Nov 18 2021 <mackerel-developers@hatena.ne.jp> - 0.41.2
- Bump github.com/shirou/gopsutil/v3 from 3.21.9 to 3.21.10 (by dependabot[bot])

* Thu Oct 14 2021 <mackerel-developers@hatena.ne.jp> - 0.41.1
- Bump github.com/aws/aws-sdk-go from 1.39.4 to 1.40.59 (by dependabot[bot])
- Bump github.com/shirou/gopsutil/v3 from 3.21.6 to 3.21.9 (by dependabot[bot])
- Bump github.com/fsouza/go-dockerclient from 1.7.3 to 1.7.4 (by dependabot[bot])
- Bump github.com/lib/pq from 1.10.2 to 1.10.3 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.6 to 0.3.7 (by dependabot[bot])
- Bump github.com/go-ldap/ldap/v3 from 3.3.0 to 3.4.1 (by dependabot[bot])

* Wed Oct 6 2021 <mackerel-developers@hatena.ne.jp> - 0.41.0
- update golib, checkers (by yseto)
- [check-log] add search-in-directory option (by yseto)
- [check-redis] migrate redis client library to redigo (by pyto86pri)

* Wed Sep 29 2021 <mackerel-developers@hatena.ne.jp> - 0.40.1
- check-mysql: Closes `checkReplication` rows (by mechairoi)

* Tue Aug 24 2021 <mackerel-developers@hatena.ne.jp> - 0.40.0
- [check-mysql] add --tls, --tls-root-cert and --tls-skip-verify options (by lufia)
- Bump github.com/aws/aws-sdk-go from 1.38.68 to 1.39.4 (by dependabot[bot])

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
