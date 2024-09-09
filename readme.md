
# ByteBurger - Sistema Digital de Processamento de Pedidos de Restaurante

Bem-vindo ao ByteBurger, um restaurante digital que utiliza tecnologia moderna para processar pedidos de forma eficiente. Este projeto demonstra como utilizar **Go**, **Web Workers** e **PostgreSQL** para gerenciar e processar vários pedidos de restaurante simultaneamente, lidando com diversas tarefas como cortar ingredientes, grelhar hambúrgueres, montar pratos e preparar bebidas.

## Visão geral do projeto

O ByteBurger simula um sistema utilizado por um restaurante onde os pedidos passam por diversas tarefas até serem concluídos. O sistema utiliza Go para gerenciamento de concorrência e banco de dados, além de PostgreSQL para armazenar e gerenciar os dados dos pedidos.

Cada pedido é composto por diferentes tarefas, como:

- Corte de Ingredientes
- Grelhar Hambúrgueres
- Montagem do Prato
- Preparação de Bebidas

Para lidar de maneira eficiente com vários pedidos, é utilizado o conceito de Web Workers, onde cada worker é responsável por processar uma tarefaa relacionada a um passo do pedido. Além disso, o sistema possui funcionalidades como:

- **Tratamento de Prioridade**: Permite que certos pedidos passem na frente de outros.
- **Cancelamento de Pedidos**: Possibilidade de cancelar um pedido em andamento, interrompendo o processo e removendo as ações seguintes da fila.
- **Estágios do Processamento**: Visualização das etapas de cada pedido, incluindo o progresso de cada worker.
- **Processo de solicitar conta**: Capacidade de retornar, a partir de um ID, os pedidos do cliente e sua respectiv conta.
- **Estimativa de Tempo de Conclusão**: O sistema calcula o tempo estimado para conclusão de cada pedido com base nas tarefas pendentes e no tempo médio de cada etapa do processo.


## Instalação

### Pré-requisitos
- Go 1.20 ou superior

### Passos para Instalação
Clone o repositório e execute o comando:
```bash
git clone https://github.com/seu-usuario/byteburger.git
cd byteburger
```

Para executar o código, basta:
```bash
go run cmd/main.go
```
    