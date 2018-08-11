FROM cargo.caicloudprivatetest.com/caicloud/debian:jessie

# Set the default timezone to Shanghai
RUN echo "Asia/Shanghai" > /etc/timezone
RUN dpkg-reconfigure -f noninteractive tzdata

WORKDIR /caicloud/dashboard-admin
COPY ./bin/admin /caicloud/dashboard-admin/

EXPOSE 2587

CMD ["/caicloud/dashboard-admin/admin"]
