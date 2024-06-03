#url := flag.String("url", "https://www.google.com/", "URL do serviço a ser testado")
#requests := flag.Int("requests", 50, "Número total de requests")
#concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")

run:
	@echo "Starting the application..."
	@echo "Example:"
	@echo "url: http://google.com"
	@echo "requests: 50"
	@echo "concurrency: 10"
	go run main.go

build:
	@echo "Building Docker Image..."
	docker build -t stress-test .

.PHONY: run build