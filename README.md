A NATS hook for Logrus

This lib will let you use any nats connection to connect, but does provide some convience methods to help connection setup.

The helpers all assume you're using a TLS connection.

Examples:

setup and add hook to default logger. Connection is assumed TLS with no error handler.

``` go
  conf := &HookConf{
    Subject: "logs.group.hostname",
    Dimensions: map[string]interface{}{
      "hostname": "server-1",
    },
    NatsConfig: NatsConfig{
      CAFiles: []string{"ca.pem"},
      KeyFile: "key.pem",
      CertFile: "cert.pem",
      Servers: []string{"server-2", "server-3"},
    },
  }
  nc, hook, err := AddNatsHook(conf)
```

setup with an error handler and TLS

``` go
  natsConfig := &NatsConfig {
    CAFiles: []string{"ca.pem"},
    KeyFile: "key.pem",
    CertFile: "cert.pem",
    Servers: []string{"server-2", "server-3"},
  }
  e := func(nc *nats.Conn, sub *nats.Subject, err error) {
    logrus.WithError(err).Fatal("error!!!")
  }
  nc, err := ConnectToNatsWithError(natsConfig, e)

  hook, err := NewNatsHook(nc, "logs.group.host")
  logrus.AddHook(hook)
```

setup with your own connection

``` go
  nc, err := nats.Connect("server-2")
  hook, err := NewNatsHook(nc, "logs.group.host")
  logrus.AddHook(hook)
```

You can modify the formatter on the hook and the level that it is connected to as well. By default it will use the logrus's `JSONFormatter`.

