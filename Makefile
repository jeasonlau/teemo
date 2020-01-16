BUILD_ENV := CGO_ENABLED=0
BUILD=`date +%FT%T%z`
VENDOR=-mod=vendor

TARGET_EXEC := teemo

#.PHONY: all clean setup build-linux build-osx build-windows setup-linux setup-osx setup-windows pack-linux pack-osx pack-windows
.PHONY: all clean setup build-linux build-windows setup-linux setup-windows pack-linux pack-windows

#all: clean setup build-linux build-osx build-windows
all: clean setup build-linux build-windows

#release: all pack-linux pack-osx pack-windows
release: all pack-linux pack-windows

clean:
	rm -rf build

#setup: setup-linux setup-osx setup-windows
setup: setup-linux setup-windows

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${VENDOR} -o build/linux/${TARGET_EXEC}

#build-osx: setup-osx
#	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${VENDOR} -o build/osx/${TARGET_EXEC}

build-windows:
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${VENDOR} -o build/win/${TARGET_EXEC}.exe

pack-linux:
	upx build/linux/${TARGET_EXEC}
#	upx build/linux/${TARGET_EXEC} && zip -r teemo.zip build/linux

#pack-osx:
#	upx build/osx/${TARGET_EXEC} && zip -r teemo.zip build/osx

pack-windows:
	upx build/win/${TARGET_EXEC}.exe
#	upx build/win/${TARGET_EXEC}.exe && zip -r teemo.zip build/win