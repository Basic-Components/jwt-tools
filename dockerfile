FROM alpine:3.10
ADD bin/linux-amd64/jwtcenter /code/jwtcenter
WORKDIR /code
CMD ["./jwtcenter"]