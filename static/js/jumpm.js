(function() {
    // Robert Penner's easeInOutQuad

// find the rest of his easing functions here: http://robertpenner.com/easing/
// find them exported for ES6 consumption here: https://github.com/jaxgeller/ez.js

    var easeInOutQuad = function easeInOutQuad(t, b, c, d) {
        t /= d / 2;
        if (t < 1) return c / 2 * t * t + b;
        t--;
        return -c / 2 * (t * (t - 2) - 1) + b;
    };

    var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) {
        return typeof obj;
    } : function (obj) {
        return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;
    };

    var jumper = function jumper() {
        // private variable cache
        // no variables are created during a jump, preventing memory leaks

        var element; // element to scroll to                   (node)

        var start; // where scroll starts                    (px)
        var stop; // where scroll stops                     (px)

        var offset; // adjustment from the stop position      (px)
        var easing; // easing function                        (function)
        var a11y; // accessibility support flag             (boolean)

        var distance; // distance of scroll                     (px)
        var duration; // scroll duration                        (ms)

        var timeStart; // time scroll started                    (ms)
        var timeElapsed; // time spent scrolling thus far          (ms)

        var next; // next scroll position                   (px)

        var callback; // to call when done scrolling            (function)
        var elem;

        // scroll position helper

        function location() {
            return window.scrollY || window.pageYOffset;
        }

        // element offset helper

        function top(element) {
            return element.getBoundingClientRect().top + start;
        }

        // rAF loop helper

        function loop(timeCurrent) {
            // store time scroll started, if not started already
            if (!timeStart) {
                timeStart = timeCurrent;
            }

            // determine time spent scrolling so far
            timeElapsed = timeCurrent - timeStart;

            // calculate next scroll position
            next = easing(timeElapsed, start, distance, duration);

            // scroll to it
            elem.scrollTop = next;

            // check progress
            timeElapsed < duration ? window.requestAnimationFrame(loop) // continue scroll loop
                : done(); // scrolling is done
        }

        // scroll finished helper

        function done() {
            // account for rAF time rounding inaccuracies
            window.scrollTo(0, start + distance);

            // if scrolling to an element, and accessibility is enabled
            if (element && a11y) {
                // add tabindex indicating programmatic focus
                element.setAttribute('tabindex', '-1');

                // focus the element
                element.focus();
            }

            // if it exists, fire the callback
            if (typeof callback === 'function') {
                callback();
            }

            // reset time for next jump
            timeStart = false;
        }

        // API

        function jump(e, target) {
            var options = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : {};

            // resolve options, or use defaults
            duration = options.duration || 1000;
            offset = options.offset || 0;
            callback = options.callback; // "undefined" is a suitable default, and won't be called
            easing = options.easing || easeInOutQuad;
            a11y = options.a11y || false;

            // cache starting position
            start = location();
            elem = e;

            // resolve target
            switch (typeof target === 'undefined' ? 'undefined' : _typeof(target)) {
                // scroll from current position
                case 'number':
                    element = undefined; // no element to scroll to
                    a11y = false; // make sure accessibility is off
                    stop = start + target;
                    break;

                // scroll to element (node)
                // bounding rect is relative to the viewport
                case 'object':
                    element = target;
                    stop = top(element);
                    break;

                // scroll to element (selector)
                // bounding rect is relative to the viewport
                case 'string':
                    element = document.querySelector(target);
                    stop = top(element);
                    break;
            }

            // resolve scroll distance, accounting for offset
            distance = stop - start + offset;

            // resolve duration
            switch (_typeof(options.duration)) {
                // number in ms
                case 'number':
                    duration = options.duration;
                    break;

                // function passed the distance of the scroll
                case 'function':
                    duration = options.duration(distance);
                    break;
            }

            // start the loop
            window.requestAnimationFrame(loop);
        }

        // expose only the jump method
        return jump;
    };

    window.jump = jumper();
})();