# 🚀 Crypto Trading Platform - Technical Overview

A scalable, fault-tolerant, and secure real-time cryptocurrency trading platform. This document outlines the system architecture, strategies used to achieve scalability and resilience, and ongoing plans for future improvements.

---

# 🚀 How to Run

### Prerequisites
- Docker and Docker Compose
- Git
- Make (optional, for convenience commands)

## Clone the repository
   ```bash
   git clone https://github.com/SametAvcii/cyrpto-trade.git
   cd cyrpto-trade 
   ```

## Run Project

#### Running with Docker Compose (Hot Reload Development):

For a smoother development experience with hot reloading, a separate docker-compose-hot.yml configuration is provided. This setup likely mounts your local code into the container, allowing changes to be reflected without rebuilding the entire container. It also utilizes a specific configuration file (config-hot.yml).

```bash
   make hot
   ```
 Alternatively:

```bash
   docker compose -p crypto-trade -f docker-compose-hot.yml up
   ```


#### Running with Docker Compose (Directly):

For a smoother development experience with hot reloading, a separate docker-compose-hot.yml configuration is provided. This setup likely mounts your local code into the container, allowing changes to be reflected without rebuilding the entire container. It also utilizes a specific configuration file (config-hot.yml).

```bash
   make run
   ```
 Alternatively:

```bash
   docker compose -f docker-compose.yml up -d --build
   ```


## 📈 Scalability Approach

### 🧩 Microservices Architecture
- Decoupled **API services** and **consumer services**
- Independent scaling of services based on system load

### ⚙️ Event-Driven Architecture
- **Apache Kafka** as central message broker
- Asynchronous market data processing
- Topic-based message routing (e.g., `orderbooks`, `trades`, `candlesticks`)

### 🗃️ Database Scaling
- **PostgreSQL**: Transactional data with connection pooling
- **MongoDB**: High-volume time-series data (e.g., trades, orderbooks)
- **Redis**: Caching & real-time data access

### 🌐 Horizontal Scaling
- Docker-based containerization
- **Docker Compose** for local development
- Kubernetes-ready stateless service deployment

### ⚡ Performance Optimization
- Database connection pooling
- Redis for frequently accessed data
- Paginated API responses
- Prometheus metrics for performance insights

---

## 🛡️ Fault Tolerance Features

### 🏥 Health Checks & Recovery
- Continuous health monitoring
- Auto-reconnect for PostgreSQL, Redis, and Kafka
- Health check routines:
  - `RedisAlive()`
  - `CheckPgAlive()`
  - `CheckKafkaAlive()`

### 🧬 Data Redundancy
- Critical data persisted in PostgreSQL & MongoDB
- Kafka-based **event sourcing** for data recovery

### 🔧 Graceful Degradation
- Services run with reduced functionality when dependencies fail
- **Circuit breaker** pattern in critical paths

### 🐛 Error Handling & Logging
- `ctlog` for structured and contextual logging
- Comprehensive error metrics with Prometheus

### 🧪 Testing
- Unit tests with mock dependencies


---

## 🔐 Security Implementation

### 🛡️ Data Protection
- Encrypted database connections
- Secrets managed via environment variables
- Secure password hashing & storage

### 🌍 Network Security
- Internal-only service access via Docker networks
- Public endpoint rate limiting

### 📉 Monitoring & Alerting
- Real-time metrics (Prometheus + Grafana)
- Alerting on unusual patterns
- Audit logging for traceability

### 📦 Dependency Security
- Regular updates of packages
- Secure coding and static analysis

---

## 🔧 Challenges Faced & Solutions

### 📈 Moving Average Strategy Evolution

## Initial Approach & Learnings

At the beginning, calculating the moving average (MA) came with several design uncertainties. Initially, we considered using real-time price data (AggTrade events) due to their immediacy. However, this method led to inconsistent results, mainly because:

- AggTrade streams are extremely frequent and lack aggregation

- Price fluctuations and volume shifts made averages unstable

- There was no guarantee of a finalized price per interval

## Strategic Shift to Candlestick Data

To improve signal quality, we transitioned to using candlestick (Kline) data, which aggregates price movements over set timeframes. This change offered:

- Structured OHLCV data (Open, High, Low, Close, Volume)

- Stable and finalized Close prices ideal for MA calculations

- Better alignment with traditional market analysis practices

- This adjustment resulted in more reliable MA values that better reflected actual market conditions instead of momentary changes.

## Real-time vs. Interval Close Debate

One key decision point was whether to use the live price or wait for the candle to close. While the real-time price could generate faster signals, it introduced noise and premature decisions. Ultimately, we chose to:

Wait for the candle to officially close (e.g., every 1m, 5m, etc.)

Then compute and update the MA value

This method improved signal accuracy and followed best practices in algorithmic trading.

## ✅ Redis-Based Moving Average Optimization

To optimize performance:

- We store recent close prices (for MA50 and MA200) in Redis lists per symbol and interval

- Each time a candle closes, we append the price and trim the list

- Redis commands like RPush and LTrim help maintain fixed-length data

## Benefits:

- Quick updates to MA values

- Low-latency signal generation

- Efficient use of memory and simple time windowing

After computation:

- Signals are persisted in PostgreSQL and MongoDB

- The latest signal is also cached in Redis to avoid sending duplicates


### Real-time Data Processing
**Challenge**: Managing high-throughput data from exchanges  
**Solution**:
- Buffered processing pipeline with Kafka
- Specialized consumers for each data stream
- Rate-limited API calls & optimized data structures

### Database Connection Stability
**Challenge**: Handling unstable connections under load  
**Solution**:
- Connection pooling with limits
- Reconnection logic with exponential backoff
- Periodic health checks and batching

### Signal Processing Performance
**Challenge**: Real-time computation of signals  
**Solution**:
- Redis-based sliding window for MAs
- Pre-computed technical indicators
- Goroutine-based parallel processing

### 5️⃣ System Monitoring
**Challenge**: Observability of system internals  
**Solution**:
- Full-stack Prometheus metrics
- Custom Grafana dashboards
- Middleware logging for HTTP requests

---

## 🔭 Future Improvements

- Add **distributed tracing** (OpenTelemetry)
- Integrate **circuit breakers** for external APIs
- Deploy on **Kubernetes** with HPA (horizontal pod autoscaling)
- Introduce **A/B testing framework** for strategy evaluation

---

## ✅ Summary

This platform is engineered to handle high data throughput, resist failures, and uphold security best practices. It’s prepared for scale, extensibility, and rapid evolution in the volatile world of crypto trading.

---

### 📬 Contact
For questions, collaboration, or more technical details, feel free to reach out to the development team.