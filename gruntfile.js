module.exports = function (grunt) {

    // Project configuration.
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        sass: {
            options: {
                sourceMap: false
            },
            dist: {
                files: [{
                    expand: true,
                    cwd: 'assets/sass',
                    src: '*.scss',
                    dest: 'assets/css/sass',
                    ext: '.generated.css'
                }]
            }
        },
        concat: {
            js: {
                src: [
                    'assets/js/third-party/*.js',
                    'assets/js/*.js'
                ],
                dest: 'assets/tmp/concatenate.js'
            },
            css: {
                src: [
                    'assets/css/third-party/*.css',
                    'assets/css/sass/*.css'
                ],
                dest: 'assets/tmp/concatenate.css'
            }
        },
        cssmin: {
            target: {
                files: {
                    'assets/compiled.min.css': ['assets/tmp/concatenate.css']
                }
            }
        },
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */\n'
            },
            build: {
                src: 'assets/tmp/*.js',
                dest: 'assets/compiled.min.js'
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
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-sass');

    // Default tasks.
    grunt.registerTask('default', [
        // CSS
        'sass',
        'concat:css',
        'cssmin',

        // JS
        'concat:js',
        'uglify'

        // Watch
        // 'watch'
    ]);
};
