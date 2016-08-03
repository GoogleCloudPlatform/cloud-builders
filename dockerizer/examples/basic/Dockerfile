FROM alpine

RUN echo 'echo "Hello, world!"' > /entry.ash && chmod +x /entry.ash

ENTRYPOINT ["ash", "/entry.ash"]
