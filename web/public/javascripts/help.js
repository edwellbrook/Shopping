var lists = []

window.setInterval(function() {
  var d = new Date()
  d.setSeconds(d.getSeconds() - 10)

  console.log(lists)

  lists.forEach(function(list, idx, arr) {
    if (list.lastPing < d) {
      arr.splice(idx, 1)
      console.log('removing item')
    }
  })

  console.log(lists)
}, 1500)

mqtt
  .connect('<%= mqttAddress %>')
  .on('connect', function() {
    this.subscribe('/help')
    this.subscribe('/active')
  })
  .on('message', function(topic, data) {
    var json;

    try {
      json = JSON.parse(data)
    } catch (_) {}

    switch (topic) {
    case '/help':
      return handleHelp(json)
    case '/active':
      return handleActive(json)
    default:
      console.log('Missing message handler for topic:', topic)
    }
  })

function handleHelp(data) {
  console.log(data)
}

function handleActive(data) {
  lists.push({
    data: data,
    lastPing: new Date()
  })
}
