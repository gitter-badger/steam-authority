if ($('#queue-page').length > 0) {

    var time = 5;

    var interval = setInterval(function () {

        time--;

        $('#live-badge').html('Live (' + time + ')');

        if (time === 0) {
            clearInterval(interval);
            location.reload();
        }

    }, 1000);

}
