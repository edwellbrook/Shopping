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

    database.query('SELECT id, password FROM cards WHERE id = $1', [card], function(err, result) {
      if (err != null) {
        return next(err)
      }

      if (result.rowCount != 1) {
        let err = new Error('Invalid login')
        err.status = 401

        return next(err)
      }

      const hash = result.rows[0].password
      bcrypt.compare(pass, hash, function(err, success) {
        if (err != null) {
          return next(err)
        }

        if (success === false) {
          return res.render('customer/login', {
            attemptedLogin: true
          })
        }

        req.session.card = card
        res.redirect('/customer/list')
      })
    })
  })


  //
  // List
  //

  router.get('/list', function(req, res, next) {
    const cardId = req.session.card

    if (cardId == null) {
      let err = new Error('Must log in')
      err.status = 401

      return next(err)
    }

    database.query('SELECT id, list FROM cards WHERE id = $1', [cardId], function(err, result) {
      if (err != null) {
        return next(err)
      }

      if (result.rowCount != 1) {
        // something messed up in the database. force logout
        req.session.destroy()

        let err = new Error('Invalid number of rows returned')
        err.status = 500

        return next(err)
      }

      const list = result.rows[0].list
      res.render('customer/list', {
        list: list
      })
    })
  })

  return router
}
