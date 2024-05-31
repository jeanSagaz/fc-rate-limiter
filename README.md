## Dêe uma estreal! :star:
Se vc gostou do projeto 'go-rate-limiter', por favor dêe uma estrela

## Como executar:
Execute no prompt de comando na pasta raiz:  
docker-compose up -d  
Para subir o 'redis' e o 'server' na porta 8080.  

Caso queira rodar local execute o comando:  
go run ./cmd/server/main.go  
Lembrando que deve alterar os arquivos:  
**configs/config.go**  
Comentar a linha 'docker' e descomentar a linha 'local'.  
**cmd/server/.env**  
Comentar a linha 'docker' e descomentar a linha 'local'.  

Instale o 'REST Client' no 'VS Code' e execute os testes da pasta:  
./tests/api.http  

**Configuração do Token:**  
Foi criada uma chave no arquivo .env chamada 'TOKEN_CONFIGURATION' conforme o esquema abaixo:  
[{"Token": "defa69ef-a390-4fad-b319-922e325c9efd", "NumberRequests": 5, "Seconds": 5}  
{"Token": "a71fe9c5-efaf-4267-9bbc-3ccd4dba3b61", "NumberRequests": 5, "Seconds": 7},   
{"Token": "5ebff2ac-e380-4e1d-b32a-9e77b7644ddd", "NumberRequests": 5, "Seconds": 9},...]  

Token: Corresponde ao valor que deve ser informado no 'Header' da requisição  
NumberRequests: Número total de request que pode ser enviado dentro do tempo configurado  
Seconds: Tempo que vai durar a chave no Redis  

**Configuração do IP:**  
Foi criada duas chaves no arquivo .env chamadas 'NUMBER_REQUESTS' e 'SECONDS':  

NUMBER_REQUESTS: Número total de request que pode ser enviado dentro do tempo configurado  
SECONDS: Tempo que vai durar a chave no Redis  

## Tecnologias implementadas:

go 1.22
 - Router [chi](https://github.com/go-chi/chi)
 - Viper
 - DI
 - Redis
 - testing
 - docker
 