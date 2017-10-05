var gulp = require('gulp');
var doc = require('gulp-task-doc').patchGulp();
var git = require('git-rev');
var exec = require('child_process').exec;
var fs = require('fs');

//var qtbin = "/opt/Qt/5.9.1/gcc_64/bin/";
var qtbin = "/opt/Qt5.8.0/5.8/gcc_64/bin/";

gulp.task('check', function(cb) {
    if (!fs.existsSync(qtbin)) {
        console.log('Qt Binary directory not found at ' + qtbin);
        console.log('Please edit gulpfile.js and adjust');
        process.exit(1);
    }
    cb();
})

// @internal
gulp.task('default', ['check','help']);

/**
 * Display this help
 */
gulp.task('help', doc.help());

/**
 * Compile NeoChess Help Documentation
 */
gulp.task('buildhelp', function (cb) {
    exec(qtbin + 'qcollectiongenerator helpsrc/neochess_US.qhcp -o helpsrc/neochess_US.qhc', function (err, stdout, stderr) {
        console.log(stdout);
        console.log(stderr);
    });
    gulp.src('helpsrc/**/*.{qhc,qch}').pipe(gulp.dest('./qml/help'));
});

/**
 * Compile Neochess Translation Data
 */
gulp.task('buildi18n', function (cb) {
    exec('goi18n merge -outdir qml/translate translatesrc/*.all.json', function (err, stdout, stderr) {
        console.log(stdout);
        console.log(stderr);
        cb(err);
    });
});

/**
 * Fast Build NeoChessq Application
 */
gulp.task('fastbuild', function () {
    git.short(function (rev) {
        console.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -fast -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            console.log(stdout);
            console.log(stderr);
          });    
    });
});

/**
 * Build NeoChess Application
 */
gulp.task('build', function (cb) {
    git.short(function (rev) {
        console.log('Building Neochess Revision: ', rev);
        exec('qtdeploy -ldflags="-X main.REVISION=' + rev + '"', function (err, stdout, stderr) {
            console.log(stdout);
            console.log(stderr);
          });    
    });
});