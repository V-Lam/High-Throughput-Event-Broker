# High-Throughput-Event-Broker
Event broker to handle bottlenecks. A batching worker queue to ease the database pressure. Essentialy a homemade Apache Kafka/RabbitMQ/AWS SQS


Goals:
Throttle incoming traffic.
Efficient SQL to handle heavy loads 
Safe, concurrent Go code that doesn't suffer from race conditions.





+------------------+      HTTP/gRPC      +-----------------------+      SQL Queries      +--------------------+

|  1. The Producer | ------------------> | 2. The Broker (Go)    | --------------------> | 3. The Database    |
| (Python Load)    |                     | In-Memory Engine      |                       | (PostgreSQL Log)   |
+------------------+                     +-----------------------+                       +--------------------+
                                                     |
                                                     | Internal Channels
                                                     v
                                         +-----------------------+

                                         | 4. Worker Pool        |
                                         | (Processes & Logs)    |
                                         +-----------------------+





1. The Producer
Python script to mimic thousands of requests slamming the server. It will generate fake data (e.g., "User 123 clicked a button")


2. The Broker Engine (The Core Go App)
Receives the massive stream of data from the Producer. 
It can’t process everything instantly, so it uses Go Channels to queue the messages safely in memory.
It manages a Worker Pool (a designated team of background goroutines) that grabs messages off the queue one by one.


3. The Database (The Safe Storage)
Go workers take the messages and write them into the PostgreSQL database.

TODO batch them together (e.g., saving 100 messages in one single database trip instead of 100 separate trips). Because thousands of writes can slow down Postgres








bash terminal commands


docker compose up -d

docker compose down



docker exec -i backend_practice_db psql -U admin -d event_broker < schema.sql



go mod init github.com/V-Lam/High-Throughput-Event-Broker
go get github.com/jackc/pgx/v5


go run cmd/server/main.go