Cmd = go
BinariesName = ElasticView
GoBuild = $(GoCmd) build -v -ldflags="-w -s" -o $(BinariesName)
VueBuild = cd vue; npm run build:prod; cd ../
RawOS = ""

ifeq ($OS, Windows_NT)
	BinariesName = "ElasticView.exe"
	RawOS = windows
else
 	ifeq ($(shell uname),Darwin)
  		BinariesName="ElasticView"
  		RawOS = darwin
 	else
  		BinariesName="ElasticView"
  		RawOS = linux
 	endif
endif

all: build

build:
	$(VueBuild)
	$(Cmd) clean
	$(GBuild)

clean:
	$(Cmd) clean

windows:
	$(VueBuild)
	export CGO_ENABLED=0
	export GOOS=windows
	export GOARCH=amd64
	$(GBuild)
	export GOOS=$(RawOS)

darwin:
	$(VueBuild)
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	$(GBuild)
	export GOOS=$(RawOS)

linux:
	$(VueBuild)
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GBuild)
	export GOOS=$(RawOS)

help:
	echo "flag support: build, clean, windows, darwin, linux"
