SHELL:=/bin/bash

ALL_PACKAGES := bytemark.co.uk/client/lib bytemark.co.uk/client/cmds/util bytemark.co.uk/client/cmds bytemark.co.uk/client/cmd/bytemark
ALL_FILES := lib/*.go cmds/*.go cmds/util/*.go mocks/*.go util/*/*.go cmd/*/*.go

BUILD_NUMBER ?= 0

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
	GO15VENDOREXPERIMENT=1 go build -o bytemark bytemark.co.uk/client/cmd/bytemark

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

clean:
	rm -rf Bytemark.app rm $(LAUNCHER_APP)
	rm -f bytemark
	rm -f main.coverage lib.coverage
	rm -f main.coverage.html lib.coverage.html

gensrc:
	BUILD_NUMBER=$(BUILD_NUMBER) go generate ./...

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
	go test -coverprofile=$@ bytemark.co.uk/client/cmd/bytemark

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

find-bugs-todos:
	$(RGREP) -P "// BUG(.*):" . || echo ""
	$(RGREP) -P "// TODO(.*):" .

find-exits:
	$(RGREP) --exclude=exit.go --exclude=main.go -P "panic\(|os.Exit" .
