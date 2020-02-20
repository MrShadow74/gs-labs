## GS Labs тестовое задание "Экспортер ключей из базы Redis"

```
.
├── deploy
│   ├── Chart.lock
│   ├── charts
│   │   ├── redis-10.5.3.tgz
│   │   └── redis-exporter-1.0.0.tgz
│   ├── Chart.yaml
│   └── values.yaml
├── Makefile
├── prometheus
│   └── custom_values.yml
├── README.md
├── redis-exporter
│   ├── Chart.yaml
│   ├── templates
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── values.yaml
└── src-redis-exporter
    ├── docker-compose.yml
    ├── Dockerfile
    └── main.go
```

Экспортер выполнен на языке Go, для реализации функционала использованы следующие внешние библиотеки:
 - github.com/go-redis/redis
 - github.com/prometheus/client_golang/prometheus
 - github.com/prometheus/client_golang/prometheus/promhttp


Логика работы приложения:
 - объявляются и инициализируются метрики с именами PointA, PointB, PointC;
 - объявляются переменные для временного хранения значений a, b, c;
 - производится подключение к серверу Redis с адресом redis-master по порту 6379, в случае неудачного подключения перезапуск;
 - в базе Redis создаются ключи Point-A, Point-B, Point-C, каждому присваивается уникальное значение;
 - объявляются переменные val1, val2, val3, которым присваиваются соответствующие значения ключей, полученные из базы Redis;
 - в случае успешного получения значения ключа для ключа присваивается значение 0;
 - выполняется преобразование полученных значений из string в float64, результат присваивается метрикам;
 - запускается веб-сервер, который отдаёт метрики в стандартный /metrics на порт 2112.


Для приложения создан кластер Kubernetes в облаке GCP, состоящий из 3 нод g1-small.
Для удобства создан `Makefile`, в котором можно выполнить все необходимые операции для разворачивания приложений из helm-чартов:
 - для разворачивания Redis и экспортера создан helm-чарт в каталоге deploy;
 - для разворачивания Prometheus используется helm-чарт stable/prometheus с кастомным конфигурационным файлом.

```
kubectl get all
NAME                                          READY   STATUS    RESTARTS   AGE
pod/prom-prometheus-server-78c559c9c6-n6j2d   2/2     Running   0          5h39m
pod/redis-exporter-75fbfddbcd-6gjzs           1/1     Running   3          70m
pod/redis-master-0                            1/1     Running   0          70m

NAME                             TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
service/kubernetes               ClusterIP      10.99.0.1      <none>           443/TCP        9h
service/prom-prometheus-server   LoadBalancer   10.99.1.165    35.246.250.120   80:30048/TCP   5h39m
service/redis-exporter           ClusterIP      10.99.15.229   <none>           2112/TCP       70m
service/redis-headless           ClusterIP      None           <none>           6379/TCP       70m
service/redis-master             ClusterIP      10.99.12.154   <none>           6379/TCP       70m

NAME                                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/prom-prometheus-server   1/1     1            1           5h39m
deployment.apps/redis-exporter           1/1     1            1           70m

NAME                                                DESIRED   CURRENT   READY   AGE
replicaset.apps/prom-prometheus-server-78c559c9c6   1         1         1       5h39m
replicaset.apps/redis-exporter-75fbfddbcd           1         1         1       70m

NAME                            READY   AGE
statefulset.apps/redis-master   1/1     70m
```