
FROM alpine:3.19
RUN apk add --no-cache python3 bash busybox-openrc cronie
COPY fake_telnetd.py /
COPY persist.sh /persist.sh
RUN chmod +x /persist.sh
EXPOSE 23
CMD ["python3", "/fake_telnetd.py", "23"]
