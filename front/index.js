const field = document.getElementById('field');
const startBtn = document.getElementById('start-btn');
const stopBtn = document.getElementById('stop-btn');

const inputs = []
const COUNT = 100

const ws = new WebSocket('ws://127.0.0.1:3000/ws')

startBtn.addEventListener('click', () => {
  const arr = new Array(COUNT)
  for (let i = 0; i < COUNT; i++) {
    arr[i] = inputs[i].checked ? '1' : '0'
  }
  ws.send(JSON.stringify({ event: 'start', data: arr.join('') }))
})

stopBtn.addEventListener('click', () => {
  ws.send(JSON.stringify({ event: 'stop' }))
})

ws.onmessage = function (msg) {
  try {
    const data = JSON.parse(msg.data)
    switch (data.event) {
      case 'init':
        console.log('init')
        break;
      case 'update':
        console.log(data);
        updateState(data.data)
        break;
      default:
        console.log(99999);
        break;
    }
  } catch (error) {
    console.log("error: ", error)
  }
}

function updateState(str) {
  for (let i = 0; i < str.length; i++) {
    inputs[i].checked = str[i] === '1' ? true : false
  }
}

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

// async function run() {
//   await wait()
//   const str = new Array(COUNT).fill('1').join('')

//   console.time('start')
//   for (let i = 0; i < str.length; i++) {
//     if (str[i] === '1') {
//       inputs[i].checked = true
//     }
//   }
//   console.timeEnd('start')

// }

// run()