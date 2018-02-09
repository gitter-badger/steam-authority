$("[data-link]").click(function () {
    var link = $(this).attr('data-link');
    if (link) {
        window.location.href = $(this).attr('data-link');
    }
});
