FROM rabbitmq:3.6

RUN rabbitmq-plugins enable --offline rabbitmq_management
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt

COPY rabbitmq_web_mqtt-3.6.x.ez /plugins
RUN rabbitmq-plugins enable --offline rabbitmq_web_mqtt

# Fix nodename
RUN echo 'NODENAME=rabbit@localhost' > /etc/rabbitmq/rabbitmq-env.conf

EXPOSE 15672
EXPOSE 15675
EXPOSE 1883
