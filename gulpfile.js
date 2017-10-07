let gulp = require('gulp');
let doc = require('gulp-task-doc').patchGulp();
var gutil = require('gulp-util');
let git = require('git-rev');
let exec = require('child_process').exec;
let fs = require('fs');
let os = require('os');

let qtbin = "";

var deleteFolderRecursive = function(path) {
    if (fs.existsSync(path)) {
      fs.readdirSync(path).forEach(function(file, index){
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
gulp.task('clean', function(cb) {
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
gulp.task('qtfind', function(cb) {
    gutil.log("Searching for QT Binaries.")
    let qtbinarr = [
        "/opt/Qt/5.9.1/gcc_64/bin/",
        "/opt/Qt5.8.0/5.8/gcc_64/bin/",
        "C:\\msys64\\mingw64\\bin"
    ];
    qtbinfound = false;
    for (i=0;i<qtbinarr.length;i++) {
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
    cb();
})

// @internal
gulp.task('default', ['qtfind','help']);

/**
 * Display this help (default)
 */
gulp.task('help', doc.help());

/**
 * Compile NeoChess Help Documentation
 */
gulp.task('buildhelp', function (cb) {
    exec(qtbin + 'qcollectiongenerator helpsrc/neochess_US.qhcp -o helpsrc/neochess_US.qhc', function (err, stdout, stderr) {
        gutil.log(stdout);
        gutil.log(stderr);
    });
    gulp.src('helpsrc/**/*.{qhc,qch}').pipe(gulp.dest('./qml/help'));
});

/**
 * Compile Neochess Translation Data
 */
gulp.task('buildi18n', function (cb) {
    exec('goi18n merge -outdir qml/translate translatesrc/en-US.all.json', function (err, stdout, stderr) {
        gutil.log(stdout);
        gutil.log(stderr);
        cb(err);
    });
});

/**
 * Fast Build NeoChessq Application
 */
gulp.task('buildfast', function () {
    git.short(function (rev) {
        gutil.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -fast -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            gutil.log(stdout);
            gutil.log(stderr);
          });    
    });
});

/**
 * Build Main Application NeoChess Application
 */
gulp.task('build', function (cb) {
    git.short(function (rev) {
        gutil.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            gutil.log(stdout);
            gutil.log(stderr);
          });    
    });
});

/**
 * Compile Help, Compile Translations, Build Main Application
 */
gulp.task('buildall', ['qtfind','buildhelp','buildi18n','build']);

