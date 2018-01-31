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
                dest: 'assets/concatenate.js'
            },
            css: {
                src: [
                    'assets/css/*.css',
                    'assets/css/sass/*.css'
                ],
                dest: 'assets/concatenate.css'
            }
        },
        cssmin: {
            options: {
                sourceMap: true,
                roundingPrecision: -1
            },
            target: {
                files: {
                    'assets/compiled.min.css': ['assets/concatenate.css']
                }
            }
        },
        uglify: {
            options: {
                banner: '/*! <%= pkg.name %> <%= grunt.template.today("yyyy-mm-dd") %> */',
                compress: true,
                sourceMap: true,
                'mangle.properties': true
            },
            build: {
                src: 'assets/concatenate.js',
                dest: 'assets/compiled.min.js'
            }
        },
        watch: {
            sass: {
                files: ['assets/sass/*.scss'],
                tasks: ['sass', 'concat:css', 'cssmin', 'clean']
            },
            js: {
                files: ['assets/js/*.js'],
                tasks: ['concat:js', 'uglify', 'clean']
            }
        },
        clean: [
            'assets/css/sass/',
            'assets/concatenate.js',
            'assets/concatenate.css',
        ]
    });

    // Load the plugin that provides the tasks
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-watch');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-sass');

    // Default tasks.
    grunt.registerTask('default', [
        // CSS
        'sass',
        'concat:css',
        'cssmin',

        // JS
        'concat:js',
        'uglify',

        //
        'clean',
        'watch'
    ]);
};
