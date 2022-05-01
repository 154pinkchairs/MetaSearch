# Description
A meta search engine I wrote because I was unhappy with the enormity of searx.

# Config
A default configuration is given in `config.go.def`.
Simply copy this into `config.go`, which is ignored by git, and make whatever changes you desire.

This file also contains the custom listen and serve command.
If you would like to serve the site in a differ manner, using https instead of http perhaps, change this function.

# Other Def Files
The other files ending with `.def` are also meant to be modified.
Simply copy them into a file with the `.def` suffix removed and make whatever changes desiredd.

`opensearch.xml.def` contains a default [OpenSearch](https://developer.mozilla.org/en-US/docs/Web/OpenSearch) description.
Modify this with the appropriate information to allow your instance of MetaSearch to be added as a search engine to your browser.

`index.html.def` and `results.html.template.def` contain the default html which will be served to the user, modify these to change the look of the site.

# Ranking Score
The score used to rank search results comes from [this](https://doi.org/10.1145/1571941.1572114) paper.
