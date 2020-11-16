FROM tars.cppbase
RUN apt install -y node
COPY /root /
RUN chmod a+x /bin/entrypoint.sh
COPY /tars-node-agent /usr/local/app/tars/