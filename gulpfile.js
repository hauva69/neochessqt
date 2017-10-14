let gulp = require('gulp');
let doc = require('gulp-task-doc').patchGulp();
let gutil = require('gulp-util');
let download = require('gulp-download');
let unzip = require('gulp-unzip');
let git = require('git-rev');
let exec = require('child_process').exec;
let fs = require('fs');
let os = require('os');

let qtbin = "";
let qtqcollectiongenerator = "";

var deleteFolderRecursive = function (path) {
    if (fs.existsSync(path)) {
        fs.readdirSync(path).forEach(function (file, index) {
            var curPath = path + "/" + file;
            if (fs.lstatSync(curPath).isDirectory()) { // recurse
                deleteFolderRecursive(curPath);
            } else { // delete file
                fs.unlinkSync(curPath);
            }
        });
        fs.rmdirSync(path);
    }
};

/**
 * Clean out application data
 */
gulp.task('clean', function (cb) {
    gutil.log('Cleaning Application Data')
    appdata = "";
    if (os.platform() === 'win32') {
        appdata = process.env.appdata + "/NeoDevelop/NeoChess";
    }
    if (os.platform() === 'linux') {
        appdata = process.env.HOME + "/.local/share/NeoDevelop/NeoChess";
    }
    if (os.platform() === 'darwin') {
        appdata = process.env.HOME + "Library/Preferences/NeoDevelop/NeoChess";
    }
    deleteFolderRecursive(appdata);
});

/**
 * Find Qt Binaries
 */
gulp.task('qtfind', function (cb) {
    gutil.log("Searching for QT Binaries.")
    let qtbinarr = [
        "/opt/Qt/5.9.1/gcc_64/bin/",
        "/opt/Qt5.8.0/5.8/gcc_64/bin/",
        "C:\\msys64\\mingw64\\bin",
        "C:\\Qt\\5.9.1\\mingw53_32\\bin"
    ];
    qtbinfound = false;
    for (i = 0; i < qtbinarr.length; i++) {
        qtbin = qtbinarr[i];
        if (fs.existsSync(qtbin)) {
            qtbinfound = true;
            break;
        }
    }
    if (!qtbinfound) {
        gutil.log('Qt Binary directory not found.');
        gutil.log('Please edit gulpfile.js and adjust.');
        process.exit(1);
    }
    if (os.platform() === 'win32') {
        qtqcollectiongenerator = qtbin + '/qcollectiongenerator.exe';
    }
    if (os.platform() === 'linux') {
        qtqcollectiongenerator = qtbin + '/qcollectiongenerator';
    }
    if (os.platform() === 'darwin') {
        qtqcollectiongenerator = qtbin + '/qcollectiongenerator';
    }

    cb();
})

// @internal
gulp.task('default', ['qtfind', 'help']);

/**
 * Display this help (default)
 */
gulp.task('help', doc.help());

// @internal
gulp.task('copyhelp', function (done) {
    var stream = gulp.src('helpsrc/**/*.{qhc,qch}').pipe(gulp.dest('qml/help/'));
    stream.on('end', function () {
        done();
    });
});

/**
 * Install stockfish engine
 */
gulp.task('getengines', function (cb) {
    gutil.log('Grabbing Stockfis for ' + os.platform());
    var minimatch = require('minimatch');
    var flatten = require('gulp-flatten');
    stockfishurl = "";
    if (os.platform() === 'linux') {
        stockfishurl = "https://stockfish.s3.amazonaws.com/stockfish-8-linux.zip";
        download(stockfishurl).pipe(unzip({
            filter: function (entry) {
                if (minimatch(entry.path, "stockfish_8_x64", {matchBase: true})) {
                    return minimatch(entry.path, "stockfish_8_x64", {matchBase: true})
                }
                if (minimatch(entry.path, "Readme.md", {matchBase: true})) {
                    return minimatch(entry.path, "Readme.md", {matchBase: true})
                }                
                return minimatch(entry.path, "Copying.txt", {matchBase: true})                                
            }
        })).pipe(flatten()).pipe(gulp.dest("linux/"));
    }
});

// @internal
gulp.task('compilehelp', function (done) {
    var stream = exec(qtqcollectiongenerator + ' helpsrc/neochess_US.qhcp -o helpsrc/neochess_US.qhc', function (err, stdout, stderr) {
        gutil.log(stdout);
        gutil.log(stderr);
    });
    stream.on('end', function () {
        done();
    });
});

/**
 * Compile Neochess Translation Data
 */
gulp.task('buildi18n', function (done) {
    var stream = exec('goi18n merge -outdir qml/translate translatesrc/en-US.all.json', function (err, stdout, stderr) {
        gutil.log(stdout);
        gutil.log(stderr);
    });
    stream.on('end', function () {
        done();
    });
});

/**
 * Fast Build NeoChessq Application
 */
gulp.task('buildfast', function (done) {
    git.short(function (rev) {
        gutil.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -fast -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            gutil.log(stdout);
            gutil.log(stderr);
            done(err);
        });
    });
});

/**
 * Build Main Application NeoChess Application
 */
gulp.task('build', function (done) {
    git.short(function (rev) {
        gutil.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            gutil.log(stdout);
            gutil.log(stderr);
            done(err);
        });
    });
});

/**
 * Build / Compile Help files
 */
gulp.task('buildhelp', ['qtfind', 'compilehelp', 'copyhelp']);

/**
 * Compile Help, Compile Translations, Build Main Application
 */
gulp.task('buildall', ['qtfind', 'compilehelp', 'copyhelp', 'buildi18n', 'build']);

