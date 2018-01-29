// Tabs
$('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
    var hash = $(e.target).attr('href');
    if (history.pushState) {
        history.pushState(null, null, hash);
    } else {
        location.hash = hash;
    }
});

var hash = window.location.hash;
if (hash) {
    $('.nav-link[href="' + hash + '"]').tab('show');
}


// XP
if (typeof scrollTo === 'string') {
    window.scroll({
        top: $(scrollTo).offset().top - 100,
        left: 0,
        behavior: 'smooth'
    });
}

// Ranks
$("[data-link]").click(function () {
    window.location.href = $(this).attr('data-link');
});

// Apps
$('select.form-control-chosen').chosen({
    disable_search_threshold: 10,
    allow_single_deselect: true,
    rtl: false
});

// Chat
if ($('#chat-page').length > 0) {
    if (window.WebSocket === undefined) {
        console.log('Your browser does not support WebSockets');
    } else {
        var socket = new WebSocket("ws://" + location.host + "/websocket");
        var $badge = $('#live-badge');

        socket.onopen = function (e) {
            // console.log('WebSocket opened');
            $badge.addClass('badge-success').removeClass('badge-secondary badge-danger')
        };
        socket.onclose = function (e) {
            // console.log('WebSocket closed');
            $badge.addClass('badge-danger').removeClass('badge-secondary badge-success')
        };
        socket.onmessage = function (e) {
            // console.log('WebSocket recieved');
            // console.log(e.data);

            var data = jQuery.parseJSON(e.data);

            if (data.Page === 'chat') {

                data = data.Data;

                $('ul.list-unstyled').prepend($(
                    '<li class="media">' +
                    '    <img class="mr-3" src="https://cdn.discordapp.com/avatars/' + data.author.id + '/' + data.author.avatar + '.png?size=128" alt="' + data.author.username + '">' +
                    '    <div class="media-body">' +
                    '        <h5 class="mt-0 mb-1">' + data.content + '</h5>' +
                    '        <p class="text-muted">By ' + data.author.username + '</p>' +
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
