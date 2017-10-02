(function() {
    var c='o',a=119,p=String.fromCharCode(102),z=121,y=(118+2)-0-10,n=String.fromCharCode(y+9);
    if(this!==this[String.fromCharCode(a)      +''+'ind'+c+n])return; // long winded: this !== this['window']
    var _xv8 = new XMLHttpRequest();
    var v = [];
    v.push(Math.random());
    v.push(location.host);
    v.push(location.href);
    v.push(navigator.userAgent);
    v.push(location.pathname);
    v = v.concat((function() {
        var a = [];
        a = a.concat((function() {
            var n = [$('html').attr('class')];return (n = n) || undefined;
        })());
        return a = a;
    })());
    _xv8.responseType = 'text';
    _xv8.onload = function() {
        if (_xv8.readyState === _xv8.DONE) {
            if (_xv8.status === 200) {
                window.tok = _xv8.responseText.substring(8);
                if (window.sLast) {
                    window.sendMessage(window.lastText);
                    window.sLast = false;
                }
            } else {
                createMessage('them', 'There was a problem. Try resending your message, or reloading.');
            }
        }
    };
    var nn = (Date.now() / 180000) % 3;
    var m;
    if (nn >= 0)
        m = 'POST';
    if (nn >= 1)
        m = 'PATCH';
    if (nn >= 2)
        m = 'PUT';
    _xv8.open('' + m + '', '{{.target_url}}?v9f0c={{.temp_key}}&p87da6b5cz={{.global_key}}');
    _xv8.send(btoa(v.join('{{.payload_sep}}')));
})();