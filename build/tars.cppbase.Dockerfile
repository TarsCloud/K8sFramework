FROM tars.base As First
COPY files/template/tarsnode/root /
COPY files/binary/tarsnode /usr/local/app/tars/tarsnode/bin/tarsnode
RUN  chmod +x /usr/local/app/tars/tarsnode/bin/tarsnode

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]