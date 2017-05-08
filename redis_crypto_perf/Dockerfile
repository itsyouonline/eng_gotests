FROM ubuntu:16.04

# Install Redis.
RUN \
  apt-get update && \
  apt-get install -y wget make gcc && \
  cd /tmp && \
  wget https://github.com/antirez/redis/archive/4.0-rc3.tar.gz && \
  tar xvzf 4.0-rc3.tar.gz && \
  cd redis-4.0-rc3 && \
  make && \
  make install && \
  cp -f src/redis-sentinel /usr/local/bin && \
  mkdir -p /etc/redis && \
  cp -f *.conf /etc/redis && \
  rm -rf /tmp/redis-4.0-rc3* && \
  sed -i 's/^\(bind .*\)$/# \1/' /etc/redis/redis.conf && \
  sed -i 's/^\(daemonize .*\)$/# \1/' /etc/redis/redis.conf && \
  sed -i 's/^\(dir .*\)$/# \1\ndir \/data/' /etc/redis/redis.conf && \
  sed -i 's/^\(logfile .*\)$/# \1/' /etc/redis/redis.conf

# Define mountable directories.
VOLUME ["/data"]

# Define working directory.
WORKDIR /data

# custom config
COPY ./redis.conf /etc/redis/redis.conf

# Define default command.
CMD ["redis-server", "/etc/redis/redis.conf"]

# Expose ports.
EXPOSE 6379
