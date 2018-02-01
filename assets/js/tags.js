if ($('#tags-page').length > 0) {

    var options = {
        valueNames: ['tag-name'],
        listClass: 'tags-list',
        page: 1000,
        pagination: false,
        fuzzySearch: {
            searchClass: 'tags-search',
            location: 0,
            // distance: 100,
            threshold: 0.5,
            // multiSearch: true
        }
    };

    new List('tags-page', options);
}
