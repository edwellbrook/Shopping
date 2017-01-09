const express = require('express')

module.exports = function(database) {
  const router = express.Router()

  router.get('/', function(req, res, next) {
    res.render('staff/index', {})
  })

  router.get('/cards', function(req, res, next) {
    database.query('SELECT * FROM cards', function(err, result) {
      if (err != null) {
        return next(err)
      }

      res.render('staff/cards', {
        cards: result.rows
      })
    })
  })

  router.get('/help', function(req, res, next) {
    res.render('staff/help', {
      mqttAddress: `ws://${ process.env.MQTT || 'localhost:15675' }/ws`
    })
  })

  return router
}
