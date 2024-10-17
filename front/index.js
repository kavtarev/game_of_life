const field = document.getElementById('field');
const inputs = []

for (let i = 0; i < 100000; i++) {
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
  const str = new Array(100000).fill('1').join('')

  console.time('start')
  for (let i = 0; i < str.length; i++) {
    if (str[i] === '1') {
      inputs[i].checked = true
    }
  }
  console.timeEnd('start')

}

run()