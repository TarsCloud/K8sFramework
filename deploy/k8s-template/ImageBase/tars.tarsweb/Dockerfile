FROM node:lts-stretch-slim
COPY /root /
COPY /k8s-web /k8s-web
RUN chmod a+x /bin/entrypoint.sh
CMD ["/bin/entrypoint.sh"]
