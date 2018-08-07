FROM golang
WORKDIR /
RUN mkdir -p /u01/app/Oracle

# need to remove this later on
RUN apt-get update && apt-get install unzip

# need to remove this later on
COPY instantclient-basiclite-linux.x64-12.2.0.1.0.zip .
RUN unzip instantclient-basiclite-linux.x64-12.2.0.1.0.zip -d /u01/app/Oracle
RUN ls -lrt /u01/app/Oracle/instantclient_12_2

# need to remove this later on
COPY instantclient-sdk-linux.x64-12.2.0.1.0.zip .
RUN unzip instantclient-sdk-linux.x64-12.2.0.1.0.zip -d /u01/app/Oracle
RUN ls -lrt /u01/app/Oracle/instantclient_12_2

COPY oci8.pc /u01/app/Oracle/instantclient_12_2
RUN ls -lrt /u01/app/Oracle/instantclient_12_2

# to be removed
COPY libaio1_0.3.110-1_amd64.deb .
RUN dpkg-deb -R libaio1_0.3.110-1_amd64.deb /u01/app

ENV PKG_CONFIG_PATH /u01/app/Oracle/instantclient_12_2
ENV LD_LIBRARY_PATH /u01/app/Oracle/instantclient_12_2:/u01/app/lib/x86_64-linux-gnu

RUN ln -s -f /u01/app/Oracle/instantclient_12_2/libclntsh.so.12.1 /u01/app/Oracle/instantclient_12_2/libclntsh.so
RUN ln -s -f /u01/app/Oracle/instantclient_12_2/libocci.so.12.1 /u01/app/Oracle/instantclient_12_2/libocci.so

RUN ls -lrt /u01/app/Oracle/instantclient_12_2
RUN echo $PKG_CONFIG_PATH
RUN echo $LD_LIBRARY_PATH

RUN go get github.com/mattn/go-oci8

RUN rm instantclient-basiclite-linux.x64-12.2.0.1.0.zip
RUN rm instantclient-sdk-linux.x64-12.2.0.1.0.zip
RUN rm libaio1_0.3.110-1_amd64.deb