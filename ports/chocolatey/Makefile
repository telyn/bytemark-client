.PHONY: always-rebuild


bytemark.nupkg: package/bytemark.nuspec package/tools/chocolateyInstall.ps1 choco.exe
	cd package && mono ../choco.exe pack bytemark.nuspec
	mv package/bytemark.$$(<VERSION).nupkg bytemark.nupkg

package/%: %.m4 VERSION
	mkdir -p package/tools
	m4 -D VERSION=$$(<VERSION) $< > $@

choco.exe:
	./install-chocolatey.sh

VERSION: always-rebuild
	echo "$(VERSION)" > $@.tmp
	diff $@ $@.tmp || cp $@.tmp $@
	rm -f $@.tmp
