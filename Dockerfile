FROM scratch
COPY tfc-cli /
ENTRYPOINT ["/tfc-cli"]
