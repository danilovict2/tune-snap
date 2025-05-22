import { FFmpeg } from "@ffmpeg/ffmpeg";
import { fetchFile } from "@ffmpeg/util";
import Glide from "@glidejs/glide";
import { MediaRecorder, register } from "extendable-media-recorder";
import { connect } from "extendable-media-recorder-wav-encoder";
import { createElement, Mic, MicOff, Monitor, MonitorOff } from "lucide";

await register(await connect());

const ffmpeg = new FFmpeg();
await ffmpeg.load();

let audioInput = 'mic';

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

    try {
        const mediaDevice =
            audioInput === "device"
                ? navigator.mediaDevices.getDisplayMedia.bind(navigator.mediaDevices)
                : navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices);

        const stream = await mediaDevice({ audio: { channelCount: channels, sampleSize: 16 } });
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
    } catch (e) {
        console.log(e);
        reset();
    }
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
            const MAX_RESULTS = 5;
            for (const i in songs) {
                if (i == MAX_RESULTS) {
                    break
                }

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
        })
        .catch(e => console.log(e))
        .finally(() => reset());
}

function reset() {
    listenButton.classList.remove('pulse');
    listenButton.innerHTML = `<h2>Listen</h2>`
    listenButton.disabled = false;
}

const micButton = document.getElementById('mic');
const monitorButton = document.getElementById('monitor');

const mic = createElement(Mic);
const monitor = createElement(Monitor);
const micOff = createElement(MicOff);
const monitorOff = createElement(MonitorOff);

function changeAudioInput() {
    micButton.firstElementChild.remove();
    monitorButton.firstElementChild.remove();

    if (audioInput === 'mic') {
        audioInput = 'device';
        micButton.appendChild(micOff);
        monitorButton.appendChild(monitor);
    } else {
        audioInput = 'mic';
        micButton.appendChild(mic);
        monitorButton.appendChild(monitorOff);
    }
}

micButton.addEventListener('click', changeAudioInput);
monitorButton.addEventListener('click', changeAudioInput);

changeAudioInput();