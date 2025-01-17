
name="notice"

GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./$name
echo "amd64编译完成..."
echo "开始压缩..."
#upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi
docker build --platform=linux/amd64 --push -t docker.io/injoyai/notice:latest -f ./Dockerfile .

GOOS=linux GOARCH=arm64 go build -v -ldflags="-w -s" -o ./$name
echo "amd64编译完成..."
echo "开始压缩..."
#upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi
docker build --platform=linux/arm64 --push -t docker.io/injoyai/notice:latest -f ./Dockerfile .


GOOS=linux GOARCH=arm GOARM=7 go build -v -ldflags="-w -s" -o ./$name
echo "amd64编译完成..."
echo "开始压缩..."
#upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi
docker build --platform=linux/arm --push -t docker.io/injoyai/notice:latest -f ./Dockerfile .

sleep 8