package logs

/**
{
    "console": {
        "enable": true,		// wether output the log
        "level": "FINE"		// log level: FINE, DEBUG, TRACE, INFO, WARNING,ERROR, CRITICAL
    },
    "files": [{
        "enable": true,
        "level": "DEBUG",
        "filename":"./test.log",
        "category": "Test",			// different category log to different files
        "pattern": "[%D %T] [%C] [%L] (%S) %M"	// log output formmat
    },{
        "enable": false,
        "level": "DEBUG",
        "filename":"rotate_test.log",
        "category": "TestRotate",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "rotate": true,				// whether rotate the log
        "maxsize": "500M",
        "maxLines": "10K",
        "daily": true,
        "sanitize": true
    }],
    "sockets": [{
        "enable": false,
        "level": "DEBUG",
        "category": "socket",
        "pattern": "[%D %T] [%C] [%L] (%S) %M",
        "addr": "127.0.0.1:12124",
        "protocol":"udp"
    }]
}
 */
func InitLog(configPath string) {
	// load config file, it's optional
	// or LoadConfiguration("./example.json", "json")
	// config file could be json or xml
	LoadConfiguration(configPath)

	GetLogger("socket", "app").Info("category Test info test ...")
	GetLogger("socket", "app").Info("category Test info test message: %s", "new test msg")
	GetLogger("socket", "app").Debug("category Test debug test ...")

	// Other category not exist, test
	GetLogger("socket", "app").Debug("category Other debug test ...")

	// socket log test
	GetLogger("socket", "app").Debug("category TestSocket debug test ...")

	// original log4go test
	Info("normal info test ...")
	Debug("normal debug test ...")

	Close()
}

