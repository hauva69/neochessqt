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
- [ ] Create and Setup Installers for Releases for iOS, Windows, Linux
- [ ] Document Binary Game Storage Engine
- [ ] Create Import for SCID format to Binary Game Storage
- [ ] Create Import for Chessbase format (at least for known formats) to Binary Game Storage
- [ ] Review Licensing
- [ ] Package up a 1.0 relase
- [ ] Create Projects for Each Feature Past the 1.0 release

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
  - [Qt Library](http://download.qt.io/official_releases/online_installers/qt-unified-linux-x64-online.run)

```bash
$ cd Downloads
$ wget http://download.qt.io/official_releases/online_installers/qt-unified-linux-x64-online.run
$ chmod +x qt-unified-linux-x64-online.run
$ sudo ./qt-unified-linux-x64-online.run
```

*Just install all the default options.  For me this install version 5.9.1 in /opt/Qt5.9.1 directory.*

  - Add some additional libraries and compilers if needed
    - Debian/Ubuntu

```bash
$ sudo apt-get -y install build-essential libgl1-mesa-dev libpulse-dev 
```

  - Qt Binding *Note I added environ variables around here
  - .bashrc appended

```bash
  export QT_VERSION=5.9.1 
  export QT_DIR=/opt/Qt5.9.1
```

  - restarted terminal
  - Setup Qt binding for golang library

```bash
$ go get -u -v github.com/therecipe/qt/cmd/...
$ $GOPATH/bin/qtsetup
```

*If everything goes well you 5-10 example apps will test compile and popup towards the end of setup.*

  - Additional Go libraries

```bash
$ go get -u github.com/allan-simon/go-singleinstance
$ go get -u github.com/nicksnyder/go-i18n/...
$ go get -u github.com/sirupsen/logrus
```

- Clone this repository

```bash
$ go get github.com/rashwell/neochess
$ cd $GOPATH/src/github.com/rashwell/neochess
$ $GOPATH/bin/qtdeploy
```

### Windows

For now on windows I use therecipe docker images to do the builds.  The main reason is that for 64bit
building it is a bit difficult currently to use MSYS2 to get to a 64bit mingw version of Qt which will work with Go.  It is doable, just a big pain.  

- Install prerequisites

Download and install [Docker CE for Windows](https://store.docker.com/editions/community/docker-ce-desktop-windows)

This is for a 64bit build that uses the shared dll versions of the Qt library

```powershell
PS > docker pull therecipe/qt:windows_64_shared
```

Install Binding for Qt 


```bash
PS > go get -u -v github.com/therecipe/qt/cmd/...
PS > qtsetup
```

- Additional Go Libraries

```powershell
PS > go get -u github.com/allan-simon/go-singleinstance
PS > go get -u github.com/nicksnyder/go-i18n/...
PS > go get -u github.com/sirupsen/logrus
```

- Clone this repository

```powershell
PS > go get github.com/rashwell/neochess
PS > cd %GOPATH%/src/github.com/rashwell/neochess
```

- Run qtdeploy with the docker switches to build a windows 64bit version of NeoChess

*Somewhere here or when you first try to build Neochess Docker is going to want to Share a folder between the Docker Image and your host machine so that it can read your Go's work GOPATH directory*

```powershell
PS > qtdeploy -docker build windows_64_shared
```

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

