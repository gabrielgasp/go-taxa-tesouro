# Contribuindo

Obrigado pelo interesse em contribuir com este projeto.

## Desenvolvimento local

1. Copie o arquivo de exemplo de ambiente para `.env`:

```bash
cp .env.example .env
```

2. Ajuste as variáveis conforme o seu ambiente.

3. **Importante**: Seu `.env` deve conter `ENV=production`, para que a aplicação funcione em Docker.

4. Execute a aplicação localmente:

```bash
go run .
```

## Usando Docker

Se preferir rodar a aplicação com Docker:

1. Construa a imagem:

```bash
docker build -t go-taxa-tesouro .
```

2. Execute o container usando seu arquivo `.env`:

```bash
docker run --rm -p 4722:4722 --env-file .env go-taxa-tesouro
```

## Estrutura da configuração

O projeto lê as variáveis do `.env` e também aceita variáveis de ambiente do sistema.

Os valores mais importantes para desenvolvimento são:

- `PORT`: porta da API local;
- `API_KEY`: chave usada nos endpoints protegidos;
- `SCRAPER_ENABLED`: habilita o scraping em background;
- `INVEST_CSV_URL` e `REDEEM_CSV_URL`: URLs dos dados de entrada;
- `INTERVAL_MINUTES`: frequência de atualização;
- `START_DAY`, `END_DAY`, `START_HOUR`, `END_HOUR`: janela em que o scraper pode rodar;
- `SORT_ORDER`: ordem de apresentação dos títulos;
- `LOG_LEVEL`: nível de logs.

## Antes de abrir um PR

- Confira as issues abertas para evitar duplicar trabalho ou perder contexto;
- Verifique se a aplicação sobe com `go run .`;
- Confirme se os endpoints principais continuam respondendo como esperado;
- Atualize a documentação quando adicionar ou alterar comportamento público.