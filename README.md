# Gopherfile
The best way (in my opinion) to build a Docker image for Go.


## Intro
Both Go and Docker are great platforms. Knowing that Docker is built with Go, you'd think that they would *go* perfectly hand-in-hand. Don't get me wrong, they do work nicely with one another, but it can be better.


## The Issue
In the `src` directory of this repo, `main.go` contains a simple web server. If I compile this image to a standard binary, the output size is `7.3M`.
```
➜  go build main.go
➜  ls -lah main
-rwxr-xr-x  1 tanner  staff   7.3M Sep  9 20:31 main
```
If I use the [official Golang Docker library](https://hub.docker.com/_/golang/) to create an Alpine-based Docker image (which is already considered small), I end up getting an output size of `334MB`. Unfortunately, this image is bulky compared to Go's standard binary.
```
➜  docker build -t the-old-way .
Sending build context to Docker daemon  13.44MB
Step 1/5 : FROM golang:alpine
alpine: Pulling from library/golang
...
Successfully built 04444a05b502
Successfully tagged the-old-way:latest
➜  docker image ls
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
the-old-way         latest              04444a05b502        15 seconds ago      334MB
golang              alpine              ffe224d301bc        4 days ago          311MB
```
At the end of the day, a smaller image should always be desired, especially when running numerous services in a production environment.

## The Solution
Introducing **Gopherfile**. I created Gopherfile for the people of the Internet to use as a good "default" Dockerfile when dealing with Go images.
```
➜  go mod vendor
➜  docker build -t gopherimg .
Sending build context to Docker daemon  13.44MB
Step 1/15 : FROM golang:alpine as build
alpine: Pulling from library/golang
...
Successfully built 6f17f1db28ab
Successfully tagged gopherimg:latest
➜  docker image ls
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
gopherimg           latest              6f17f1db28ab        18 seconds ago      5.59MB
<none>              <none>              233dba4a5b00        20 seconds ago      451MB
golang              alpine              ffe224d301bc        4 days ago          311MB
```
When using this technique, I can generate an equivalent Docker image that is only `5.59M` in size. That's over **95%** in reductions, compared to Go's standard Docker process.


## The File
Here is the file that I've been preaching.
```
# Gopherfile
# github.com/TheTannerRyan/Gopherfile

FROM golang:alpine as gopherfile
ENV GO111MODULE on

RUN adduser -D -g '' gopher
WORKDIR /data

# certificates + timezone data
RUN apk update
RUN apk --no-cache add ca-certificates tzdata

# dependency management
RUN apk add git

FROM gopherfile as build
COPY . /data

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -ldflags="-w -s" -o /data/entrypoint

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /data /
USER gopher

ENTRYPOINT ["/entrypoint"]
```
Now this isn't doing any *real* magic, but rather is taking advantage of two technical features:
1) Dockers multi-stage builds
2) Removal of additional symbols and debug info from Go binaries

Docker multi-stage builds allow the use of multiple containers for building images. For Gopherfile, multiple containers are used for building the smallest possible binary, while a final scratch image is used for storing the binary.

During the compilation process, the standard `go build main.go` isn't being used. The large build script is performing cross compilation while stripping all unrequired objects from the binary. This allows for the smallest image (without modifying the existing program).

## Usage
With the current setup, Go modules are being used as the dependency manager. If you're using other tools such as `dep`, you will need to modify the file.

The standard build process is still used. Running something like this should work (again, modify as required):
```
cd src
go mod vendor
docker build -t gopherimg .
docker run -d --restart=unless-stopped -p 80:3000 --name server gopherimg
```

## Questions + Issues
If you have any questions or issues, you know what to do.
