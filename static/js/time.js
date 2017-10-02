var latencies = [];
var avgLatency = 0;
function doLatFetch() {
    var before = Date.now();
    fetch(new Request('ping')).then(function() {
        latencies.push(Date.now() - before);
        if (latencies.length < 10) {
            doLatFetch();
        } else {
            var t = 0;
            for (var i = 0; i < latencies.length; i++) {
                t += latencies[i];
            }
            avgLatency = t / latencies.length;
            document.querySelector('#calcl').style.display = 'none';
            document.querySelector('#display').style.display = 'inline';
            startTimeLoop();
        }
    });
}

function startTimeLoop() {
    setInterval(function() {
        fetch(new Request('ctime_ms')).then(function(resp) {
            return resp.text();
        }).then(function(resp) {
            var difference = Date.now() - parseInt(resp, 10);
            difference -= avgLatency;
            document.querySelector('#lat').innerHTML = Math.abs(difference).toFixed(4) + 'ms';
            document.querySelector('#direction').innerHTML = difference > 0 ? 'ahead of' : 'behind';
        })
    }, 2000);
}
setTimeout(doLatFetch, 250);