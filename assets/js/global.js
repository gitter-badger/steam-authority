$("[data-link]").click(function () {
    var link = $(this).attr('data-link');
    if (link) {
        window.location.href = $(this).attr('data-link');
    }
});

function clearField(evt, input) {
    var code = evt.charCode || evt.keyCode;
    if (code === 27) {
        input.value = '';
    }
}
