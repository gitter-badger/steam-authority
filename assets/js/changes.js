if ($('#changes-page').length > 0) {

    if (window.WebSocket === undefined) {
        console.log('Your browser does not support WebSockets');
    } else {
        var socket = new WebSocket("ws://" + location.host + "/websocket");
        var $badge = $('#live-badge');

        socket.onopen = function (e) {
            console.log('WebSocket opened');
            $badge.addClass('badge-success').removeClass('badge-secondary badge-danger')
        };
        socket.onclose = function (e) {
            console.log('WebSocket closed');
            $badge.addClass('badge-danger').removeClass('badge-secondary badge-success')
        };
        socket.onmessage = function (e) {
            console.log('WebSocket recieved');
            console.log(e.data);

            var data = jQuery.parseJSON(e.data);

            if (data.Page === 'changes') {

                data = data.Data;

                $('ul.list-unstyled').prepend($(
                    '<li class="media">' +
                    '    <div class="media-body">' +
                    '        <h5 class="mt-0 mb-1">Change ' + data.id + '</h5>' +
                    '        <p class="text-muted" data-livestamp="' + data.created_at + '" style="margin-bottom: 0;">' + data.created_at + '</p>' +
                    '    </div>' +
                    '</li>'));
            }
        };
        socket.onerror = function (e) {
            console.log('WebSocket error');
            $badge.addClass('badge-danger').removeClass('badge-secondary badge-success')
        };
    }
}
