FROM gogs/gogs:0.11.19
COPY ./bin/exposer /bin/exposer
COPY ./app.ini /data/gogs/conf/app.ini
COPY ./init.sh /app/gogs/docker/init.sh
COPY ./gogs.db /data/gogs/data/gogs.db
COPY ./bin/initial /data/gogs/docker/initial

# git
RUN apk upgrade --update; apk add openssh-client git
#Disable stric host check
COPY ssh_config /root/.ssh/config

ENTRYPOINT ["/app/gogs/docker/init.sh"]
CMD ["/bin/s6-svscan", "/app/gogs/docker/s6/"]
