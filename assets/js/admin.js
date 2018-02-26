if ($('#admin-page').length > 0) {

    var $actions = $('#actions a');

    $actions.on('click', function () {
        var text = $(this).find('p').text();
        return confirm(text + '?');
    });

    $actions.hover(
        function () {
            $(this).addClass('list-group-item-danger')
        },
        function () {
            $(this).removeClass('list-group-item-danger')
        }
    );
}
