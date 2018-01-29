
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
