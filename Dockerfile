FROM alpine:3.6
MAINTAINER Sebastian Döll <sebastian.doell@axelspringer.com>

ADD \
    /bin/moppi_0.0.1_linux_amd64 /bin/moppi

ADD \
    scripts/init.sh /init.sh

RUN \
    chmod +x /bin/moppi && \
    chmod +x /init.sh

EXPOSE 80

STOPSIGNAL SIGTERM

ENTRYPOINT ["/init.sh"]

CMD ["/bin/moppi"]
