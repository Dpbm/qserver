FROM python:3.12.10-alpine3.21

COPY . /
RUN chmod +x ./job.sh

# run job.sh every day at 12PM
RUN echo "0  12  *  *  *    /job.sh" > /etc/crontabs/root

RUN crontab /etc/crontabs/root

CMD ["crond", "-f"]



