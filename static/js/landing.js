$(function() {
    var rellax = new Rellax('.parallax', {
        center: true
    });
    console.log('me load');

    setTimeout(initConsent, 2000);
});

function initConsent() {
    if (window['cookieconsent'] !== undefined) {
        window.cookieconsent.initialise({
            "palette": {
                "popup": {
                    "background": "#1d8a8a"
                },
                "button": {
                    "background": "#62ffaa"
                }
            },
            "theme": "classic",
            "position": "bottom-right",
            "content": {
                "message": "We use cookies to improve your experience on our website."
            }
        })
    } else {
        setTimeout(initConsent, 250);
    }
}