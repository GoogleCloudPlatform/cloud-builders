const app = require('./index');
const should = require('chai').should();
const request = require('supertest');

describe('test.js', () => {

  describe('GET /', () => {

    it('responds with 200', (done) => {
      request(app)
        .get('/')
        .expect(200)
        .end((e, res) => {
          should.not.exist(e);
          done();
        });
    });

  });

});
