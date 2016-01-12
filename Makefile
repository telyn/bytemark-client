SHELL:=/bin/bash

ALL_PACKAGES := bytemark.co.uk/client/lib bytemark.co.uk/client/cmds/util bytemark.co.uk/client/cmds bytemark.co.uk/client
ALL_FILES := *.go lib/*.go cmds/*.go cmds/util/*.go mocks/*.go util/*/*.go

MAJORVERSION := 0
MINORVERSION := 1
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
BUILD_NUMBER ?= 0
GIT_BRANCH ?= `git rev-parse --abbrev-ref HEAD`
GIT_COMMIT ?= `git rev-parse HEAD`
VERSIONFILE := lib/version.go

OSAARCH:=x86_64
ifeq ($(GOARCH),386)
OSAARCH:=i386
endif
LAUNCHER_APP:=ports/mac/launcher.app
RGREP=grep -rn --color=always --exclude=.* --exclude-dir=Godeps --exclude=Makefile

.PHONY: test update-dependencies
.PHONY: Bytemark.app
.PHONY: find-uk0 find-bugs-todos find-exits
.PHONY: gensrc

all: bytemark

bytemark: $(ALL_FILES) gensrc
	GO15VENDOREXPERIMENT=1 go build -o bytemark bytemark.co.uk/client

Bytemark.app.zip: Bytemark.app
	zip -r $@ $<

Bytemark.app: bytemark $(LAUNCHER_APP) ports/mac/*
	mkdir -p Bytemark.app/Contents/Resources/bin
	mkdir -p Bytemark.app/Contents/Resources/Scripts
	mkdir -p Bytemark.app/Contents/MacOS
	# pilfer the applet binary, applescript and resource file from the compiled script
	cp $(LAUNCHER_APP)/Contents/Resources/Scripts/main.scpt Bytemark.app/Contents/Resources/Scripts
	cp $(LAUNCHER_APP)/Contents/Resources/applet.rsrc Bytemark.app/Contents/Resources
	cp $(LAUNCHER_APP)/Contents/MacOS/applet Bytemark.app/Contents/MacOS/launcher
	# then put in our own Info.plist which has Bytemark branding and copyright and paths and stuff
	cp ports/mac/Info.plist Bytemark.app/Contents
	# copy in the terminal profile and start script
	cp -r ports/mac/Bytemark.terminal Bytemark.app/Contents/Resources
	cp -r ports/mac/start Bytemark.app/Contents/Resources
	# copy in bytemark into its own folder (this allows us to say 
	# "add Bytemark.app/Contents/Resources/bin to your PATH" and it'll only add bytemark
	# and not the launcher too.)
	cp bytemark Bytemark.app/Contents/Resources/bin
	# make a symlink into MacOS. This step is totally unnecessary but it means all the binaries live in MacOS which is nice I guess?
	rm -f Bytemark.app/Contents/MacOS/bytemark
	ln -s ../Resources/bin/bytemark Bytemark.app/Contents/MacOS
	# sign the code? anyone? shall we sign the code?

gensrc:
	rm -f $(VERSIONFILE)
	@echo "Writing $(VERSIONFILE)"
	@echo "package lib" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  majorversion = $(MAJORVERSION)" >> $(VERSIONFILE)
	@echo "  minorversion = $(MINORVERSION)" >> $(VERSIONFILE)
	@echo "  buildnumber = $(BUILD_NUMBER)" >> $(VERSIONFILE)
	@echo "  gitcommit = \"$(GIT_COMMIT)\"" >> $(VERSIONFILE)
	set -x; GIT_BRANCH=HEAD; for i in master $$(git for-each-ref --format "%(refname)" refs/heads/release\*) develop $$(git for-each-ref --format "%(refname)" refs/heads); do \
	    [ "`git rev-parse $$i`" == "`git rev-parse HEAD`" ] && GIT_BRANCH=`git rev-parse --abbrev-ref $$i`; \
	done; \
	echo "  gitbranch = \"$$GIT_BRANCH\"" >> $(VERSIONFILE);
	@echo "  builddate = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)
	@cat $(VERSIONFILE)

clean:
	rm -rf Bytemark.app rm $(LAUNCHER_APP)
	rm -f bytemark
	rm -f main.coverage lib.coverage
	rm -f main.coverage.html lib.coverage.html

checkinstall: 
	checkinstall -D --install=no -y --maintainer="telyn@bytemark.co.uk" \
	    --pkgname=bytemark-client --pkgversion="$(MAJORVERSION).$(MINORVERSION).$(BUILD_NUMBER)" \
	    --requires="" \
	    --strip=no --stripso=no

$(LAUNCHER_APP): ports/mac/launcher-script.txt
ifeq (Darwin, $(shell uname -s))
	rm -rf $@
	osacompile -a $(OSAARCH) -x -o $@ $<
else
	echo "WARNING using old pre-built launcher."
endif

install: all
	cp bytemark /usr/bin/bytemark

coverage: lib.coverage.html main.coverage.html cmds.coverage.html 
ifeq (Darwin, $(shell uname -s))
	open lib.coverage.html
	open main.coverage.html
	open cmds.coverage.html
else
	xdg-open lib.coverage.html
	xdg-open main.coverage.html
	xdg-open cmds.coverage.html
endif

main.coverage: *.go
	go test -coverprofile=$@ bytemark.co.uk/client

%.coverage.html: %.coverage
	go tool cover -html=$< -o $@

%.coverage: % %/*
	go test -coverprofile=$@ bytemark.co.uk/client/$<

docs: doc/*.md
	for file in doc/*.md; do \
	    pandoc --from markdown --to html $$file --output $${file%.*}.html; \
	done

test: gensrc
ifdef $(VERBOSE)
	GO15VENDOREXPERIMENT=1 go test -v $(ALL_PACKAGES)
else 
	GO15VENDOREXPERIMENT=1 go test $(ALL_PACKAGES)
endif

find-uk0: 
	$(RGREP) --exclude=bytemark "uk0" .

find-bigv:
	$(RGREP) --exclude=bytemark -i bigv . | grep -v "bigv.io"

find-bugs-todos:
	$(RGREP) -P "// BUG(.*):" . || echo ""
	$(RGREP) -P "// TODO(.*):" .

find-exits:
	$(RGREP) --exclude=exit.go --exclude=main.go -P "panic\(|os.Exit" .
