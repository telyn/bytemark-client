ALL_PACKAGES := bigv.io/client/lib bigv.io/client/cmd
ALL_FILES := cmd/*.go lib/*.go
OSAARCH:=x86_64
ifeq ($(GOARCH),386)
OSAARCH:=i386
endif
LAUNCHER_APP:=ports/mac/launcher.app

.PHONY: test update-dependencies
.PHONY: BigV.app

all: go-bigv

BigV.app: go-bigv $(LAUNCHER_APP) ports/mac/*
	mkdir -p BigV.app/Contents/Resources/bin
	mkdir -p BigV.app/Contents/Resources/Scripts
	mkdir -p BigV.app/Contents/MacOS
	# pilfer the applet binary, applescript and resource file from the compiled script
	cp $(LAUNCHER_APP)/Contents/Resources/Scripts/main.scpt BigV.app/Contents/Resources/Scripts
	cp $(LAUNCHER_APP)/Contents/Resources/applet.rsrc BigV.app/Contents/Resources
	cp $(LAUNCHER_APP)/Contents/MacOS/applet BigV.app/Contents/MacOS/launcher
	# then put in our own Info.plist which has BigV branding and copyright and paths and stuff
	cp ports/mac/Info.plist BigV.app/Contents
	# copy in the terminal profile and start script
	cp -r ports/mac/BigV.terminal BigV.app/Contents/Resources
	cp -r ports/mac/start BigV.app/Contents/Resources
	# copy in go-bigv into its own folder (this allows us to say 
	# "add BigV.app/Contents/Resources/bin to your PATH" and it'll only add go-bigv
	# and not the launcher too.)
	cp go-bigv BigV.app/Contents/Resources/bin
	# make a symlink into MacOS. This step is totally unnecessary but it means all the binaries live in MacOS which is nice I guess?
	rm -f BigV.app/Contents/MacOS/go-bigv
	ln -s ../Resources/bin/go-bigv BigV.app/Contents/MacOS

clean:
	rm -rf BigV.app rm $(LAUNCHER_APP)
	rm -f go-bigv
	rm -f cmd.coverage lib.coverage
	rm -f cmd.coverage.html lib.coverage.html


go-bigv: $(ALL_FILES)
	go build -o go-bigv bigv.io/client/cmd

$(LAUNCHER_APP): ports/mac/launcher-script.txt
ifeq (Darwin, $(shell uname -s))
	rm -rf $@
	osacompile -a $(OSAARCH) -x -o $@ $<
else
	echo "WARNING using old pre-built launcher."
endif

install: all
	cp go-bigv /usr/bin/go-bigv

coverage: lib.coverage.html cmd.coverage.html
	open lib.coverage.html
	open cmd.coverage.html

%.coverage.html: %.coverage
	go tool cover -html=$< -o $@

%.coverage: % %/*
	go test -coverprofile=$@ bigv.io/client/$<

test: 
	go test $(ALL_PACKAGES)

update-dependencies: 
	go get -ut $(ALL_PACKAGES)
	godep update $(ALL_PACKAGES)
