# [Pós GoExpert - FullCycle](https://fullcycle.com.br)

## Desafios técnicos - Stress test

### Pré-requisitos
- [Golang](https://golang.org/)

### Como executar a aplicação

```bash
  # 1 - Clonar o repositório do projeto
  git clone https://github.com/andersonigorf/goexpert-stress-test.git
  
  # 2 - Acessar o diretório do projeto
  cd goexpert-stress-test

  # 3 - Executar a aplicação local
  go run main.go --url=http://google.com --requests=1000 --concurrency=10
  
  # 4 - Executar a aplicação por meio de container Docker
  make build
  
  ou
  
  docker build -t stress-test .
  
  docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```