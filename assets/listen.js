import { FFmpeg } from "@ffmpeg/ffmpeg";
import { fetchFile } from "@ffmpeg/util";
import { MediaRecorder, register } from "extendable-media-recorder";
import { connect } from "extendable-media-recorder-wav-encoder";

const ffmpeg = new FFmpeg();
ffmpeg.load();

const listenButton = document.getElementById('listen-button');
listenButton.addEventListener('click', listen);

const recordingDuration = 20000;
const channels = 2;
const sampleRate = '44100';

async function listen() {
    listenButton.classList.toggle('pulse');
    listenButton.innerHTML = `<h2>Listening...</h2>`

    if (!ffmpeg.loaded) {
        await ffmpeg.load();
    }

    await register(await connect());

    const stream = await navigator.mediaDevices.getUserMedia({ audio: { channelCount: channels, sampleSize: 16 } });
    const tracks = stream.getAudioTracks();
    const audioStream = new MediaStream(tracks);

    tracks[0].onended = () => {
        audioStream.getTracks().forEach(t => t.stop());
    };

    for (const track of stream.getVideoTracks()) {
        track.stop();
    }

    const recorder = new MediaRecorder(stream, {
        mimeType: 'audio/wav',
    });

    recorder.start();
    const chunks = [];
    recorder.ondataavailable = e => {
        chunks.push(e.data);
    };

    setTimeout(() => {
        recorder.stop();
    }, recordingDuration);

    recorder.addEventListener('stop', async () => {
        const blob = new Blob(chunks, { type: 'audio/wav' });
        const inputFile = 'input.wav';
        const outputFile = 'output.wav';

        await ffmpeg.writeFile(inputFile, await fetchFile(blob));
        const exitCode = await ffmpeg.exec(['-i', inputFile, '-c', 'pcm_s16le', '-ar', sampleRate, "-ac", channels, outputFile]);
        if (exitCode !== 0) {
            console.log(`ffmpeg exited with code: ${exitCode}`);
            return;
        }

        const data = await ffmpeg.readFile(outputFile);
        const outputBlob = new Blob([data.buffer], { type: 'audio/wav' });
        sendAudio(outputBlob);
    });
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