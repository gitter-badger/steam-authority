if ($('#settings-page').length > 0) {

    var $checkbox = $('#browser-alerts');

    $checkbox.on('click', function () {
        if ($(this).is(':checked')) {

            Push.Permission.request(
                function () {
                },
                function () {
                    alert('You have denied notification access in your browser.');
                    $(this).prop("checked", false);
                }
            );
        }
    });

    $('[data-toggle="tooltip"]').tooltip(options)
}
