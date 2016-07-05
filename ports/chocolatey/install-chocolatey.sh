#!/bin/bash

if [ `which apt-get` ]; then
	sudo apt-get install -y --force-yes curl unzip mono-runtime libmono-system-core4.0-cil libmono-system-componentmodel-dataannotations4.0-cil libmono-windowsbase4.0-cil libmono-system-xml-linq4.0-cil || exit 1
elif [ `which brew` ]; then
	brew install mono
fi

mkdir -p chocolatey

pushd chocolatey
curl $(curl 'http://chocolatey.org/install.ps1' | grep '\.nupkg' | grep 'https' | head -n 1 | grep -o 'http.*nupkg') > chocolatey.nupkg || exit 1
unzip chocolatey.nupkg || exit 1
mv tools/chocolateyInstall/choco.exe ..
popd
rm -rf chocolatey
