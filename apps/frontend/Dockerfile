# Stage 1: Build Node
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Serve Nginx
FROM nginx:alpine
# Видаляємо дефолтний конфіг
RUN rm /etc/nginx/conf.d/default.conf
# Копіюємо наш конфіг
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Копіюємо зібраний фронт
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]