# khystrix
Playing with hystrix under kubernetes

# Tomcat 7
https://wiki.archlinux.org/index.php/tomcat

# Turbine
https://github.com/Netflix/Turbine/wiki/Getting-Started-(1.x)
https://github.com/Netflix/Turbine/wiki/Configuration-(1.x)#turbine-cluster-configuration

My configuration for localhost:
```
turbine.aggregator.clusterConfig=david
turbine.ConfigPropertyBasedDiscovery.david.instances=localhost
turbine.instanceUrlSuffix.david=:8222/hystrix.stream
```

# Dashboard
https://github.com/Netflix/Hystrix/wiki/Dashboard

Using Turbine:
```
http://localhost:8080/turbine-web/turbine.stream?cluster=david
```

Using App:
```
http://localhost:8222/
```

# Hystrix-go output and turbine

As soon as the event-stream contains both *HystrixCommand* and *HystrixThreadPool* Turbine is getting confused.
So I decide to continue to remove the *HystrixThreadPool*.

To do so I had to comment one line in Hystrix-go:
```
diff --git a/hystrix/eventstream.go b/hystrix/eventstream.go
index 0b49de5..8ed94b3 100644
--- a/hystrix/eventstream.go
+++ b/hystrix/eventstream.go
@@ -78,7 +78,7 @@ func (sh *StreamHandler) loop() {
                        circuitBreakersMutex.RLock()
                        for _, cb := range circuitBreakers {
                                sh.publishMetrics(cb)
-                               sh.publishThreadPools(cb.executorPool)
+                               //sh.publishThreadPools(cb.executorPool)
                        }
                        circuitBreakersMutex.RUnlock()
                case <-sh.done:
```

#Helpers

commands to inject on the process to:

emulate a command call
```
curl -X GET -i "http://localhost:8221/start?name=toto&base=10&floor=13"
```

configure the circuit breaker for a command
```
curl -X GET -i "http://localhost:8221/configure?name=toto&timeout=100&maxConcurrentRequests=100&errorPercentThreshold=10"
```

get the status of the circuit breaker
```
curl -X GET -i "http://localhost:8221/status?name=toto"curl -X GET -i "http://localhost:8221/status?name=toto"
```

force opening the circuit breaker (and reopen)
```
curl -X GET -i "http://localhost:8221/toggleOpen?name=toto&value=false"curl -X GET -i "http://localhost:8221/toggleOpen?name=toto&value=false"

curl -X GET -i "http://localhost:8221/toggleOpen?name=toto&value=false"curl -X GET -i "http://localhost:8221/toggleOpen?name=toto&value=true"
```

close an open circuit (not useful at all...)
```
curl -X GET -i "http://localhost:8221/close?name=toto"
```
