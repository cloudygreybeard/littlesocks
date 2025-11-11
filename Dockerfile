FROM scratch

COPY littlesocks /littlesocks

EXPOSE 1080

ENTRYPOINT ["/littlesocks"]
CMD ["-addr", "0.0.0.0:1080"]

