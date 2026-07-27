# go-taxa-tesouro

O objetivo desse projeto é realizar o scraping de dados da API do Tesouro Direto através de um script que simula o acesso por um navegador web, disponibilizando os dados em um formato estruturado através de uma API própria.

A necessidade disso surgiu por conta do bloqueio do acesso direto à API do Tesouro Direto, que não permite requisições provenientes de outras aplicações ou scripts.

## Como funciona?

A aplicação sobe dois componentes principais:

1. Um scraper que busca os CSVs de investimento e resgate no Tesouro Direto, interpreta os dados e atualiza um cache em memória.
2. Uma API HTTP que lê esse cache e expõe os endpoints de consulta.

O scraper roda em segundo plano quando `SCRAPER_ENABLED=true`. Ele executa o primeiro carregamento ao iniciar o processo e depois repete a coleta em intervalos definidos por `INTERVAL_MINUTES`, respeitando a janela de execução configurada por dia e horário.

## Como utilizar a API?

A API oferece:

- `GET /health` para health check;
- `GET /bonds` para listar todos os títulos carregados;
- `GET /bonds/{bondName}` para buscar um título específico;
- `POST /bonds` para inserir manualmente os dados no cache, protegido por `Authorization: Bearer <API_KEY>` (token definido nas variáveis de ambiente).

### Exemplos

Listando os títulos:

```bash
curl http://localhost:4722/bonds
```

Buscando um título específico:

```bash
curl http://localhost:4722/bonds/tesouro-ipca+-2035
```

Se você quiser usar os dados em uma planilha, site ou automação, a forma mais simples é consultar esses endpoints e ler o campo desejado no JSON retornado.

Se o nome do título estiver diferente do esperado, a busca pode não encontrar o item.

## Postman

Há uma coleção pronta em [api.postman_collection.json](api.postman_collection.json) para testar os endpoints principais.

## Contribuindo

Consulte o arquivo [CONTRIBUTING.md](CONTRIBUTING.md) para instruções de desenvolvimento local e fluxo de contribuição.

## Licença

Este projeto está licenciado sob os termos do arquivo [LICENSE](LICENSE).