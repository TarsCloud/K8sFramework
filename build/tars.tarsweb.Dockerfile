FROM node:10-stretch-slim AS First
COPY files/sources.list /etc/apt/sources.list
COPY files/template/tarsweb/root /
COPY files/k8s-web /k8s-web

RUN apt update -y && apt install python build-essential busybox -y && busybox --install
RUN cd /k8s-web && rm -f package-lock.json 
#&& npm install --registry=http://registry.upchinaproduct.com

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