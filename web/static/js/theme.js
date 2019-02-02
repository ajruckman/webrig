// Cookie functions from: https://www.quirksmode.org/js/cookies.html
function createCookie(name, value, days) {
    let expires;

    if (days) {
        let date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    } else {
        expires = "";
    }

    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name, "", -1);
}

//////////

function updateTheme() {
    const theme = $("#themeSelector").val();
    createCookie("theme", theme, 30);
    location.reload();
}

function loadTheme() {
    let theme = readCookie("theme");
    if (theme == null) {
        createCookie("theme", "dark");
        theme = "dark";
    }

    ['structure', 'style'].forEach(function (sheet) {
        linkStylesheet("/static/css/build/" + sheet + "." + theme + ".css");
    });

    $(function () {
        $("#themeSelector").val(theme);
    })
}

function linkStylesheet(url) {
    const head = document.getElementsByTagName('head')[0];
    const link = document.createElement('link');
    link.rel = 'stylesheet';
    link.type = 'text/css';
    link.href = url;
    head.appendChild(link);
}
