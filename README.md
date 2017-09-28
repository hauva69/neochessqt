# neochess
Chess Database System written in Go using Qt

![Early Prototype](/screenshots/EarlyProtoType.png?raw=true "Early Prototype")

## Key Fatures and Progress

- [x] Setup Github Environment 
- [x] Commit early pilot assets
- [ ] Refactor current code prior to commiting to repository
- [ ] Commit all of the current application code
- [x] Initial Move Gen Engine for use with Game Board View
- [ ] Add Support for UCI engines
- [ ] Create test suite for Move Gen Engine
- [ ] Create initial PGN Editor based on Qt Webengine
- [ ] Document Binary Game Storage Engine
- [ ] Create Import for PGN files to Binary Game Storage
- [ ] Create Import for SCID format to Binary Game Storage
- [ ] Create Import for Chessbase format (at least for known formats) to Binary Game Storage
- [ ] Review Licensing 

## How to build

### Prerequisets

- [Git](https://git-scm.com) 
- [Go](https://golang.org) Currently developing and testing with version 1.9
- [Qt](https://www.qt.io) Qt Framework for your development environment
- [Go Qt Binding](https://github.com/therecipe/qt/) Note therecipe also provides docker images for targeting each of the environments, if you don't want to configure this yourself

### Linux

High level steps for now:

- Install prequisets
- Clone this repository
- Run qtdeploy in directory

### Windows

- Install prequisets
- Clone this repository
- Run qtdeploy in directory

### OS X

- Install prequisets
- Clone this repository
- Run qtdeploy in directory

## Credits

- [Qt Framework](https://www.qt.io/) Qt Framework
- [therecipe/qt](https://github.com/therecipe/qt/) Binding Library used to develop Go GUI applications utilizing the Qt framework.
- [Qt Styling](https://github.com/ColinDuquesnoy/QDarkStyleSheet) Base style sheet with modifications

## Inspiring Chess Related Go Project

- [Donna](https://github.com/michaeldv/donna)

