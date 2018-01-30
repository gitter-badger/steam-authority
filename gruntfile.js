module.exports = function (grunt) {

    // Project configuration.
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n'
            },
            build: {
                src: 'assets/tmp/*.js',
                dest: 'assets/_compiled.min.js'
            }
        },
        less: {
            options: {
                paths: ['assets/css/*.less'],
                modifyVars: {
                    imgPath: '"http://mycdn.com/path/to/images"',
                    bgColor: 'red'
                }
            },
            files: {
                'path/to/result.css': 'path/to/source.less'
            }
        },
        sass: {
            dist: {
                options: {
                    sourcemap: 'none'
                },
                files: {
                    'assets/tmp/concatenate.scss.css': 'assets/tmp/concatenate.scss',
                }
            }
        },
        concat: {
            js: {
                src: [
                    'assets/js/third-party/*.js',
                    'assets/js/*.js',
                ],
                dest: 'assets/tmp/concatenate.js',
            },
            sass: {
                src: [
                    'assets/sass/*.scss',
                ],
                dest: 'assets/tmp/concatenate.scss'
            },
            css: {
                src: [
                    'assets/css/third-party/*.css',
                    'assets/css/*.css',
                    'assets/tmp/concatenate.scss.css',
                ],
                dest: 'assets/tmp/all-css.css'
            }
        },
        cssmin: {
            target: {
                files: {
                    'assets/_compiled.min.css': ['assets/tmp/all-css.css']
                }
            }
        },
        watch: {
            css: {
                files: ['assets/css/*.css', 'assets/css/third-party/*.css', 'assets/sass/*.scss'],
                tasks: ['concat:sass', 'sass', 'concat:css', 'cssmin']
            },
            js: {
                files: ['assets/js/*.js', 'assets/js/third-party/*.js'],
                tasks: ['concat:js', 'uglify']
            }
        }
    });

    // Load the plugin that provides the tasks
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-sass');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-cssmin');

    //
    // grunt.task.run('notify_hooks');

    // Default tasks.
    grunt.registerTask('default', [
        // CSS
        'concat:sass',
        'sass',
        'concat:css',
        'cssmin',

        // JS
        'concat:js',
        'uglify',

        // Watch
        'watch'
    ]);
};
