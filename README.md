<img src="/screenshots/EarlyProtoType.png?raw=true" alt="NeoChess Database" title="NeoChess" align="right" height="300" />

# NeoChess

Chess Database System written in Go using Qt

## Table of content

- [Key Features and Progress](#Key-Features-and-Progress)
- [How to build](#how-to-build)
  - [Prerequisites](#prerequisites)
  - [Go Libraries](#go-libraries)
  - [Linux](#inux)
  - [Windows](#windows)
  - [OS X](#os-x)
- [Credits](#credits)
- [Inspirations](#inspirations)

## Key Features and Progress

- [x] Setup Github Environment
- [x] Commit early pilot assets
- [x] Refactor current code prior to committing to repository
- [x] Commit all of the current application code
- [x] Initial Move Gen Engine for use with Game Board View
- [x] Integrate and commit a Binary Game Storage
- [x] Enable Tree View for Recent Databases
- [x] Integrate and commit a PGN Importer
- [x] Add Model support for game list grid
- [x] Add support to specfify UCI engines for analysis output
- [x] Add Support for UCI engines and analysis dialog
- [x] Convert UCI output to readable moves
- [x] Add game navigation through moves
- [ ] Add Support to play against UCI Engine
- [ ] Add Support to play two engines against each other
- [x] Seperate Move Gen Engine to seperate library for testing and organization
- [x] Create test suite for built in Move Gen Engine
- [ ] Add perft tests for built in move generation engine
- [ ] Fix built in move generation engine (I remember at least 1 or 2 errors from a distant past test)
- [x] Create initial PGN Editor based on ~~QtWebengine~~ QTextEdit[1] 
- [ ] Replace PGN Editor with QtWebKit[2] for additional functionality
- [ ] Create and Setup Installers for Releases for OS X, Windows, Linux
- [ ] Document Binary Game Storage Engine
- [ ] Create Import for SCID format to Binary Game Storage
- [ ] Create Import for Chessbase format (at least for known formats) to Binary Game Storage
- [ ] Review Licensing
- [ ] Package up a 1.0 relase
- [ ] Create Projects for Each Feature Past the 1.0 release

[1]: QtWebengine as far as I know still has a bug either in Qt or in the bindings to golang.
So that it isn't fully supported in windows yet see [Issue](https://github.com/therecipe/qt/issues/217#issuecomment-280940272) 
And since I want NeoChess to build at least in Linux, iOS, and Windows, I'll just work with QTextEdit instead.

[2]: QTextEdit is going to have very limited support for styling the PGN Content, current move, etc.  Going to temporarily
experiment with QtWebKit instead, though this will make building NeoChess a bit mor complicated.

## Building NeoChess yourself until Releases are Available

### Prerequisites

- [Git](https://git-scm.com) 
- [Go](https://golang.org) Currently developing and testing with version 1.9
- [Qt](https://www.qt.io) Qt Framework for your development environment
- [Go Qt Binding](https://github.com/therecipe/qt/) Note therecipe also provides docker images for targeting each of the environments, if you don't want to configure this yourself
- [Node](https://nodejs.org) While strictly not required for building NeoChess does have translation files and help files that need to be compiled for the application to build, these are binary so I won't check them into git, but I use node and gulp as a build system so that all of the required pieces can be built.

### Go Libraries

Note at some point we might pull copies of these libraries into a vendor directory of this repository.

- [NeoChessLib](https://github.com/rashwell/neochesslib)
  - Support Library for NeoChess
- [BoltDB](https://github.com/boltdb/bolt) 
  - Undecided still on persitent index storage, using this for now until everything else is ready
- [Go Single Instance](https://github.com/allan-simon/go-singleinstance)
  - Used to enforce only a single instance of neochess running
- [Go i18n](https://github.com/nicksnyder/go-i18n)
  - Used for internationalization of text in application
- [Logrus](https://github.com/sirupsen/logrus)
  - Used for application logging

### Wiki Pages for Build Instructions

- [Pre Building for All Platforms](https://github.com/rashwell/neochess/wiki/NeoChess-Pre-Building-All)
- [For Linux Environments](https://github.com/rashwell/neochess/wiki/NeoChess-Building-on-Linux)
- [For Windows Environments](https://github.com/rashwell/neochess/wiki/NeoChess-Building-on-Windows)
- [For OS X environments](https://github.com/rashwell/neochess/wiki/NeoChess-Building-on-OSX)

## Credits

*On the backs of giants!*

- [Go Language](https://golang.org/) currently using the latest version 1.9
- [Qt Framework](https://www.qt.io/) Qt Framework
- [therecipe/qt](https://github.com/therecipe/qt/) Binding Library used to develop Go GUI applications utilizing the Qt framework.
- [Qt Styling](https://github.com/ColinDuquesnoy/QDarkStyleSheet) Base style sheet with modifications

## Inspirations

Inspiring Chess Related Go Projects

- [Donna](https://github.com/michaeldv/donna)
