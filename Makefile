publish: build_exporter build_images build_push

deploy: install_prometheus install_redis

update_deps:
	cd ./deploy && helm dep update

install_prometheus:
	helm upgrade --install prom --wait \
	stable/prometheus \
	--set server.ingress.hosts={prometheus} \
	--namespace default \
	-f ./prometheus/custom_values.yml

uninstall_prometheus:
	helm del prom || exit 0

install_redis:
	helm upgrade \
	--install redis \
	--namespace default \
	./deploy

uninstall_redis:
	helm del redis || exit 0

build_exporter:
	cd src-redis-exporter && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

build_images:
	docker build -t mrshadow74/redis-exporter ./src-redis-exporter

build_push:
	docker push mrshadow74/redis-exporter