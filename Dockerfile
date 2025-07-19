FROM gcr.io/distroless/static:nonroot

COPY semantic-search-api /semantic-search-api 

ENV LOG_LEVEL=warn

ENTRYPOINT [ "/semantic-search-api" ]