const express = require('express')
const bcrypt = require('bcrypt')

const BCRYPT_ROUNDS = 10

module.exports = function(database) {
  const router = express.Router()

  //
  // Index
  //

  router.get('/', function(req, res, next) {
    res.render('customer/index', {})
  })


  //
  // Register
  //

  router.get('/register', function(req, res, next) {
    res.render('customer/register', {
      didRegister: false
    })
  })

  router.post('/register', function(req, res, next) {
    const card = (req.body.card_id || '').trim()
    const pass = (req.body.card_pass || '')

    if (card == '' || pass == '') {
      let err = new Error('Invalid Card ID or password')
      err.status = 400

      return next(err)
    }

    bcrypt.hash(pass, BCRYPT_ROUNDS, function(err, hash) {
      if (err != null) {
        return next(err)
      }

      database.query('INSERT INTO cards VALUES ($1, $2)', [card, hash], function(err) {
        if (err != null) {
          return next(err)
        }

        res.render('customer/register', {
          didRegister: true
        })
      })
    })
  })


  //
  // Login
  //

  router.get('/login', function(req, res, next) {
    res.render('customer/login', {
      attemptedLogin: false
    })
  })

  router.post('/login', function(req, res, next) {
    const card = (req.body.card_id || '').trim()
    const pass = (req.body.card_pass || '')

    if (card == '' || pass == '') {
      let err = new Error('Invalid Card ID or password')
      err.status = 400

      return next(err)
    }

    database.query('SELECT * FROM cards WHERE card_id = $1', [card], function(err, results) {
      if (err != null) {
        return next(err)
      }

      if (results.rowCount != 1) {
        let err = new Error('Invalid login')
        err.status = 401

        return next(err)
      }

      const hash = results.rows[0].password
      bcrypt.compare(pass, hash, function(err, success) {
        if (err != null) {
          return next(err)
        }

        if (success === false) {
          return res.render('customer/login', {
            attemptedLogin: true
          })
        }

        res.redirect('/customer/list')
      })
    })
  })


  //
  // List
  //

  router.get('/list', function(req, res, next) {
    res.send('This is where you can edit your shopping list.')
  })

  return router
}
