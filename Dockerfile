FROM scratch
ADD main /

ENV GIN_MODE release
CMD ["/main"]
