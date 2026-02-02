# RUN
FROM node:14-alpine

WORKDIR /app
COPY package*.json ./
RUN npm install --production && npm cache clean --force
COPY . .
EXPOSE 3000
CMD ["node","app.js"]