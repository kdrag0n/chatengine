var last60 = [];
var last5 = [];
function avg(arr) {
    var t = 0;
    for (var i = 0; i < arr.length; i++) {
        t += arr[i];
    }
    return t / arr.length;
}
function update() {
    var before = Date.now();
    fetch(new Request('ping')).then(function() {
        var lat = Date.now() - before;
        document.querySelector('#current').innerHTML = lat + 'ms';

        if (last5.length === 5)
            last5.shift();
        last5.push(lat);
        document.querySelector('#min_avg').innerHTML = avg(last5).toFixed(4) + 'ms';

        if (last60.length === 60)
            last60.shift();
        last60.push(lat);
        document.querySelector('#last5_avg').innerHTML = avg(last60).toFixed(4) + 'ms';
    });
}

setTimeout(function() {
    setInterval(update, 1000);
}, 250);