var gulp = require('gulp');
var doc  = require('gulp-task-doc').patchGulp();

// @internal
gulp.task('default', ['help']);

/**
 * Display this help
 */
gulp.task('help', doc.help());

/**
 * Compile NeoChess Help Documentation
 */
gulp.task('buildhelp', function() {
});

/**
 * Compile Neochess Translation Data
 */
gulp.task('buildi18n', function() {
});

/**
 * Fast Build NeoChess Application
 */
gulp.task('fastbuild', function() {
});

/**
 * Build NeoChess Application
 */
gulp.task('build', function() {

});