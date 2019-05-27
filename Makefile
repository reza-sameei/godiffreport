.PHONY: help deps test build mods

unexport GOPATH

help :
	@echo "deps          Update dependencies list (go-mod/go-get)"
	@echo "build         Build main module"
	@echo "test          Test main module"
	@echo "mods          Show a list of modules included in the project: the main and deps"

build :
	go build ./...

test :
	@# go test -v -cover gitlab.com/trialblaze/etl-common/pkg/data/sdxmldef
	go test -v -cover gitlab.com/trialblaze/etl-common/pkg/operations/sddm/build -run Test_Build
	@# go test -v -cover gitlab.com/trialblaze/etl-common/pkg/data/formdef
	@# go test -v -cover gitlab.com/trialblaze/etl-common/pkg/data/diffreport
	@# go test ./...

# A Place to put early access modules/dependencies to make it possible to update
# them easily; but for stablished dependencies consider adding them directly in
# go.mod file
deps :
	@# echo "No early access dep"
	go get -u https://github.com/nsf/jsondiff@master


mods :
	go list -m all
