const express = require('express')

module.exports = function(database) {
  const router = express.Router()

  router.get('/', function(req, res, next) {
    res.send('respond with a resource')
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

  return router
}
