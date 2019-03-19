FROM nginx:1.15.9

COPY nginx.conf /etc/nginx/nginx.conf
COPY mime.types /etc/nginx/mime.types

ADD frontend_dir/public /www/public

# stub. remove this and mount volume for logs
RUN mkdir /www/logs

EXPOSE 80