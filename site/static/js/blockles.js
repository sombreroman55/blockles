function getCookieByName(name) {
  if (!cookieTable[name]) {
    let cookies = document.cookie.split(";");
    for (var i = 0; i < cookies.length; i++) {
      let c = cookies[i].trim();
      let p = c.split("=");
      cookieTable[p[0].trim()] = p[1].trim();
    }
  }
  return cookieTable[name];
}

var cookieTable = {};

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

// TODO: Echo user text into the tetris buffer for testing, only in that area and make it scrollable
const host = window.location.host;
const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
const wshost = protocol + host + "/solows?id=" + gameId;
console.log(wshost);
let ws = new WebSocket(wshost);
ws.addEventListener("open", (event) => {
  let testCommand = {
    sender: getCookieByName("playerId"),
    command: "Hello World!"
  };
  let s = JSON.stringify(testCommand);
  console.log(testCommand);
  console.log(s);
  ws.send(s);
});
