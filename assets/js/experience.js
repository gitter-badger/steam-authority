(function ($, window, document, undefined) {
    'use strict';

    function scroll() {

        if (typeof scrollTo === 'string') {
            window.scroll({
                top: $(scrollTo).offset().top - 100,
                left: 0,
                behavior: 'smooth'
            });

            $('tr').removeClass('table-success');
            $(scrollTo).addClass('table-success');
        }
    }

    scroll();

    $('#xp-page table tr td').click(function () {

        var level = $(this).parent().attr('data-level');

        if (history.pushState) {
            history.pushState('data', '', '/experience/' + level);
        }

        scrollTo = 'tr[data-level=' + level + ']';
        scroll();
    });

})(jQuery, window, document);
