ARG NODE_VERSION=21.7.1

FROM node:${NODE_VERSION}-alpine

WORKDIR /usr/src/app

COPY package.json .

RUN npm install

# Copy the rest of the source files into the image.
COPY . .

# Expose the port that the application listens on.
EXPOSE 3000
EXPOSE 3001

# Run the application.
CMD ["npm", "run", "dev"]