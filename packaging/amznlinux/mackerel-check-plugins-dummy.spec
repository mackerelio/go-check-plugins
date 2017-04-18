%define _binaries_in_noarch_packages_terminate_build   0
%define _localbindir /usr/local/bin
%define __targetdir /usr/bin
%define __oldtargetdir /usr/local/bin

Name:      mackerel-check-plugins
Version:   0.5.0
Release:   2
License:   Commercial
Summary:   macekrel.io check plugins
URL:       https://mackerel.io
Group:     Hatena
Packager:  Hatena
BuildArch: noarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
dummy package

%prep

%build

%install

%clean

%pre

%post

%preun

%files
