const { request } = require('http');

const buff = new ArrayBuffer(10);



for (let i = 0; i < 10; i++) {
  if (i & 1) {
    buff[i] = 12
  } else {
    buff[i] = 10
  }
}

console.log(buff);

function foo() {
  var post_options = {
    host: '127.0.0.1',
    port: '3000',
    path: '/buff',
    method: 'POST',
  };

  var post_req = request(post_options, function (res) {
    res.on('data', function (chunk) {
      console.log('Response: ' + chunk);
    });
  });

  // post the data
  post_req.write(new Uint8Array(buff));
  post_req.end();
}

foo()