var colors = ['is-danger', 'is-warning', 'is-info', 'is-primary', 'is-success'];
var emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
var bar;
var startScore = 1;
var targetScore = 1;
var setTime;
var animStartTime;

var reqAnimFrame = window.requestAnimationFrame ||
    window.mozRequestAnimationFrame ||
    window.webkitRequestAnimationFrame ||
    window.msRequestAnimationFrame;

function easeInSine(t, b, c, d) {
    return -c * Math.cos(t/d * (Math.PI/2)) + c + b;
}

function easeOutSine(t, b, c, d) {
    return -b * Math.cos(t/d * (Math.PI/2)) + b + c;
}

function pwdChanged() {
    var elem = $('#passwd');
    var value = elem.val();
    var score, label, barClass;

    if ((value || '').length < 1) {
        score = 1;
        label = "";
    } else {
        var zx = zxcvbn(value.substring(0, 101));
        score = Math.min(zx.guesses_log10 / 10, 1.0) * 100;
        label = zx.crack_times_display.online_throttling_100_per_hour;

        barClass = colors[Math.min(Math.floor(score / 20), 4)];
    }

    bar = $('#pwd-strength');
    startScore = parseInt(bar.attr('value'));
    targetScore = score;
    setTime = true;
    reqAnimFrame(barStep);
    bar.attr('class', 'progress ' + barClass);

    var psl = $('#ps-label');
    psl.attr('class', 'label ' + barClass);
    psl.html('Estimated crack time: ' + label);
}

function barStep(time) {
    if (setTime) {
        animStartTime = time;
        setTime = false;
    }

    var f = targetScore < startScore ? easeOutSine : easeInSine;
    var newValue = f(time - animStartTime, startScore, targetScore, 300);
    bar.attr('value', newValue);
    if (newValue < targetScore) {
        reqAnimFrame(barStep);
    }
}

function recCallback() {
    var btn = $('.submit.button');
    if (btn.hasClass('disabled'))
        btn.removeClass('disabled');
}

function recExpCallback() {
    var btn = $('.submit.button');
    if (!btn.hasClass('disabled'))
        btn.addClass('disabled');
}

function emailUpdate() {
    var help = $('.email-field > .help');
    if (!($('.email-field input').val() || '').match(emailRegex)) {
        help.attr('class', 'help is-danger');
        help.html('Please enter a valid e-mail address!');
    } else {
        help.attr('class', 'help is-success');
        help.html('Your e-mail address is valid.');
    }
}

function passwdUpdate() {
    var help = $('.password-field > .help');
    var len = $('.password-field input').val() || '';

    if (len === 0) {
        help.attr('class', 'help is-danger');
        help.html('Please enter a password!');
    } else if (len < 6) {
        help.attr('class', 'help is-danger');
        help.html('Your password must have at least 6 characters!');
    } else if (len > 36) {
        help.attr('class', 'help is-danger');
        help.html('Your password may not be longer than 36 characters!');
    } else {
        help.attr('class', 'help is-success');
        help.html('Your password is valid.');
    }
}

function confirmUpdate() {
    var help = $('.confirm-field > .help');
    var confirm = $('.confirm-field input').val() || '';

    if (confirm.length === 0) {
        help.attr('class', 'help is-danger');
        help.html('Please enter a password!');
    } else if (confirm !== ($('.password-field input').val() || '')) {
        help.attr('class', 'help is-danger');
        help.html("Your passwords don't match!");
    } else {
        help.attr('class', 'help is-success');
        help.html('Your passwords match.');
    }
}