const express = require('express')

module.exports = function(database) {
  const router = express.Router()

  router.get('/', function(req, res, next) {
    res.render('customer/index', {})
  })

  router.get('/register', function(req, res, next) {
    res.render('customer/register', {})
  })

  router.post('/register', function(req, res, next) {
    let card = req.body.card_id
    if (card == null || card.trim() == '') {
      let err = new Error('Invalid Card ID')
      err.status = 400

      return next(err)
    }

    database.query('INSERT INTO cards VALUES ($1, $2)', [card, -1], function(err) {
      if (err != null) {
        return next(err)
      }

      res.send('nice one bruvva')
    })
  })

  return router
}
