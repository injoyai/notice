FROM alpine:latest

LABEL maintainer="injoy"

###############################################################################
#                                INSTALLATION
###############################################################################

# 设置固定的项目路径
ENV WORKDIR /root/notice

# 添加应用可执行文件，并设置执行权限
ADD ./notice   $WORKDIR/notice
RUN chmod +x $WORKDIR/notice
RUN mkdir  $WORKDIR/config

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD $WORKDIR/notice
