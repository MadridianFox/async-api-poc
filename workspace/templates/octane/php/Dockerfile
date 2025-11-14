FROM dockerhub.greensight.ru/ensi-tech/php-base-image:8.2-master-2023aug21-1-swoole

COPY --from=composer:2.5 /usr/bin/composer /usr/bin/composer

RUN apk add npm git bash htop vim mc

RUN docker-php-ext-install pdo_mysql && \
    docker-php-ext-enable pdo_mysql

ARG USER_ID
ARG GROUP_ID

RUN addgroup -g $GROUP_ID developer
RUN adduser -u $USER_ID -S -D -H -G developer developer
RUN addgroup www-data developer
RUN addgroup developer www-data