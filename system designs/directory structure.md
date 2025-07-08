WeTalk/
├── README.md # overview + “how to run”
├── go.mod
├── go.sum
├── .gitignore

├── cmd/
│ └── we-talk/
│ └── main.go # mounts both REST & GraphQL on chi

├── internal/ # import-only-here domain logic
│ ├── config/ # load YAML/ENV
│ │ └── config.go
│ ├── db/ # migrations & DB client
│ │ ├── migrations/
│ │ └── postgres.go
│ │
│ ├── users/ # user service (CRUD + business rules)
│ ├── halls/ # “server” service
│ ├── rooms/ # room service
│ ├── chat/ # message hub + pubsub for subscriptions
│ ├── rest/ # REST handlers & router setup
│ │ ├── router.go
│ │ └── handlers.go # users, halls, rooms, roles…
│ └── graph/ # GraphQL gateway
│ ├── schema.graphql # type Query/Mutation/Subscription
│ ├── generated/ # gqlgen output
│ └── resolver/ # map GraphQL → internal/\* calls

├── web/
│ └── react-app/
│ ├── package.json
│ ├── public/
│ └── src/
│ ├── apolloClient.js # HTTP + WebSocket link to /graphql
│ ├── hooks/
│ ├── pages/
│ └── components/

├── scripts/ # helpers
│ ├── build.sh # build Go + React
│ ├── migrate.sh # run DB migrations
│ └── deploy.sh # rsync/scp build → server

├── configs/ # server-side config & service files
│ ├── we-talk.yaml # your service config
│ ├── nginx/
│ │ └── we-talk.conf # reverse-proxy + static files
│ └── systemd/
│ └── we-talk.service # run your Go binary as a daemon

└── data/ # on-server DB files, migrations, etc.
