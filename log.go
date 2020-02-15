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
        "MaxLines": "10K",
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


