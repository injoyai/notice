
name="notice"

GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./$name
echo "amd64编译完成..."
echo "开始压缩..."
upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi


docker build --push -t 192.168.10.102:30002/qianlang-public/notice:latest -f ./Dockerfile .

sleep 8