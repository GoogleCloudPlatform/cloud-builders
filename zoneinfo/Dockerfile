FROM alpine:3 as zoneinfo
RUN apk --no-cache add ca-certificates
# ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip .
# 2020f
ADD https://github.com/golang/go/tree/20d0991b86f533a734cad96b2002678d9750e6d8/lib/time/zoneinfo.zip .

FROM scratch
# https://golang.org/pkg/time/#LoadLocation
# https://golang.org/src/time/zoneinfo_unix.go
# Many systems use /usr/share/zoneinfo, Solaris 2 has
# /usr/share/lib/zoneinfo, IRIX 6 has /usr/lib/locale/TZ.
ENV ZONEINFO /usr/share/zoneinfo

# Copy the time zone database.
COPY --from=zoneinfo /zoneinfo.zip /usr/share/zoneinfo
