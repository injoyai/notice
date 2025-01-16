
# 编译

name="notice"
GOOS=windows GOARCH=amd64 go build -v -ldflags="-s -w -extldflags -static" -o ./$name.exe
echo "$name 编译完成..."
echo "开始压缩..."
upx -9 -k "./$name.exe"
if [ -f "./$name.ex~" ]; then
  rm "./$name.ex~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi
echo "压缩完成..."

name="$notice"_linux_arm
GOOS=linux GOARCH=arm GOARM=7 go build -v -ldflags="-w -s" -o ./$name
echo "$name 编译完成..."
echo "开始压缩..."
upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi

name="$notice"_linux_arm64
GOOS=linux GOARCH=arm64 go build -v -ldflags="-w -s" -o ./$name
echo "$name 编译完成..."
echo "开始压缩..."
upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi

name="$notice"_linux_amd64
GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./$name
echo "$name 编译完成..."
echo "开始压缩..."
upx -9 -k "./$name"
if [ -f "./$name.~" ]; then
  rm "./$name.~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi


echo "稍后退出.."
sleep 8