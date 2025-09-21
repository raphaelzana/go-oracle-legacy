# Template Go + Oracle XE 11

Template mínimo em Go para conectar ao Oracle (ex.: XE 11g) e expor um endpoint HTTP de verificação.

## Requisitos
- Go 1.24.3+
- Oracle acessível (host, porta 1521, service/SID: XE, usuário e senha)

## Configuração
Defina suas variáveis (exemplo):
```shell script
export ORACLE_USER="<USUARIO>"
export ORACLE_PASS="<SENHA>"
export ORACLE_HOST="<HOST>"
export ORACLE_PORT="1521"
export ORACLE_SERVICE="XE"
```


## Rodando
```shell script
go mod download
go run .
```


## Testando
```shell script
curl -s http://localhost:8080/
```


Resposta esperada (exemplo):
```json
{"message":"Hello, Buddy!","db_time":"YYYY-MM-DD HH:MM:SS"}
```
