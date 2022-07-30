> Instrumentation documentation

No endpoint `/metrics` a API expoem metricas no formato Prometheus para que um certo **scraper** venha realizar a coleta das mesmas.

Atualmente a instrumentação da aplicação provê visões de regras de negócio que serão listadas abaixo e visões **default** vindas da própria biblioteca do Prometheus:

```yaml
- Tempo em que a API levou para criar/deletar uma publicação/usuario e quantidade total:
    Nome: sm_time_to_create_publication
    Descricao: Tempo que levou para criar uma publicaco
    Tipo: Histogram

    Nome: sm_time_to_create_user
    Descricao: Tempo que levou para criar um usuario
    Tipo: Histogram

    Nome: sm_created_publications_total
    Descricao: Total de publicaoes criadas
    Tipo: Counter

    Nome: sm_created_users_total
    Descricao: Total de usuarios criadas
    Tipo: Counter

    Nome: sm_deleted_publications_total
    Descricao: Total de publicaoes deletadas
    Tipo: Counter

    Nome: sm_deleted_users_total
    Descricao: Total de usuarios deletados
    Tipo: Counter

- Quantidade de conexões TCP abertas:
    Nome: sm_open_connections
    Descricao: Numero total de conexoes abertos
    Tipo: Gauge

- A duracao dos requests para a api:
    Nome: sm_handlers_duration_seconds
    Descricao: Duracao dos requests em segundos
    Tipo: Histogram

- Numero total de requests feitos para a API por status:
    Nome: sm_requests_total
    Descricao: Numero total de requests por status que a api processou
    Tipo: Counter

- Numero total de requests com erro:
    Nome: sm_errors
    Descricao: Numero total de requests que deram erro api processou
    Tipo: Counter
```