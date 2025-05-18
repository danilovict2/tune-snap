import { FFmpeg } from "@ffmpeg/ffmpeg";
import { fetchFile } from "@ffmpeg/util";
import Glide from "@glidejs/glide";
import { MediaRecorder, register } from "extendable-media-recorder";
import { connect } from "extendable-media-recorder-wav-encoder";

const ffmpeg = new FFmpeg();
ffmpeg.load();

const listenButton = document.getElementById('listen-button');
listenButton.addEventListener('click', listen);

const recordingDuration = 20000;
const channels = 2;
const sampleRate = '44100';

const glide = document.querySelector('.glide');
const slides = document.querySelector('.glide__slides');
const bullets = document.querySelector('.glide__bullets');

async function listen() {
    listenButton.classList.toggle('pulse');
    listenButton.innerHTML = `<h2>Listening...</h2>`
    listenButton.disabled = true;

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
        .then(r => r.json())
        .then(songs => {
            slides.innerHTML = '';
            bullets.innerHTML = '';
            for (const i in songs) {
                bullets.innerHTML += `
                    <button class="glide__bullet" data-glide-dir="=${i}"></button>
                `

                slides.innerHTML += `
                <li class="glide__slide">
                    <iframe width="560" height="315" src="https://www.youtube.com/embed/${songs[i].SongID}?si=sw6IEZY0IB1BUofG" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
                </li>`
            }

            glide.classList.remove('d-none');

            new Glide('.glide', {
                type: 'slider',
                focusAt: 'center',
            }).mount();

            listenButton.classList.toggle('pulse');
            listenButton.innerHTML = `<h2>Listen</h2>`
            listenButton.disabled = false;
        })
        .catch(e => console.log(e));
}