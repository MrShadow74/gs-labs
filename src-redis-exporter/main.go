package main

import (
        "fmt"
        "net/http"
        "strconv"

        "github.com/go-redis/redis"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    PointA = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "Point_A",
        Help: "The current value Point A",
    })

    PointB = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "Point_B",
        Help: "The current value Point B",
    })

    PointC = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "Point_C",
        Help: "The current value Point C",
    })

    a float64
    b float64
    c float64
)

func init() {
    prometheus.MustRegister(PointA)
    prometheus.MustRegister(PointB)
    prometheus.MustRegister(PointC)
}

func main() {
        client := redis.NewClient(&redis.Options{
                Addr: "redis-master:6379",
                Password: "",
                DB: 0,
        })

        pong, err := client.Ping().Result()
        if err != nil {
            fmt.Println("Redis is not Connect")
            panic(err)
        }else {
            fmt.Println(pong, err)
            fmt.Println("Redis is Connected")
        }

        err = client.Set("Point-A", "123", 0).Err()
        if err != nil {
            fmt.Println(err)
        }

        err = client.Set("Point-B", "456", 0).Err()
        if err != nil {
            fmt.Println(err)
        }

        err = client.Set("Point-C", "789", 0).Err()
        if err != nil {
            fmt.Println(err)
        }

        val1, err := client.Get("Point-A").Result()
        if err != nil {
            fmt.Println(err)
        }else{
            client.Set("Point-A", 0, 0)
        }

        val2, err := client.Get("Point-B").Result()
        if err != nil {
            fmt.Println(err)
        }else{
            client.Set("Point-B", 0, 0)
        }

        val3, err := client.Get("Point-C").Result()
        if err != nil {
            fmt.Println(err)
        }else{
            client.Set("Point-C", 0, 0)
        }

        a, err := strconv.ParseFloat(val1,64)
        if err != nil {
            fmt.Println(err)
        }

        PointA.Set(a)

        b, err := strconv.ParseFloat(val2,64)
        if err != nil {
            fmt.Println(err)
        }

        PointB.Set(b)

        c, err := strconv.ParseFloat(val3,64)
        if err != nil {
            fmt.Println(err)
        }

        PointC.Set(c)

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}
