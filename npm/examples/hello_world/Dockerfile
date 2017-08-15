FROM gcr.io/google-appengine/nodejs

WORKDIR /hello

COPY package.json /hello/
RUN npm install
COPY . /hello/

EXPOSE 3000

CMD ["npm", "start"]
