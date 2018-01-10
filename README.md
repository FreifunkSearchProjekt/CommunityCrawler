# CommunityCrawler

The crawler deployed at the communities to fetch local pages

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

1. Get golang via your package manager or follow instructions at https://golang.org/dl/
2. Download and install [gb](https://getgb.io/) `go get github.com/constabulary/gb/...`

### Installing

A step by step series of examples that tell you have to get a development env running

1. Clone this Repo somewhere on your system

```
git clone https://github.com/FreifunkSearchProjekt/CommunityCrawler.git
```

2. Move into the folder

```
cd CommunityCrawler
```

3. Build the Source got with the help of gb

```
gb build
```

4. Create a config. The default filename is "crawler.yaml"

Example Config:
```
community_id: ffslfl # The Community ID you registered before at the UI page to get the Access Token
community_access_token: null # YOU WILL GET THIS VIA THE UI AT https://not-yet-done.de
indexer:
 - http://riot.nordgedanken.de:9999 # The main indexer URL. Selfhosted indexer can also be added to the array.
network:
 - 10.24.32.0/19 # A valid network. This should match the Community network used as prefix in the site.
# External Pages are pages either in the Internet or not matched by the network array. For example the Website of your Community. You can either add the Page or the Feed of the Page as URL.
external_pages:
 - https://schleswig-flensburg.freifunk.net/feed/
 - https://schleswig-flensburg.freifunk.net
```

5. Run it

```
bin/main-crawler
```

## Running the tests

This Project currently doesn't have working Tests

## Deployment

This should run on all Platforms and doesn't require either much CPU or RAM.
Also this is a oneshot script. You should run this roughly one time a month via a Crontab

## Built With

* [go](https://golang.com/) - The programming language used
* [gb](https://getgb.io/) - Dependency Management and Project Buildsystem
* Plenty more... See the vendors folder

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/FreifunkSearchProjekt/CommunityCrawler/tags).

## Authors

* **Marcel Radzio** - *Initial work* - [MTRNord](https://github.com/MTRNord)

See also the list of [contributors](https://github.com/FreifunkSearchProjekt/CommunityCrawler/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* None yet. Be special and you might get listed! :D
