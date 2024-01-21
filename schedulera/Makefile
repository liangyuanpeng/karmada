build:
	go build 
	mv schedulera karmada-scheduler
	docker build -t ghcr.io/liangyuanpeng/karmada-schedulera:latest .
	rm -f karmada-scheduler
	docker push ghcr.io/liangyuanpeng/karmada-schedulera:latest