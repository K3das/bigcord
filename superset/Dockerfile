FROM apache/superset

USER root
COPY setup.sh /app/custom-setup.sh
RUN chmod +x /app/custom-setup.sh

RUN mkdir /data && chown -R superset:superset /data

USER superset

ENV PIP_NO_CACHE_DIR=1
RUN pip3 install clickhouse-connect

COPY config.py /app/custom_config.py
ENV SUPERSET_CONFIG_PATH /app/custom_config.py
ENV SQLALCHEMY_DATABASE_URI="sqlite:////data/superset.db"

ENTRYPOINT [ "/app/custom-setup.sh" ]
