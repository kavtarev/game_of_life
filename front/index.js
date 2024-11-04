const field = document.getElementById('field');
const startBtn = document.getElementById('start-btn');
const stopBtn = document.getElementById('stop-btn');
const nextBtn = document.getElementById('next-btn');
const clearBtn = document.getElementById('clear-btn');

const inputs = []
const COUNT = 10000

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

nextBtn.addEventListener('click', async () => {
  const arr = new Array(COUNT)
  for (let i = 0; i < COUNT; i++) {
    arr[i] = inputs[i].checked ? '1' : '0'
  }
  const r = await fetch('http://localhost:3000/next', {
    method: 'POST', headers: {
      'content-type': 'application/json'
    }, body: JSON.stringify({ data: arr.join('') })
  })

  const js = await r.json()
  updateState(js.data)
})

clearBtn.addEventListener('click', () => {
  inputs.forEach(f => f.checked = false)
})

ws.onmessage = function (msg) {
  try {
    const data = JSON.parse(msg.data)

    switch (data.event) {
      case 'init':
        console.log('init')
        break;
      case 'update':
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

let isMouseDown = false;

document.addEventListener('mousedown', (event) => {
  // Проверяем, что нажата левая кнопка мыши (event.button === 0)
  if (event.button === 0) {
    isMouseDown = true;
  }
});

document.addEventListener('mouseup', () => {
  isMouseDown = false;
});

for (let i = 0; i < COUNT; i++) {
  const ch = document.createElement('input');
  ch.addEventListener('mouseover', () => {
    if (isMouseDown) {
      ch.checked = true
    }
  })

  ch.type = 'checkbox';
  inputs.push(ch)
  field.insertAdjacentElement('beforeend', ch);
}


async function wait() {
  return new Promise(r => {
    setTimeout(() => r(), 2000)
  })
}
