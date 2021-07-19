FROM python:3.9-alpine AS compile-image

COPY requirements.txt .
# Make sure we use the virtualenv:
ENV PATH="/opt/venv/bin:$PATH"

RUN apk update && apk add --no-cache \
    build-base \
    gcc \
    && python -m venv /opt/venv \
    && pip install -r requirements.txt

### === builder image ===
FROM python:3.9-alpine AS build-image

COPY --from=compile-image /opt/venv /opt/venv
# Make sure we use the virtualenv:
ENV PATH="/opt/venv/bin:$PATH"

COPY . /api
WORKDIR /api
RUN chmod u+x run.sh

CMD ['./run.sh']