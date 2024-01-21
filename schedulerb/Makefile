build:
	go build 
	docker build -t ghcr.io/liangyuanpeng/karmada-schedulerb:latest .
	rm -f schedulerb
	docker push ghcr.io/liangyuanpeng/karmada-schedulerb:latest