# Tahap 1: Build Aplikasi React
FROM node:18 AS build

# Set work directory di dalam container
WORKDIR /app

# Salin file package.json dan package-lock.json (atau yarn.lock) untuk mengoptimalkan caching dependensi
COPY package*.json ./

# Install dependencies
RUN npm install

# Salin semua file aplikasi ke dalam container
COPY . .

# Build aplikasi untuk produksi
RUN npm run build

# Tahap 2: Menyajikan Aplikasi dengan Nginx
FROM nginx:alpine

# Salin hasil build dari tahap 1 ke dalam direktori yang sesuai di nginx
COPY --from=build /app/build /usr/share/nginx/html

# Salin konfigurasi NGINX
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80 untuk akses web
EXPOSE 3000

# Gunakan Nginx untuk menyajikan aplikasi
CMD ["nginx", "-g", "daemon off;"]
