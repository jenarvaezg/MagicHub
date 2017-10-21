export CGO_ENABLED=0
export GOOS=linux
export MAGICBOX_VERSION=$(cat version.txt)

go build -a -installsuffix cgo -o magicbox .
docker build -t jenarvaezg/magicbox:$MAGICBOX_VERSION .

