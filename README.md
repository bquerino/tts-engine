# 🗣️ TTS Engine

TTS Engine is a **text-to-speech (TTS) API** built with **Golang**, using **AWS Polly** and designed for integration with **Twitch** and other platforms. It supports **multiple TTS providers** and offers **load balancing** to distribute requests efficiently.

## 🚀 Features
- 🗣️ **Text-to-Speech API** using AWS Polly.
- 🎯 **Supports multiple TTS providers** (AWS Polly, Google TTS, Azure TTS - expandable).
- 🔄 **Load balancing** to distribute usage between providers.
- 🗂️ **Persists usage statistics** using Redis.
- 📊 **Observability & Tracing** with OpenTelemetry.
- 📦 **Dockerized for easy deployment**.

---

## 🛠️ Getting Started

### 1️⃣ **Clone the Repository**
```sh
git clone https://github.com/bquerino/tts-engine.git
cd tts-engine
```
### 2️⃣ Set Up Environment Variables
Create a .env file in the root directory:

```sh
AWS_ACCESS_KEY_ID=your_aws_key
AWS_SECRET_ACCESS_KEY=your_aws_secret
AWS_REGION=us-east-1
REDIS_HOST=localhost
REDIS_PORT=6379
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
LOG_FORMAT=json
```

## 🏗️ Running Locally

### 3️⃣ Run with Docker
```sh
docker-compose up --build
```

This will start:

* TTS Engine API
* Redis
* OpenTelemetry Collector

### 4️⃣ Run Manually (Without Docker)
Install dependencies:

```sh
go mod tidy
```

Start Redis (if not using Docker):
```sh
redis-server
```

Run the app:
```sh
go run main.go
```

## 🔥 Usage
### 🎤 Synthesize Speech

**Endpoint**: POST /synthesize
**Example Request**:

```json
{
  "text": "Hello, world!",
  "language": "en-US",
  "voice": "Joanna"
}
```
**Example Response**:

```json
{
  "audio_url": "/tmp/Joanna.mp3"
}
```

## 🛠️ Technologies Used
* Golang - Core language
* Fiber - Web framework
* AWS Polly - TTS provider
* Redis - Caching & usage tracking
* OpenTelemetry - Observability & Tracing
* Docker & Docker Compose - Containerization

## 🤝 Contributing
Pull requests are welcome! Open an issue if you find a bug or have an improvement idea.

## 📄 License
This project is licensed under the MIT License.