FROM tars.base
RUN  apt update ; apt install -y curl; curl -sSL https://get.daocloud.io/docker | sh
COPY /root /
RUN chmod a+x /bin/entrypoint.sh
CMD ["/bin/entrypoint.sh"]