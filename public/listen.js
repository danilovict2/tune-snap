const listenButton = document.getElementById('listen-button');

listenButton.addEventListener('click', listen)

function listen() {
    listenButton.classList.toggle('pulse');
    listenButton.innerHTML = `<h2>Listening...</h2>`
}