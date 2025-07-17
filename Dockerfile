FROM gcr.io/distroless/static:nonroot

COPY semantic-search-api /semantic-search-api 

ENTRYPOINT [ "/semantic-search-api" ]