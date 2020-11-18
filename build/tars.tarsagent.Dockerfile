FROM tars.base As First
COPY files/template/tarsagent/root /
COPY files/binary/tarsagent /usr/local/app/tars/tarsagent/bin/tarsagent
RUN  chmod +x /usr/local/app/tars/tarsagent/bin/tarsagent

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]