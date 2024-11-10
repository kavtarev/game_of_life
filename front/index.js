const field = document.getElementById('field');
const startBtn = document.getElementById('start-btn');
const stopBtn = document.getElementById('stop-btn');
const nextBtn = document.getElementById('next-btn');
const clearBtn = document.getElementById('clear-btn');
const computeBtn = document.getElementById('compute-btn');
const handleStringBtn = document.getElementById('handle-string-btn');
const handleBytesBtn = document.getElementById('handle-bytes-btn');

const inputs = []
const COUNT = 10000
const host = '127.0.0.1:3000'

const ws = new WebSocket(`ws://${host}/ws`)

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
  const r = await fetch(`http://${host}/next`, {
    method: 'POST', headers: {
      'content-type': 'application/json'
    }, body: JSON.stringify({ data: defaultGetState() })
  })

  const js = await r.json()
  updateState(js.data)
})

handleStringBtn.addEventListener('click', async () => {
  const r = await fetch(`http://${host}/handle-string`, {
    method: 'POST', headers: {
      'content-type': 'application/json'
    }, body: JSON.stringify({ data: defaultGetState() })
  })
})
handleBytesBtn.addEventListener('click', async () => {
  const r = await fetch(`http://${host}/handle-byte`, {
    method: 'POST', headers: {
    }, body: getStateAsArrayBuffer()
  })
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
        throw new Error('should not be here ever');
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
  ch.checked = Math.random() > 0.4
  inputs.push(ch)
  field.insertAdjacentElement('beforeend', ch);
}

async function wait() {
  return new Promise(r => {
    setTimeout(() => r(), 2000)
  })
}

function defaultGetState() {
  const arr = new Array(COUNT)
  for (let i = 0; i < COUNT; i++) {
    arr[i] = inputs[i].checked ? '1' : '0'
  }

  return arr.join('')
}

function getStateAsArrayBuffer() {
  const buf = new Array(COUNT / 8)
  let index = 0

  for (let i = 0; i < COUNT; i += 8) {
    let temp = 0;
    for (let j = 0; j < 8 && j + i < COUNT; j++) {
      temp |= inputs[i + j].checked << 7 - j
    }
    buf[index] = temp
    index++
  }

  return new Uint8Array(buf);
}

function getStateAsCountOccurrence() {
  const buf = []
  let amount = 1

  for (let i = 1; i < COUNT; i++) {
    if (inputs[i].checked === inputs[i - 1].checked) {
      amount++
      continue;
    }
    buf.push(amount, Number(inputs[i - 1].checked))
    amount = 1
  }

  return buf
}

function getStateAsPhil() {
  const obj = []
  let amountChecked = 0
  let amountUnchecked = 0
  const objReverse = []

  for (let i = 0; i < COUNT; i++) {
    if (inputs[i].checked) {
      obj.push(i)
      amountChecked++
    } else {
      objReverse.push(i)
      amountUnchecked++
    }
  }

  if (amountChecked > amountUnchecked) {
    return { isRevers: true, obj: obj }
  }

  return { isRevers: false, obj: objReverse }

}

function getStateAsPhilString() {
  const obj = ""
  let amountChecked = 0
  let amountUnchecked = 0
  const objReverse = ""

  for (let i = 0; i < COUNT; i++) {
    if (inputs[i].checked) {
      obj.concat(i)
      amountChecked++
    } else {
      objReverse.concat(i)
      amountUnchecked++
    }
  }

  if (amountChecked > amountUnchecked) {
    return { isRevers: true, obj: obj }
  }

  return { isRevers: false, obj: objReverse }

}

computeBtn.addEventListener('click', () => {
  const amount = 1000;

  console.time("compute as string")
  for (let i = 0; i < amount; i++) {
    defaultGetState()
  }
  console.timeEnd("compute as string")

  console.time("compute as array buffer")
  for (let i = 0; i < amount; i++) {
    getStateAsArrayBuffer()
  }
  console.timeEnd("compute as array buffer")

  console.time("compute as amount in order")
  for (let i = 0; i < amount; i++) {
    getStateAsCountOccurrence()
  }
  console.timeEnd("compute as amount in order")

  console.time("compute as Phil suggested as array")
  for (let i = 0; i < amount; i++) {
    getStateAsPhil()
  }
  console.timeEnd("compute as Phil suggested as array")

  console.time("compute as Phil suggested as string")
  for (let i = 0; i < amount; i++) {
    getStateAsPhilString()
  }
  console.timeEnd("compute as Phil suggested as string")

  console.time("PART2 compute as string")
  for (let i = 0; i < amount; i++) {
    defaultGetState()
  }
  console.timeEnd("PART2 compute as string")

  console.time("PART2 compute as array buffer")
  for (let i = 0; i < amount; i++) {
    getStateAsArrayBuffer()
  }
  console.timeEnd("PART2 compute as array buffer")

  console.time("PART2 compute as amount in order")
  for (let i = 0; i < amount; i++) {
    getStateAsCountOccurrence()
  }
  console.timeEnd("PART2 compute as amount in order")

  console.time("PART2 compute as Phil suggested as array")
  for (let i = 0; i < amount; i++) {
    getStateAsPhil()
  }
  console.timeEnd("PART2 compute as Phil suggested as array")

  console.time("PART2 compute as Phil suggested as string")
  for (let i = 0; i < amount; i++) {
    getStateAsPhilString()
  }
  console.timeEnd("PART2 compute as Phil suggested as string")
})