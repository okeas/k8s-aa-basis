run:
	GOOS=linux GOARCH=amd64 go build -o myapi main.go
	cp myapi /Users/zx/Desktop/vagrant-k8s/project/
	kubectl delete -f yamls/deploy.yaml
	kubectl apply -f yamls/deploy.yaml
	rm -rf ./myapi
