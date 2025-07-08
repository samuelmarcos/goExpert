# Order API

## Descrição
Esta API permite criar e listar ordens, contendo os campos **Price**, **Tax** e **Final Price**. O projeto é escrito em Go e expõe endpoints HTTP, gRPC e GraphQL para integração.

## Funcionalidades
- **Criar uma ordem** com os campos:
  - `Price` (float): valor do pedido
  - `Tax` (float): imposto
  - `Final Price` (float): valor final (calculado)
- **Listar todas as ordens** cadastradas

## Tecnologias
- **Linguagem:** Go (Golang)
- **Endpoints:**
  - HTTP REST
  - gRPC
  - GraphQL

## Endpoints

### HTTP
- **Criar ordem:**
  - `POST /order`
  - Exemplo de payload:
    ```json
    {
      "id": "a",
      "price": 100.5,
      "tax": 0.5
    }
    ```
- **Listar ordens:**
  - `GET /orders`

### gRPC
- Serviço: `OrderService`
- Métodos:
  - `CreateOrder(CreateOrderRequest) returns (CreateOrderResponse)`
  - `ListOrders(ListOrdersRequest) returns (ListOrdersResponse)`
- Mensagens definidas em: `internal/infra/grpc/protofiles/order.proto`

### GraphQL
- **Criar ordem:**
  - Mutation: `createOrder(input: OrderInput): Order`
- **Listar ordens:**
  - Query: `listOrders: [Order!]!`
- Schema: `internal/infra/graph/schema.graphqls`

## Testes de API
- O arquivo para testar os endpoints HTTP está localizado em:
  - `api/create_order.http`
- Você pode usar o VSCode REST Client ou ferramentas como Insomnia/Postman para executar os testes.

---

> Para rodar a aplicação, siga as instruções do seu ambiente Go e utilize os endpoints conforme descrito acima. 