// Сбор состояний чекбоксов
function collectCheckboxStates() {
  const checkboxes = document.querySelectorAll("input[type='checkbox']");
  const bitsArray = Array.from(checkboxes).map(checkbox => checkbox.checked ? 1 : 0);
  return bitsArray;
}

// Упаковка массива битов в байты
function packBits(bitsArray) {
  const bytes = [];
  for (let i = 0; i < bitsArray.length; i += 8) {
    let byte = 0;
    for (let j = 0; j < 8 && i + j < bitsArray.length; j++) {
      byte |= bitsArray[i + j] << (7 - j); // Сдвигаем биты, чтобы упаковать их в байт
    }
    bytes.push(byte);
  }
  return new Uint8Array(bytes); // Создаем массив байтов
}

// Отправка данных на сервер
function sendCheckboxStates() {
  const bitsArray = collectCheckboxStates();
  const packedBytes = packBits(bitsArray);

  fetch('/receive', {
    method: 'POST',
    headers: { 'Content-Type': 'application/octet-stream' },
    body: packedBytes
  }).then(response => response.text()).then(data => console.log(data));
}
