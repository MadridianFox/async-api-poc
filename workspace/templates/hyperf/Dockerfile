FROM hyperf/hyperf:8.3-alpine-v3.19-swoole

COPY --from=composer:2.5 /usr/bin/composer /usr/bin/composer

RUN apk add git bash htop vim mc less