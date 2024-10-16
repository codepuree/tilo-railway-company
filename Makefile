OUT_PATH=bin
BIN_PATH=${OUT_PATH}/trc

build:
	rm -fr /var/www && mkdir /var/www && ln -s /workspace/tilo-railway-company/web/static /var/www/static && ln -s /workspace/tilo-railway-company/customScripts /var/www/custom

build_real:
	go build -o "${BIN_PATH}" cmd/main.go

release:
	GOOS=linux GOARCH=arm GOARM=7 go build -a -tags netgo -ldflags '-w' -o ${BIN_PATH} ./cmd/

# Extract trclib for yaegi
extract_trclib:
	mkdir -p ./pkg/trclib && cd ./pkg/trclib && \
	yaegi extract -name trclib github.com/codepuree/tilo-railway-company/pkg/traincontrol && \
	sed -i "s/\(func init() {\)/var Symbols = map[string]map[string]reflect.Value {}\n\n\1/" ./github_com-codepuree-tilo-railway-company-pkg-traincontrol.go
	
clean:
	rm -fr ${OUT_PATH}
