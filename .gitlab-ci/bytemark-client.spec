Name: bytemark-client
Version: %{version}
Release: %{release} 
Summary: Command-line client for Bytemark Hosting's self-service hosting products.
License: MIT
URL: http://github.com/BytemarkHosting/bytemark-client
Source0: bytemark
Source1: bytemark.1

%description
bytemark-client provides an executable named 'bytemark' which is used to access Bytemark Hosting's self-service products.
This tool allows you to create, view, alter and delete any servers and groups you may have hosted by Bytemark, as well as
modify your account details.

%prep

%build

%install
install -d %{buildroot}%{_bindir}
install -d %{buildroot}%{_mandir}/man1
install -p -m 755 %{SOURCE0} %{buildroot}%{_bindir}
install -p -m 644 %{SOURCE1} %{buildroot}%{_mandir}/man1

%files
%defattr(-,root,root)
%{_bindir}/bytemark
%{_mandir}/man1/bytemark.1*

%changelog
