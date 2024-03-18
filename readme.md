Stori Card Challenge
Este proyecto es una aplicación que procesa un archivo CSV de transacciones de débito y crédito, guarda los datos en una base de datos PostgreSQL y envía un resumen por correo electrónico. Está desarrollado en Go y utiliza Docker para su ejecución.

Estructura del proyecto
````
stori-card-challenge/
├── README.md
├── cmd/
│   └── main.go
├── docker-compose.yml
├── readme.md
├── env-example
├── .gitignore
├── go.mod
├── go.sum
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── email/
│   │   ├── email.go
│   │   ├── logo.png
│   │   ├── styles.css
│   │   └── template.html
│   ├── file/
│   │   └── file.go
│   └── transaction/
│       └── transaction.go
└── transactions.csv
````
