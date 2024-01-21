ARG ARCH=amd64
FROM --platform=linux/${ARCH} gcr.io/distroless/static-debian12@sha256:4a2c1a51ae5e10ec4758a0f981be3ce5d6ac55445828463fce8dff3a355e0b75

ADD karmada-scheduler /usr/local/bin/

#EXPOSE 2379 2380

# Define default command.
CMD ["/usr/local/bin/karmada-scheduler"]
