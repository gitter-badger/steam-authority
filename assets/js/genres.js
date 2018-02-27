if ($('#genres-page').length > 0) {

    var options = {
        valueNames: ['genre-name', 'genre-games'],
        listClass: 'genres-list',
        page: 1000,
        fuzzySearch: {
            searchClass: 'genres-search',
            location: 0,
            threshold: 0.5
        }
    };

    new List('genres-page', options);
}
