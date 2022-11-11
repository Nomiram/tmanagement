brokers_cnt = 3
with open("docker-compose.yml",'w') as f:
    f.write('''
version: '3.8'
services:
  db:
    user: "postgres"
    image: postgres:14.1-alpine
    restart: always
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    ports:
      - '5432:5432'
    volumes: 
      - db:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - backend
  redis:
    build: ./redis
    depends_on:
      - db
    networks:
      - backend
    
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - 22181:2181
    networks:
      - backend  
    ''')
    for i in range(brokers_cnt):
        f.write(
    f'''
  kafka_{i}:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_BROKER_ID: {i+1}
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka_{i}:{9092+i},PLAINTEXT_HOST1://kafka:{29092+i}
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST1:PLAINTEXT
    healthcheck:
      test: nc -z localhost {9092+i} || exit -1 
      interval: 15s
      timeout: 10s
      retries: 5  
    networks:
      - backend''')
    f.write('''
  api:
    build: ./tmanagement
    environment:
      DBADDR: "db"
    expose:
      # - 5432
      - 8080
    ports:
      - 8080:8080
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
      serv2:
        condition: service_started
    networks:
      - backend
  serv2:
    build: ./durationCount
    environment:
      DBADDR: "db"
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    networks:
      - backend
volumes:
  db:
    driver: local
networks:
  backend:
    driver: 
      bridge
    ''')