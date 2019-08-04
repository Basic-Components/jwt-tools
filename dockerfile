FROM alpine:3.10
ADD bin/linux-amd64/jwtrpc /code/jwtrpc
WORKDIR /code
CMD ["./jwtrpc"]