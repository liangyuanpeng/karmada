build:
	go build 
	docker build -t ghcr.io/liangyuanpeng/karmada-schedulera:latest .
	rm -f schedulera
	docker push ghcr.io/liangyuanpeng/karmada-schedulera:latest