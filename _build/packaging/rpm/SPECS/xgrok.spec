Name:           %{_product_name}
Version:        %{_product_version}

%if 0%{?rhel} >= 5
Release:        1.el%{?rhel}
%else
Release:        1%{?dist}
%endif

Summary:        Introspected tunnels to localhost.
Group:          Development/Tools
License:        Apache License, Version 2.0
Source0:        %{name}_linux_amd64.zip
Source1:        %{name}.sysconfig
Source2:        %{name}.service
Source3:        %{name}.init
Source4:        %{name}.config.yml
Source5:        %{name}.logrotate
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
Introspected tunnels to localhost.

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp xgrok %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/sysconfig
cp %{SOURCE1} %{buildroot}/%{_sysconfdir}/sysconfig/%{name}
mkdir -p %{buildroot}/%{_sysconfdir}/logrotate.d/
mkdir -p %{buildroot}/%{_sysconfdir}/%{name}
cp %{SOURCE4} %{buildroot}/%{_sysconfdir}/%{name}/%{name}.yml
cp %{SOURCE5} %{buildroot}/%{_sysconfdir}/logrotate.d/%{name}
mkdir -p %{buildroot}/var/log/%{name}
mkdir -p %{buildroot}/%{_initrddir}
cp %{SOURCE3} %{buildroot}/%{_initrddir}/xgrok

%pre
getent group xgrok >/dev/null || groupadd -r xgrok
getent passwd xgrok >/dev/null || \
    useradd -r -g xgrok -d /var/lib/xgrok -s /sbin/nologin \
    -c "xgrok user" xgrok
exit 0

%post
/sbin/chkconfig --add %{name}

%preun
if [ "$1" = 0 ] ; then
    /sbin/service %{name} stop >/dev/null 2>&1
    /sbin/chkconfig --del %{name}
fi

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%attr(644, root, root) %{_sysconfdir}/%{name}/%{name}.yml
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}
%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
%{_unitdir}/%{name}.service
%else
%attr(755, root, root) %{_initrddir}/%{name}
%endif
%attr(755, root, root) %{_bindir}/%{name}
%attr(644, root, root) %{_sysconfdir}/logrotate.d/%{name}
%dir %attr(755, xgrok, xgrok) /var/log/%{name}

%doc
