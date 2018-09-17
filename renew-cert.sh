docker run -it --rm --name letsencrypt \
 -v "/var/whoismining/nginx/ssl:/etc/letsencrypt" \
 --volumes-from nginx \
 certbot/certbot \
 certonly \
 --webroot \
 --webroot-path /var/www/whoismining.com \
 --agree-tos \
 --renew-by-default \
 -d whoismining.com \
 -m hi@whoismining.com

 docker kill --signal=HUP nginx
