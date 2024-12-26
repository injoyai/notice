name="notice"
name_linux_arm="notice_linux_arm"
name_linux_amd64="notice_linux_amd64"

# 编译
GOOS=windows GOARCH=amd64 go build -v -ldflags="-s -w -extldflags -static" -o ./$name.exe
echo "windows编译完成..."
echo "开始压缩..."
upx -9 -k "./$name.exe"
if [ -f "./$name.ex~" ]; then
  rm "./$name.ex~"
fi
if [ -f "./$name.000" ]; then
  rm "./$name.000"
fi
echo "压缩完成..."

GOOS=linux GOARCH=arm GOARM=7 go build -v -ldflags="-w -s" -o ./$name_linux_arm
echo "arm7编译完成..."
echo "开始压缩..."
upx -9 -k "./$name_linux_arm"
if [ -f "./$name_linux_arm.~" ]; then
  rm "./$name_linux_arm.~"
fi
if [ -f "./$name_linux_arm.000" ]; then
  rm "./$name_linux_arm.000"
fi

GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./$name_linux_amd64
echo "amd64编译完成..."
echo "开始压缩..."
upx -9 -k "./$name_linux_amd64"
if [ -f "./$name_linux_amd64.~" ]; then
  rm "./$name_linux_amd64.~"
fi
if [ -f "./$name_linux_amd64.000" ]; then
  rm "./$name_linux_amd64.000"
fi


# 上传到minio
echo "上传到minio..."
CMD.exe /C "in upload minio ./$name.exe"

echo "稍后退出.."
sleep 8