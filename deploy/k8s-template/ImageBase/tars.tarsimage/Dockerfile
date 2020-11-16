FROM tars.base
RUN  apt update -y docker
COPY /root /
RUN chmod a+x /bin/entrypoint.sh
CMD ["/bin/entrypoint.sh"]