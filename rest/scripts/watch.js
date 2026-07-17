document.addEventListener("DOMContentLoaded", () => {

    let video = document.getElementsByTagName('video')[0];

    if (video) {
        video.currentTime = {currentTime};

        let lastProgressUpdate = new Date();
        video.addEventListener('timeupdate', (e) => {
            let now = new Date();
            let elapsed = now - lastProgressUpdate;
            if (elapsed > 5000) {
                fetch('/progress', {
                    method: 'post',
                    headers: {
                        'Content-Type': 'application/json'},
                    body: JSON.stringify({
                        v: '{videoId}',
                        t: video.currentTime.toString()})
                }).then((resp) => { if (resp && !resp.ok) {
                    console.log(resp)}
                });
                lastProgressUpdate = now;
            }});

        video.addEventListener('ended', (e) => {
            fetch('/end/{videoId}/completed', {method: 'get'}).then((resp) => { if (resp && !resp.ok) {
                console.log(resp)}});
                if (prg) {prg.value = prg.max}
        });
    }

    document.body.addEventListener('keydown', (e) => {
                switch (e.keyCode) {
            // ArrowRight
                    case 39:
                    e.preventDefault();
                    video.currentTime += 15;
                    break;
            // ArrowLeft
                    case 37:
                    e.preventDefault();
                    video.currentTime -= 15;
                    break;
            // Space
                    case 32:
                    e.preventDefault();
                    video.paused ? video.play() : video.pause();
                    break;
                };
    });

})

