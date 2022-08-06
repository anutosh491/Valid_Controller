FROM alpine

COPY Valid_Controller /Valid_Controller

ENTRYPOINT [ "./Valid_Controller" ]
