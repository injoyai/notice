FROM alpine:latest

LABEL maintainer="injoy"

###############################################################################
#                                INSTALLATION
###############################################################################

# 设置固定的项目路径
ENV WORKDIR /root/notice

# 安装时区数据
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 添加应用可执行文件，并设置执行权限
ADD ./notice   $WORKDIR/notice
RUN chmod +x $WORKDIR/notice
RUN mkdir  $WORKDIR/config

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD $WORKDIR/notice
