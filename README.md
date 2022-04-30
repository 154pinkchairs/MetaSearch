# Description
A meta search engine I wrote because I was unhappy with the enormity of searx.

# Config
A default configuration is given in `config.go.def`.
Simply copy this into `config.go`, which is ignored by git, and make whatever changes you desire.

# TLS Certificate
The `genTestCertKey.sh` can be used to generate a self signed certificate and key.
However these should only be used for testing.
You can obtain a legitimate TLS certificate from [Let's Encrypt](https://letsencrypt.org/) to use when deploying the server.

# Ranking Score
The score used to rank search results comes from [this](https://doi.org/10.1145/1571941.1572114) paper.
