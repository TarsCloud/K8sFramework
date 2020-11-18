FROM tars.base As First
COPY files/template/tarscontrol/root /
COPY files/binary/tarscontrol /usr/local/app/tars/tarscontrol/bin/tarscontrol
RUN  chmod +x /usr/local/app/tars/tarscontrol/bin/tarscontrol

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]