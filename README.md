[![Build Status](https://travis-ci.org/JNPRAutomate/gogoLogs.svg?branch=master)](https://travis-ci.org/JNPRAutomate/gogoLogs)

gogoLogs
========

A tool to read and send log data to a syslog server

Command Line Options
====================
```
Depricated:
  -C=false: Specify if the file should be continously read from (default: false)
  -F=0: Specify syslog priority value (default:0)
  -P="UDP": Specify protocol to send data over (default: UDP)
  -S=0: Specify syslog priority value (default:0)
  -d="127.0.0.1": Specify destination IP (default: 127.0.0.1)
  -f="": Specify file to read from (default: None)
  -p="514": Sepecify port (default: 514)
  -r=5: Specify rate (default: 5/s)
  -s="127.0.0.1": Specify source hostname/IP for syslog header (default: 127.0.0.1)
  -w=false: Enable WebIU for log sender (default: false)
  -wD="": Specify the directory to serve logs from (default: None, used with -w)
  -wP=8080: Specify the port to listen on (default: 8080) (Used with -w)

New:
  -d="": Optional: Specify the destination hosts to prepopulate this field in the UI. Use comma to provide multiple values. (Example: 1.2.3.4 or 1.2.3.4,2.3.4.5)
  -s="": Optional: Specify the source hosts to prepopulate this field in the UI. Use comma to provide multiple values. (Example: foo or foo,bar)
  -wD="": Specify the directory to serve logs from (default: None, used with -w)
  -wP=8080: Specify the port to listen on (default: 8080) (Used with -w)
```

Web Mode
=========
```
  NOTE: Web mode is the default and the -w flag is no longer required

  To enable web mode you need to enable the WebUI and then specify a directory to read from

  gogoLogs -wD /var/log

  You can also specify the log destinatio hosts and source host names in the UI. The flags are option and can be used together or alone.

  This is not required but it may help the end user by predefining these fields.

  gogoLogs -wD /var/log/ -d 10.0.1.100,10.0.2.2 -s LogHost1,MrLogger
```
