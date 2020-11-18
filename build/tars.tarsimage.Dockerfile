FROM tars.base As First
RUN  apt update ; apt install -y curl; curl -sSL https://get.daocloud.io/docker | sh
COPY files/template/tarsimage/root /
COPY files/binary/tarsimage /usr/local/app/tars/tarsimage/bin/tarsimage
RUN  chmod +x /usr/local/app/tars/tarsimage/bin/tarsimage

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]