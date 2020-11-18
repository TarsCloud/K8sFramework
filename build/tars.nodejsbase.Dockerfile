FROM node:10-stretch-slim AS First
COPY files/sources.list /etc/apt/sources.list
COPY files/template/tarsnode/root /
COPY files/binary/tarsnode /usr/local/app/tars/tarsnode/bin/tarsnode
RUN  chmod +x /usr/local/app/tars/tarsnode/bin/tarsnode

# 设置时区
RUN  ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN  echo Asia/Shanghai > /etc/timezone

RUN  apt update -y
RUN  apt install busybox -y
RUN  busybox --install

#COPY files/tars-node/src/tars-node-agent /usr/local/app/tars/tars-node-agent
#RUN cd /usr/local/app/tars/tars-node-agent && npm install --registry=http://registry.upchinaproduct.com

# 设置别名，兼容使用习惯
RUN echo alias ll=\'ls -l\' >> /etc/bashrc

# 清理多余文件
RUN  apt purge -y
RUN  apt clean all
RUN  rm -rf /var/lib/apt/lists/*
RUN  rm -rf /var/cache/*.dat-old
RUN  rm -rf /var/log/*.log /var/log/*/*.log

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]