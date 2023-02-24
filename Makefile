

.PHONY:	scan
scan: 
	go list -json -deps |  nancy sleuth
	trivy fs . 

.PHONY: build
build: 
	goreleaser build --clean

.PHONY: build-snapshot
build-snapshot: 
	goreleaser build --clean --snapshot --single-target


.PHONY: release-skip-publish
release-skip-publish: 
	goreleaser release --clean --skip-publish 

.PHONY: release-snapshot
release-snapshot: 
	goreleaser release --clean --skip-publish --snapshot


.PHONY: lint
lint: 
	golangci-lint run ./... 


.PHONY: changelog
changelog: 
	git-chglog -o CHANGELOG.md 

##TODO: 
.PHONY: test
test:
	echo todo 
	


.PHONY: gen-doc
gen-doc: 
	mkdir -p ./doc
	./dist/oidc-server-demo_linux_amd64_v1/oidc-server documentation  --dir ./doc

.PHONY: doc-site
doc-site: 
	podman  run --rm -it -p 8000:8000 -v ${PWD}/www:/docs:z squidfunk/mkdocs-material 
