# neochess

Chess Database System written in Go using Qt

![Early Prototype](/screenshots/EarlyProtoType.png?raw=true "Early Prototype")

## Key Features and Progress

- [x] Setup Github Environment
- [x] Commit early pilot assets
- [x] Refactor current code prior to committing to repository
- [x] Commit all of the current application code
- [x] Initial Move Gen Engine for use with Game Board View
- [ ] Integrate and commit a Binary Game Storage
- [ ] Enable Tree View for Recent Databases
- [ ] Integrate and commit a PGN Importer
- [ ] Add Support for UCI engines
- [ ] Create test suite for Move Gen Engine
- [ ] Create initial PGN Editor based on ~~QtWebengine~~ QTextEdit[1]
- [ ] Document Binary Game Storage Engine
- [ ] Create Import for SCID format to Binary Game Storage
- [ ] Create Import for Chessbase format (at least for known formats) to Binary Game Storage
- [ ] Review Licensing

[1]: QtWebengine as far as I know still has a bug either in Qt or in the bindings to golang.
So that it isn't fully supported in windows yet see [Issue](https://github.com/therecipe/qt/issues/217#issuecomment-280940272) 
And since I want NeoChess to build at least in Linux, iOS, and Windows, I'll just work with QTextEdit instead.

## How to build

### Prerequisites

- [Git](https://git-scm.com) 
- [Go](https://golang.org) Currently developing and testing with version 1.9
- [Qt](https://www.qt.io) Qt Framework for your development environment
- [Go Qt Binding](https://github.com/therecipe/qt/) Note therecipe also provides docker images for targeting each of the environments, if you don't want to configure this yourself

#### Go Libraries

Note at some point we might pull copies of these libraries into a vendor directory of this repository.

- [Go Single Instance](https://github.com/allan-simon/go-singleinstance)
  - Used to enforce only a single instance of neochess running
- [Go i18n](https://github.com/nicksnyder/go-i18n)
  - Used for internationalization of text in application
- [Logrus](https://github.com/sirupsen/logrus)
  - Used for application logging

### Linux

High level steps for now:

- Install prerequisites

```bash
$ go get -u github.com/allan-simon/go-singleinstance
$ go get -u github.com/nicksnyder/go-i18n/...
$ go get -u github.com/sirupsen/logrus
```

- Clone this repository

```bash
$ go get github.com/rashwell/neochess
$ cd $GOPATH/src/github.com/rashwell/neochess
```

- Run qtdeploy in directory

```bash
$ $GOPATH/bin/qtdeploy
```

### Windows

- Install prerequisites
- Clone this repository
- Run qtdeploy in directory

### OS X

- Install prerequisites
- Clone this repository
- Run qtdeploy in directory

## Credits

- [Qt Framework](https://www.qt.io/) Qt Framework
- [therecipe/qt](https://github.com/therecipe/qt/) Binding Library used to develop Go GUI applications utilizing the Qt framework.
- [Qt Styling](https://github.com/ColinDuquesnoy/QDarkStyleSheet) Base style sheet with modifications

## Inspiring Chess Related Go Project

- [Donna](https://github.com/michaeldv/donna)

