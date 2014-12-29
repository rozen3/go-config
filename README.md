go-config
=========

1. package rdcfg is designed to read and parse configures

2. the config file must follow FORMAT 1 or FORMAT 2

###FORMAT 1
  ip =   127.0.0.1
  port=1234


###FORMAT 2:
   ip = 127.0.0.1

  port  =  7890

  [broker]
    listen_port = 7777
  time = 10

  [logger]
  listen_port = 1888

Usage:
=========================================
* use cfg, err := NewRDCFG(configpath) to create a RDCFG obj
* use value, err := cfg.GetDefault("ip") to get string value in form1 (or use cfg.Get("", "ip") intead)
* use value, err := cfg.Get("broker", "listen_port") to get string value in form2
*
* other functions for use: GetInt、GetFloat64、 GetIntDefault、GetFloat64、Default
