if (typeof scrollTo === 'string') {
    window.scroll({
        top: $(scrollTo).offset().top - 100,
        left: 0,
        behavior: 'smooth'
    });
}
