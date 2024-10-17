const field = document.getElementById('field');
const inputs = []
const COUNT = 100

alert(1)
const ws = new WebSocket('ws://127.0.0.1:3000/ws')
alert(2)

ws.onmessage = function (msg) {
  console.log(msg);
}
alert(13)

console.log(ws);

for (let i = 0; i < COUNT; i++) {
  const ch = document.createElement('input');
  ch.type = 'checkbox';
  inputs.push(ch)
  field.insertAdjacentElement('beforeend', ch);
}


async function wait() {
  return new Promise(r => {
    setTimeout(() => r(), 2000)
  })
}

async function run() {
  await wait()
  const str = new Array(COUNT).fill('1').join('')

  console.time('start')
  for (let i = 0; i < str.length; i++) {
    if (str[i] === '1') {
      inputs[i].checked = true
    }
  }
  console.timeEnd('start')

}

run()