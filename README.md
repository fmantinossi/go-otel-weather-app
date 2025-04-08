# 🌎 Sistema de Clima por CEP com Observabilidade (Go + OpenTelemetry + Zipkin)

Este projeto demonstra a implementação de um sistema de microserviços em Go, com observabilidade completa utilizando OpenTelemetry e Zipkin para tracing distribuído.

## 📦 Serviços

### 🔹 Service A
- Recebe um input de CEP via POST (formato JSON).
- Valida o CEP (deve conter 8 dígitos numéricos).
- Encaminha a requisição para o **Service B** via HTTP.

### 🔹 Service B
- Consulta o **ViaCEP** para encontrar a cidade a partir do CEP.
- Consulta a **WeatherAPI** para obter a temperatura atual da cidade.
- Responde com a cidade e as temperaturas convertidas para:
  - Celsius
  - Fahrenheit
  - Kelvin

### 🔍 Observabilidade
- Tracing distribuído com OpenTelemetry (OTEL)
- Exportação para Zipkin via OTEL Collector
- Spans para:
  - Chamada entre serviços
  - Requisição ao ViaCEP
  - Requisição à WeatherAPI

---

## 🚀 Executando a aplicação localmente

### ✅ Pré-requisitos

- Docker
- Docker Compose
- Chave de API válida da [WeatherAPI](https://www.weatherapi.com/)

---

### 📁 Clonar o repositório

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
cd seu-repositorio
```

---

### 🔧 Configurar a chave da WeatherAPI

Crie o arquivo .env e inclua a linha:

```
WEATHER_API_KEY=your_api_key_here
```

---

### ▶️ Subir os serviços

```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up
```

---

## 📡 Testando a API

### 🔹 Requisição válida

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```

Resposta esperada:

```json
{
  "city": "Linhares",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

---

### 🔸 CEP inválido (formato incorreto)

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "123"}'
```

Resposta: `422 - invalid zipcode`

---

### 🔸 CEP com formato válido, mas inexistente

```bash
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "00000000"}'
```

Resposta: `404 - can not find zipcode`

---

## 🔍 Acessar os traces no Zipkin

Abra o navegador em:

👉 [http://localhost:9411](http://localhost:9411)

1. Clique em **"Run Query"**
2. Verifique os traces distribuídos entre os serviços
3. Exemplo de fluxo:
```
Service A
 └── Call Service B
      ├── ViaCEP Lookup
      └── WeatherAPI Lookup
```

---

## 🧱 Estrutura de pastas

```
.
├── docker-compose.yml
├── otel/
│   └── otel-collector-config.yaml
├── service-a/
│   ├── handlers/
│   ├── routes/
│   ├── otel/
│   ├── main.go
│   └── go.mod
└── service-b/
    ├── handlers/
    ├── services/
    ├── clients/
    ├── otel/
    ├── main.go
    └── go.mod
```

---

## 🔧 Execução manual (sem Docker)

Você pode rodar os serviços localmente usando Go:

### 1. Inicie os serviços:

```bash
cd service-a && go run main.go
# Em outro terminal
cd service-b && go run main.go
```

### 2. Exporte as variáveis de ambiente:

```bash
export SERVICE_B_URL=http://localhost:8081
export WEATHER_API_KEY=sua api key
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

### 3. Suba o OTEL Collector e o Zipkin com Docker:

```bash
docker-compose up otel-collector zipkin
```

---

## 🛠️ Tecnologias utilizadas

- Go 1.24
- Docker / Docker Compose
- OpenTelemetry SDK
- OTEL Collector
- Zipkin
- WeatherAPI
- ViaCEP

---
