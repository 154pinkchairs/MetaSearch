# Description
A meta search engine I wrote because I was unhappy with the enormity of searx.

# Config
A default configuration is given in `config.go.def`.
Simply copy this into `config.go`, which is ignored by git, and make whatever changes you desire.

This file also contains the custom listen and serve command.
If you would like to serve the site in a differ manner, using https instead of http perhaps, change this function.

# Open Search
A default open search specification is provided in `opensearch.xml.def`.
As with the config file, simply copy this into `opensearch.xml` and make the appropriate changes for your setup.

# Ranking Score
The score used to rank search results comes from [this](https://doi.org/10.1145/1571941.1572114) paper.
