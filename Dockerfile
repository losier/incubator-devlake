# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM mericodev/lake-builder:0.0.4 as builder

# docker build --build-arg GOPROXY=https://goproxy.io,direct -t mericodev/lake .
ARG GOPROXY=
# docker build --build-arg HTTPS_PROXY=http://localhost:4780 -t mericodev/lake .
ARG HTTP_PROXY=
ARG HTTPS_PROXY=

WORKDIR /app
COPY . /app
ENV GOBIN=/app/bin

RUN make clean && make all

FROM --platform=linux/amd64 alpine:3.15
RUN apk add --no-cache musl-dev libgit2-dev libffi-dev \
    && apk add --no-cache gcc

ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache python3-dev && ln -sf python3 /usr/bin/python
RUN python3 -m ensurepip
RUN pip3 install --no-cache --upgrade pip setuptools
RUN pip3 install dbt-mysql
RUN apk add --no-cache tar

EXPOSE 8080

WORKDIR /app

COPY --from=builder /app/bin /app/bin
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV PATH="/app/bin:${PATH}"

CMD ["lake"]

# Notes: Docker for Mac(M1) sets up qemu emulation, you can try to use the amd64 image by adding the --platform=linux/amd64 flag. 
# Such as: FROM --platform=linux/amd64 alpine:3.15
