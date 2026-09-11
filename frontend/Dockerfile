# ---- build ----
FROM node:24-alpine AS build
WORKDIR /app

# npm ci, not npm install: installs exactly the lockfile and fails if the
# lockfile and manifest disagree.
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# ---- serve ----
# The unprivileged image runs as non-root and therefore listens on 8080,
# not 80, since a non-root process cannot bind a privileged port.
FROM nginxinc/nginx-unprivileged:1.31-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
