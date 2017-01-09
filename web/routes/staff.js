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

  router.get('/items', function(req, res, next) {
    database.query('SELECT * FROM items', function(err, result) {
      if (err != null) {
        return next(err)
      }

      res.render('staff/items', {
        items: result.rows
      })
    })
  })

  router.get('/items/add', function(req, res, next) {
    res.render('staff/add', {})
  })

  router.post('/items/add', function(req, res, next) {
    const name = (req.body.item_name || '').trim()
    const location = (req.body.item_location || '').trim()

    if (name == '' || location == '') {
      let err = new Error('Invalid name or location')
      err.status = 400

      return next(err)
    }

    database.query('INSERT INTO items VALUES ($1, $2)', [name, location], function(err, result) {
      if (err != null) {
        return next(err)
      }

      res.redirect('/staff/items')
    })
  })

  router.get('/help', function(req, res, next) {
    res.render('staff/help', {
      mqttAddress: `ws://${ process.env.MQTT || 'localhost:15675' }/ws`
    })
  })

  return router
}
