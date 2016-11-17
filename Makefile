SHELL:=/bin/bash
RGREP=grep -rn --color=always --exclude=.* --exclude-dir=Godeps --exclude-dir=vendor --exclude=Makefile

.PHONY: find-uk0 find-bugs-todos find-exits

all:
	@echo "Don't try to use the makefile to compile bytemark-client! Set up"
	@echo "your go environment as you would for any other project, then run"
	@echo "this command to download, build and install bytemark-client"
	@echo 
	@echo "    go get github.com/BytemarkHosting/bytemark-client/cmd/bytemark"
	@echo
	@echo "To run all the tests, run go test ./..."

%.pdf: %.ps
	ps2pdf $< $@

doc/bytemark-client.ps: doc/bytemark.1
	groff -mandoc -T ps $< > $@

# find instances of uk0 in the code (to ensure that URLs aren't hardcoded in (too) many places
find-uk0: 
	$(RGREP) --exclude=bytemark "uk0" .

# find instances of BUG and TODO comments.
find-bugs-todos:
	$(RGREP) -P "// BUG(.*):" . || echo ""
	$(RGREP) -P "// TODO(.*):" .

# find every line that calls panic or os.Exit (from when I was writing util.ProcessError
find-exits:
	$(RGREP) --exclude=exit.go --exclude=main.go -P "panic\(|os.Exit" .
