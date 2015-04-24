ALL_PACKAGES := bigv.io/client/lib bigv.io/client/cmd

.PHONY: test update-dependencies
.PHONY: BigV.app

all: go-bigv

BigV.app: go-bigv
	mkdir -p BigV.app/Contents/Resources/bin
	mkdir -p BigV.app/Contents/MacOS
	cp ports/mac/main BigV.app/Contents/MacOS
	cp ports/mac/PkgInfo BigV.app/Contents
	cp ports/mac/Info.plist BigV.app/Contents
	cp -r ports/mac/Resources/* BigV.app/Contents/Resources
	cp go-bigv BigV.app/Contents/Resources/bin
	ln -s ../Resources/bin/go-bigv BigV.app/Contents/MacOS

clean:
	rm -rf BigV.app
	rm -f go-bigv

go-bigv: cmd/*.go lib/*.go
	go build -o go-bigv bigv.io/client/cmd

install: all
	cp go-bigv /usr/bin/go-bigv

test:
	go test $(ALL_PACKAGES)

update-dependencies: 
	go get -ut $(ALL_PACKAGES)
	godep update $(ALL_PACKAGES)
