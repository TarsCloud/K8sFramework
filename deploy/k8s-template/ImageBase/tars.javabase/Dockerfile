#　第一阶段
FROM tars.cppbase AS First

# 第二阶段
FROM openjdk:8-slim-stretch AS Second
COPY source.list /etc/apt/sources.list
RUN  apt update -y && apt dist-upgrade -y
COPY --from=First / /
COPY root /
RUN chmod a+x /bin/entrypoint.sh

# 清理多余文件
RUN  apt update -y
RUN  apt autoremove -y 
RUN  apt purge -y
RUN  apt clean all
RUN  rm -rf /var/lib/apt/lists/*
RUN  rm -rf /var/cache/*.dat-old
RUN  rm -rf /var/log/*.log /var/log/*/*.log

#　第三阶段
FROM scratch
COPY --from=Second / /
ENV JAVA_HOME=/usr/local/openjdk-8
ENV PATH=${PATH}:${JAVA_HOME}/bin