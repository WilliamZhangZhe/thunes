FROM mysql as mysql
RUN rm -f /var/lib/mysql/* 

EXPOSE 3306 33060

CMD [ "mysqld", "-u", "root" ]