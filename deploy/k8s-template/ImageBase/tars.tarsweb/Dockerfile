FROM node:lts-stretch-slim
COPY /root /
COPY /tars-web-k8s /tars-web-k8s
RUN chmod a+x /bin/entrypoint.sh
CMD ["/bin/entrypoint.sh"]
