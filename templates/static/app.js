function toHHMMSS(sec_num) {
    sec_num = parseInt(sec_num, 10);
    let hours = Math.floor(sec_num / 3600),
        minutes = Math.floor((sec_num - (hours * 3600)) / 60),
        seconds = sec_num - (hours * 3600) - (minutes * 60);

    if (hours < 10) {
        hours = '0' + hours;
    }
    if (minutes < 10) {
        minutes = '0' + minutes;
    }
    if (seconds < 10) {
        seconds = '0' + seconds;
    }
    return hours + ':' + minutes + ':' + seconds;
}

(function (d) {
    let display = d.querySelector('#countdown .display'),
        timeLeft = parseInt(display.getAttribute('data-left')),
        timer = setInterval(function () {
            if (--timeLeft >= 0) {
                display.innerHTML = toHHMMSS(timeLeft)
            } else {
                clearInterval(timer)
                window.location.reload()
            }
        }, 1000)
})(document)