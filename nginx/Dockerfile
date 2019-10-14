FROM alpine:latest
LABEL maintainer="Anton Illarionov"

RUN apk add nginx && mkdir -p /run/nginx

COPY default.conf /etc/nginx/conf.d/
COPY index.html /var/lib/nginx/html/

CMD [ "nginx", "-g", "daemon off;" ]
