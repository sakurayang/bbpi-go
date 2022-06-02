FROM Alpine:3.14
RUN mkdir ~/pi
WORKDIR ~/pi
COPY dist/main ~/pi/
ENTRYPOINT ["main"]
