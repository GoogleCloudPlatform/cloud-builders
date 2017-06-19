const app = require('./index');
const expect = require('chai').expect;
const request = require('supertest');

describe('test.js', () => {

  describe('GET /', () => {

    it('responds with 200', (done) => {
      request(app)
        .get('/')
        .expect(200)
        .end((e, res) => {
          expect(e).to.not.exist;
          done();
        });
    });

  });

});
