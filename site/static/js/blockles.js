const url = new URL(window.location.href);
const params = new URLSearchParams(url.search);
const gameId = params.get("id");

const canvasId = "canvas-" + gameId;
const canvas = document.getElementById(canvasId);
const context = canvas.getContext("2d");

let WIDTH = canvas.scrollWidth;
let HEIGHT = canvas.scrollHeight;
let BLOCK_SIZE = Math.min(Math.floor(WIDTH / 10), Math.floor(HEIGHT / 20));
console.log(BLOCK_SIZE)
context.clearRect(0, 0, WIDTH, HEIGHT);
context.fillStyle = "black";
context.fillRect(0, 0, WIDTH, HEIGHT);
