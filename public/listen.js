const listenButton = document.getElementById('listen-button');
listenButton.addEventListener('click', listen);

const recordingDuration = 3000;
const channels = 2;

function listen() {
    listenButton.classList.toggle('pulse');
    listenButton.innerHTML = `<h2>Listening...</h2>`

    navigator.mediaDevices.getUserMedia({ audio: true })
        .then((stream) => {
            const audioContext = new AudioContext();
            const input = audioContext.createMediaStreamSource(stream);
            const recorder = new Recorder(input, { numChannels: channels });

            recorder.record();

            setTimeout(() => {
                recorder.stop();
                recorder.exportWAV(sendAudio, null);
            }, recordingDuration);
        })
        .catch(e => console.log(e));
}

function sendAudio(audio) {
    const formData = new FormData();
    formData.append('sample', audio);

    fetch('/api/recognize', {
        method: 'POST',
        body: formData,
    })
    .then(r => console.log(r))
    .catch(e => console.log(e));
}