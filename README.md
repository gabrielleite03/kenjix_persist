# Kenjix Persist

Small Go service for Kenjix persistence layer.

## Requirements

- Go 1.24+
- Docker (optional)

## Build & Run (local)

Run with `go run` from repository root:

```bash
go run ./cmd/api
```

Build binary:

```bash
go build -o bin/kenjix ./cmd/api
./bin/kenjix    # or .\bin\kenjix.exe on Windows
```

The server listens on `PORT` (default `8080`). Health endpoint: `/health`.

## Docker

Build image (run from repo root where `go.mod` lives):

```bash
# run from repository root
docker build -t kenjix:latest -f docker/Dockerfile .
docker run --rm -p 8080:8080 kenjix:latest
```

If you run `docker build` from inside the `docker/` folder, pass `..` as context:

```bash
cd docker
docker build -t kenjix:latest -f Dockerfile ..
```

## Docker Compose

Start (from repo root):

```bash
docker compose -f docker/docker-compose.yml up -d --build
```

Stop and remove:

```bash
docker compose -f docker/docker-compose.yml down -v
```

## VS Code Launch

The project includes a `.vscode/launch.json` configured to launch `cmd/api`.

## Generating `go.sum`

If `go.sum` is missing, run from repo root:

```bash
go mod tidy
```

1️⃣ Visão geral (modelo conceitual)

Entidades principais:

Product → o que você vende

Category → organização dos produtos

Warehouse (ou Location) → onde o estoque fica

Stock → quantidade atual por produto/local

StockMovement → histórico de entradas e saídas

Supplier → quem fornece os produtos

PurchaseOrder → compras (entrada)

SalesOrder → vendas (saída)

1️⃣ Onde entram os custos de importação no fluxo

Fluxo real:

Fornecedor internacional
   ↓
Pedido de compra internacional
   ↓
Custos de importação (frete, imposto, seguro…)
   ↓
Rateio dos custos nos produtos
   ↓
Custo médio do estoque


👉 Custos de importação NÃO são do produto,
eles pertencem ao processo de compra/importação.


1️⃣ Onde entram os custos de importação no fluxo

Fluxo real:

Fornecedor internacional
   ↓
Pedido de compra internacional
   ↓
Custos de importação (frete, imposto, seguro…)
   ↓
Rateio dos custos nos produtos
   ↓
Custo médio do estoque


👉 Custos de importação NÃO são do produto,
eles pertencem ao processo de compra/importação.

🟦 ImportCost

Cada custo individual da importação.

ImportCost
- id
- import_process_id
- type
- description
- amount
- currency


Exemplos de type:

FREIGHT

INSURANCE

CUSTOMS_TAX

STORAGE

BROKER_FEE

OTHER

🟦 ImportCostAllocation

Rateio do custo por produto.

ImportCostAllocation
- import_cost_id
- product_id
- allocated_amount

Incoterm (International Commercial Terms) são regras internacionais que definem quem é responsável por quê numa operação de compra e venda internacional.

Em uma frase:
👉 Incoterms dizem até onde vai a responsabilidade do vendedor e a partir de onde começa a do comprador.

🔍 O que exatamente um Incoterm define?

Ele deixa claro, de forma padronizada:

📦 Quem paga o frete

🛃 Quem cuida do desembaraço aduaneiro

🛡️ Quem contrata o seguro

⚠️ Quando o risco da mercadoria é transferido

🚚 Até onde o vendedor entrega

Importante: Incoterm não trata de pagamento, propriedade da mercadoria ou contrato comercial — só logística, custos e riscos.

🌍 Quem criou?

A Câmara de Comércio Internacional (ICC).
A versão mais usada hoje é a Incoterms 2020.

📦 Exemplos comuns (na prática)
🔹 EXW – Ex Works

Vendedor: só disponibiliza o produto

Comprador: faz tudo

👉 Muito usado em importações da China

🔹 FOB – Free On Board

Vendedor entrega no porto de embarque

Comprador assume dali em diante

👉 Clássico para transporte marítimo

🔹 CIF – Cost, Insurance and Freight

Vendedor paga frete + seguro até o destino

Risco ainda transfere no embarque

👉 Muito comum no Brasil

🔹 DDP – Delivered Duty Paid

Vendedor entrega na porta do comprador, com impostos pagos

Comprador quase não se envolve

👉 Raro, mas simples para quem compra



## Notes

- Ensure you run Docker builds with the repository root as the build context so `go.mod` is included.
- Avoid `container_name` in `docker-compose.yml` if you plan to scale services.