var requests = {}

window.setInterval(function() {
  var d = new Date()
  d.setSeconds(d.getSeconds() - 5)

  Object.keys(requests).forEach(function(id) {
    var req = requests[id]

    if (req.ping < d) {
      $(`[data-id=${req.id}]`).remove()
      delete requests[id]
    }
  })
}, 1500)

mqtt
  .connect(window.mqttAddress)
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
  req = {
    id: data.id,
    location: data.location,
    ping: new Date()
  }

  if (requests[req.id]) {
    $(`#helper [data-id=${req.id}] [data-loc]`).text(req.location)
  } else {
    $('#helper').append(`<li data-id='${req.id}'><strong><code>${req.id}</code></strong>: <code data-loc>${req.location}</code></li>`)
  }

  requests[req.id] = req
}

function handleActive(data) {
  lists.push({
    data: data,
    lastPing: new Date()
  })
}
