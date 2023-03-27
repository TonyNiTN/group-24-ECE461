FROM ubuntu:20.04

WORKDIR /usr/src/app
COPY package*.json ./

RUN apt-get install golang-go golang-go.uber-zap-dev
RUN go get github.com/patrickmn/go-cache  github.com/spf13/viper github.com/google/go-github/github github.com/shurcooL/githubv4 golang.org/x/oauth2
RUN npm install

COPY . .

EXPOSE 8080

CMD ["npm", "run", "start"]