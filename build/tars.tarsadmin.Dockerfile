FROM tars.base As First
COPY files/template/tarsadmin/root /
COPY files/binary/tarsadmin /usr/local/app/tars/tarsadmin/bin/tarsadmin
RUN  chmod +x /usr/local/app/tars/tarsadmin/bin/tarsadmin

#　第二阶段
FROM scratch
COPY --from=First / /
CMD ["/bin/entrypoint.sh"]